package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"time"
)

//go:embed settings.html
var settingsHTML string

func setup(s store) error {
	page := template.Must(template.New("settings").Parse(settingsHTML))
	id, err := randomID()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer l.Close()
	host := l.Addr().String()
	origin := "http://" + host
	path := "/" + id
	ctx, stop := context.WithTimeout(context.Background(), 10*time.Minute)
	defer stop()
	server := &http.Server{ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, MaxHeaderBytes: 16384}
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path || r.Host != host {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self'; frame-ancestors 'none'")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		c, err := s.read()
		if err != nil {
			http.Error(w, "无法读取配置", 500)
			return
		}
		message := ""
		switch r.Method {
		case "GET":
		case "POST":
			if r.Header.Get("Origin") != origin {
				http.Error(w, "Forbidden", 403)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, 16384)
			if r.ParseForm() != nil {
				http.Error(w, "Invalid form", 400)
				return
			}
			values := map[string]string{}
			for _, key := range []string{"address", "username", "password"} {
				values[key] = r.PostForm.Get(key)
				if values[key] == "" {
					http.Error(w, incomplete, 400)
					return
				}
			}
			c, err = s.configure(values)
			if err != nil {
				message = err.Error()
			} else {
				message = "已登录，可以关闭此页。"
				defer stop()
			}
		default:
			http.Error(w, "Method not allowed", 405)
			return
		}
		_ = page.Execute(w, map[string]string{"Address": c.Address, "Username": c.Username, "Message": message})
	})
	fmt.Println(origin + path)
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		end, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(end)
		close(done)
	}()
	err = server.Serve(l)
	stop()
	<-done
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
