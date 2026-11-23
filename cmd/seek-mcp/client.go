package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/flock"
)

const incomplete = "请在插件设置中填好 NAS 地址、账号和密码。"

type connection struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	DeviceID string `json:"deviceId"`
	Token    string `json:"token,omitempty"`
}
type store struct{ path string }

func (s store) read() (connection, error) {
	var c connection
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	return c, err
}
func (s store) save(c connection) error {
	f, err := os.CreateTemp(filepath.Dir(s.path), "connection-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	err = json.NewEncoder(f).Encode(c)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), s.path)
}
func normalize(address string) (string, error) {
	address = strings.TrimSpace(address)
	bare := !strings.Contains(address, "://")
	if bare {
		address = "http://" + address
	}
	u, err := url.Parse(address)
	if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || strings.ContainsAny(address, "?#\\") || strings.IndexFunc(address, unicode.IsSpace) >= 0 {
		return "", errors.New("请输入 IP、域名或完整的 HTTP(S) 服务地址。")
	}
	if bare && u.Port() == "" {
		u.Host = net.JoinHostPort(u.Hostname(), "5666")
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/seek"
	}
	return strings.TrimRight(u.String(), "/"), nil
}
func (c connection) values() obj {
	return obj{"address": c.Address, "username": c.Username, "password": ""}
}
func (c connection) ready() bool { return c.Address != "" && c.Username != "" && c.Token != "" }

var client = &http.Client{Transport: &http.Transport{}, Timeout: 7 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func (c connection) request(method, path string, body json.RawMessage) (json.RawMessage, []*http.Cookie, error) {
	req, err := http.NewRequest(method, c.Address+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.AddCookie(&http.Cookie{Name: "seek-token", Value: strings.Trim(c.Token, `"`)})
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("连接 Seek 失败：%w", err)
	}
	defer r.Body.Close()
	if r.StatusCode == 401 {
		return nil, nil, errors.New("认证失败，请重新输入密码。")
	}
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("Seek 返回 HTTP %d，请核对地址、端口和权限。", r.StatusCode)
	}
	var envelope struct {
		Code *int            `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err = json.NewDecoder(io.LimitReader(r.Body, 32<<20)).Decode(&envelope); err != nil || envelope.Code == nil {
		return nil, nil, errors.New("响应不是 Seek JSON，请检查服务地址。")
	}
	if *envelope.Code != 0 {
		return nil, nil, fmt.Errorf("Seek 拒绝请求（业务码 %d）。", *envelope.Code)
	}
	return envelope.Data, r.Cookies(), nil
}
func randomID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}
func (s store) configure(fields map[string]string) (connection, error) {
	var c connection
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return c, err
	}
	lock := flock.New(filepath.Join(filepath.Dir(s.path), ".settings.lock"))
	if err := lock.Lock(); err != nil {
		return c, err
	}
	defer lock.Unlock()
	c, err := s.read()
	if err != nil {
		return c, err
	}
	old := c
	for k, v := range fields {
		switch k {
		case "address":
			c.Address = strings.TrimSpace(v)
			if c.Address != "" {
				c.Address, err = normalize(c.Address)
				if err != nil {
					return c, err
				}
			}
		case "username":
			c.Username = strings.TrimSpace(v)
		case "password":
		default:
			return c, errors.New("连接字段无效。")
		}
	}
	password := fields["password"]
	if old.Address != c.Address || old.Username != c.Username || password != "" {
		c.Token = ""
	}
	if c.DeviceID == "" {
		c.DeviceID, err = randomID()
		if err != nil {
			return c, err
		}
	}
	if err = s.save(c); err != nil || password == "" {
		return c, err
	}
	if c.Username == "" || c.Address == "" {
		return c, errors.New(incomplete)
	}
	payload, _ := json.Marshal(obj{"username": c.Username, "password": fmt.Sprintf("%x", md5.Sum([]byte(password))), "deviceId": c.DeviceID})
	_, cookies, err := c.request("POST", "/api/v1/auth/user/loginByPassword", payload)
	if err != nil {
		return c, err
	}
	for _, cookie := range cookies {
		if cookie.Name == "seek-token" && cookie.Value != "" && (!cookie.Secure || strings.HasPrefix(c.Address, "https://")) {
			c.Token = cookie.Value
			return c, s.save(c)
		}
	}
	return c, errors.New("登录响应缺少可用于此地址的 seek-token。")
}
func (c connection) api(method, path string, body json.RawMessage) (json.RawMessage, error) {
	if !c.ready() {
		return nil, errors.New(incomplete)
	}
	u, err := url.ParseRequestURI(path)
	if err != nil || !strings.HasPrefix(path, "/api/v1/") || strings.Contains(u.Path, "\\") || strings.Contains("/"+u.Path+"/", "/../") {
		return nil, errors.New("请输入 /api/v1/ 下的 Seek 接口路径。")
	}
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		return nil, errors.New("请求方法无效。")
	}
	data, _, err := c.request(method, path, body)
	return data, err
}
