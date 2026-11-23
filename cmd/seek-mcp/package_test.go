package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// Exercise the shipped executable via the same cwd, command and arguments as Codex.
func TestPackagedLoginAndProjects(t *testing.T) {
	dir := filepath.Join("..", "..", "plugins", "seek")
	if runtime.GOOS == "windows" {
		dir += "-windows"
	}
	if root := os.Getenv("SEEK_TEST_INSTALLED_ROOT"); root != "" {
		dir = root
	}
	b, err := os.ReadFile(filepath.Join(dir, ".mcp.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		McpServers map[string]struct {
			Command string
			Args    []string
			Cwd     string
		}
	}
	if err = json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	launch := manifest.McpServers["seek"]
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	content := []byte{0, 255, 13, 10, 42}
	filename := "上传 测试.bin"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/seek/api/v1/auth/user/loginByPassword":
			var creds map[string]string
			_ = json.NewDecoder(r.Body).Decode(&creds)
			if r.Method != "POST" || creds["username"] != "demo" || creds["deviceId"] == "" || creds["password"] != fmt.Sprintf("%x", md5.Sum([]byte("test-password"))) {
				http.Error(w, "invalid login", 401)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "seek-token", Value: "test-session", HttpOnly: true})
			fmt.Fprint(w, `{"code":0,"data":null}`)
		case "/seek/upload/api/v1/asset":
			cookie, err := r.Cookie("seek-token")
			if r.Method != "POST" || err != nil || cookie.Value != "test-session" || r.ContentLength <= int64(len(content)) {
				http.Error(w, "invalid upload headers", 400)
				return
			}
			parts, err := r.MultipartReader()
			if err != nil {
				http.Error(w, "need multipart", 400)
				return
			}
			fields := map[string]string{}
			for {
				part, err := parts.NextPart()
				if err != nil {
					http.Error(w, "missing file", 400)
					return
				}
				data, err := io.ReadAll(part)
				if err != nil {
					http.Error(w, "incomplete part", 400)
					return
				}
				if part.FormName() != "file" {
					fields[part.FormName()] = string(data)
					continue
				}
				if part.FileName() != filename || !bytes.Equal(data, content) || fields["size"] != strconv.Itoa(len(content)) || fields["projectGuid"] != "abc" || fields["parentFolderGuid"] != "folder" || fields["onNameConflict"] != "keepBoth" {
					http.Error(w, "invalid multipart fields or file", 400)
					return
				}
				break
			}
			if _, err = parts.NextPart(); err != io.EOF {
				http.Error(w, "file must be last", 400)
				return
			}
			fmt.Fprint(w, `{"code":0,"data":{"guid":"uploaded-id"}}`)
		case "/seek/api/v1/project/list":
			cookie, e := r.Cookie("seek-token")
			if e != nil || cookie.Value != "test-session" {
				http.Error(w, "no session", 401)
				return
			}
			fmt.Fprint(w, `{"code":0,"data":{"list":[{"name":"Demo project","guid":"abc"}]}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	// Restart the plugin for each request to cover persistence across Codex sessions.
	invoke := func(method string, params obj) obj {
		t.Helper()
		data, _ := json.Marshal(obj{"id": 1, "jsonrpc": "2.0", "method": method, "params": params})
		cmd := exec.Command(launch.Command, launch.Args...)
		cmd.Dir = filepath.Join(dir, launch.Cwd)
		cmd.Stdin = bytes.NewReader(append(data, '\n'))
		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("start package: %v %s", err, b)
		}
		var result struct {
			Result obj
			Error  any
		}
		if err = json.Unmarshal(b, &result); err != nil || result.Error != nil || result.Result["isError"] == true {
			t.Fatalf("%s: %s", method, b)
		}
		return result.Result
	}
	init := invoke("initialize", obj{"protocolVersion": "2025-11-25"})
	caps := init["capabilities"].(obj)["experimental"].(obj)["openai/settings"].(obj)
	if caps["readTool"] != readTool || caps["updateTool"] != updateTool {
		t.Fatal(caps)
	}
	tool := func(name string, args obj) obj {
		return invoke("tools/call", obj{"name": name, "arguments": args})["structuredContent"].(obj)
	}
	settings := tool(readTool, obj{})
	values := settings["values"].(obj)
	if len(values) != 3 || values["address"] != "" || values["username"] != "" || values["password"] != "" {
		t.Fatal(values)
	}
	for _, key := range []string{"address", "username", "password"} {
		if settings["schema"].(obj)["properties"].(obj)[key] == nil {
			t.Fatalf("missing %s", key)
		}
	}
	tool(updateTool, obj{"set": obj{"address": server.URL, "username": "demo"}})
	tool(updateTool, obj{"set": obj{"password": "test-password"}})
	if tool("seek_status", obj{})["configured"] != true {
		t.Fatal("session not saved")
	}
	projects := tool("seek_projects", obj{})["list"].([]any)
	if len(projects) != 1 || projects[0].(obj)["name"] != "Demo project" {
		t.Fatal(projects)
	}
	data := tool("seek_api", obj{"method": "GET", "path": "/api/v1/project/list"})
	if data["data"].(obj)["list"] == nil {
		t.Fatal(data)
	}
	source := filepath.Join(home, filename)
	if err := os.WriteFile(source, content, 0600); err != nil {
		t.Fatal(err)
	}
	if result := tool("seek_upload", obj{"projectGuid": "abc", "parentFolderGuid": "folder", "filePath": source}); result["guid"] != "uploaded-id" {
		t.Fatal(result)
	}
	saved, err := os.ReadFile(filepath.Join(home, ".config", "trim-seek-codex", "connection.json"))
	if err != nil || strings.Contains(string(saved), "test-password") {
		t.Fatal("unexpected stored credentials")
	}
	tool(updateTool, obj{"set": obj{"username": "other"}})
	if tool("seek_status", obj{})["configured"] != false {
		t.Fatal("old session remained active")
	}
}
