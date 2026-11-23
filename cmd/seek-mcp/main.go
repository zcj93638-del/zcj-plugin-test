package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type obj = map[string]any

const version = "0.4.0"
const readTool = "seek_settings_read"
const updateTool = "seek_settings_update"

func schema(properties obj, required ...string) obj {
	s := obj{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

var fields = obj{"address": obj{"type": "string", "title": "NAS 地址", "description": "例如 http://192.168.9.118:5666/seek/。仅填 IP 时默认使用 HTTP 和端口 5666。"}, "username": obj{"type": "string", "title": "Seek 账号", "description": "当前测试阶段暂用 Seek 账号和密码登录，后续将改用 API Key。"}, "password": obj{"type": "string", "title": "密码", "description": "填好地址和账号后输入密码登录。"}}

func tools() []obj {
	object := obj{"type": "object"}
	set := schema(fields)
	set["minProperties"] = 1
	readOnly := obj{"readOnlyHint": true}
	return []obj{
		uploadTool(),
		{"name": readTool, "description": "读取 Seek 连接设置。", "annotations": readOnly, "inputSchema": schema(obj{}), "outputSchema": schema(obj{"schema": object, "values": object, "layout": obj{"type": "array", "items": object}}, "schema", "values")},
		{"name": updateTool, "description": "更新连接设置，提交密码时登录。", "inputSchema": schema(obj{"set": set}, "set"), "outputSchema": schema(obj{"values": object}, "values")},
		{"name": "seek_status", "description": "检查本地是否已配置登录会话。", "annotations": readOnly, "inputSchema": schema(obj{})},
		{"name": "seek_projects", "description": "查询可见项目的名称和 GUID。", "annotations": readOnly, "inputSchema": schema(obj{})},
		{"name": "seek_api", "description": "使用已登录会话调用 Seek JSON 接口。", "inputSchema": schema(obj{"method": obj{"type": "string", "enum": []string{"GET", "POST", "PUT", "PATCH", "DELETE"}}, "path": obj{"type": "string"}, "payload": obj{}}, "method", "path")},
	}
}

type callParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func (s store) call(p callParams) (obj, error) {
	if p.Name == updateTool {
		var args struct {
			Set map[string]string `json:"set"`
		}
		if json.Unmarshal(p.Arguments, &args) != nil || len(args.Set) == 0 {
			return nil, errors.New("连接字段无效。")
		}
		c, err := s.configure(args.Set)
		if err != nil {
			return nil, err
		}
		return obj{"values": c.values()}, nil
	}
	c, err := s.read()
	if err != nil {
		return nil, err
	}
	switch p.Name {
	case readTool:
		items := []obj{}
		for _, k := range []string{"address", "username", "password"} {
			items = append(items, obj{"kind": "property", "property": k})
		}
		return obj{"schema": schema(fields), "values": c.values(), "layout": []obj{{"kind": "group", "title": "Seek 连接", "items": items}}}, nil
	case "seek_status":
		return obj{"configured": c.ready()}, nil
	case "seek_upload":
		var a uploadArgs
		if json.Unmarshal(p.Arguments, &a) != nil {
			return nil, errors.New("上传参数无效。")
		}
		return c.upload(a)
	case "seek_projects":
		data, err := c.api("GET", "/api/v1/project/list", nil)
		if err != nil {
			return nil, err
		}
		var projects struct {
			List []struct {
				Name *string `json:"name"`
				GUID string  `json:"guid"`
			} `json:"list"`
		}
		if json.Unmarshal(data, &projects) != nil || projects.List == nil {
			return nil, errors.New("项目列表响应格式无效。")
		}
		for _, p := range projects.List {
			if p.Name == nil {
				return nil, errors.New("项目名称缺失。")
			}
		}
		return obj{"list": projects.List}, nil
	case "seek_api":
		var a struct {
			Method, Path string
			Payload      json.RawMessage
		}
		if json.Unmarshal(p.Arguments, &a) != nil {
			return nil, errors.New("请求参数无效。")
		}
		data, err := c.api(a.Method, a.Path, a.Payload)
		return obj{"data": data}, err
	default:
		return nil, errors.New("找不到指定工具。")
	}
}

type message struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func rpcError(code int, msg string) obj { return obj{"code": code, "message": msg} }
func reply(m message, s store) obj {
	res := obj{"jsonrpc": "2.0", "id": m.ID}
	switch m.Method {
	case "initialize":
		var p struct{ ProtocolVersion string }
		_ = json.Unmarshal(m.Params, &p)
		v := p.ProtocolVersion
		switch v {
		case "2024-11-05", "2025-03-26", "2025-06-18", "2025-11-25":
		default:
			v = "2025-11-25"
		}
		res["result"] = obj{"protocolVersion": v, "serverInfo": obj{"name": "trim-seek", "version": version}, "capabilities": obj{"tools": obj{}, "experimental": obj{"openai/settings": obj{"readTool": readTool, "updateTool": updateTool}}}}
	case "ping":
		res["result"] = obj{}
	case "tools/list":
		res["result"] = obj{"tools": tools()}
	case "tools/call":
		var p callParams
		if err := json.Unmarshal(m.Params, &p); err != nil {
			res["error"] = rpcError(-32602, "Invalid params")
			break
		}
		result, err := s.call(p)
		if err != nil {
			res["result"] = obj{"isError": true, "content": []obj{{"type": "text", "text": err.Error()}}}
			break
		}
		text, _ := json.Marshal(result)
		res["result"] = obj{"structuredContent": result, "content": []obj{{"type": "text", "text": string(text)}}}
	default:
		res["error"] = rpcError(-32601, "Method not found")
	}
	return res
}
func serve(in io.Reader, out io.Writer, s store) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 65536), 4<<20)
	encoder := json.NewEncoder(out)
	for scanner.Scan() {
		var m message
		var res obj
		if json.Unmarshal(scanner.Bytes(), &m) != nil {
			res = obj{"jsonrpc": "2.0", "id": nil, "error": rpcError(-32700, "Parse error")}
		} else if len(m.ID) == 0 {
			continue
		} else {
			res = reply(m, s)
		}
		if err := encoder.Encode(res); err != nil {
			return err
		}
	}
	return scanner.Err()
}
func run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	s := store{filepath.Join(home, ".config", "trim-seek-codex", "connection.json")}
	if len(os.Args) > 1 {
		if len(os.Args) == 2 && os.Args[1] == "setup" {
			return setup(s)
		}
		return errors.New("usage: seek-mcp [setup]")
	}
	return serve(os.Stdin, os.Stdout, s)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
