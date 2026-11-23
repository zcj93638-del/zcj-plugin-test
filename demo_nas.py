"""Local fixture for testing the Seek Codex plugin."""
import hashlib
import json
import secrets
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer


class DemoNAS(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != "/seek/api/v1/auth/user/loginByPassword":
            return self.send_error(404)
        size = int(self.headers.get("Content-Length", 0))
        if not 0 < size < 16384:
            return self.send_error(400)
        credentials = json.loads(self.rfile.read(size))
        if (credentials.get("username") != "demo" or not credentials.get("deviceId")
                or credentials.get("password") != hashlib.md5(b"seek-demo").hexdigest()):
            return self.send_error(401)
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Set-Cookie", "seek-token=" + self.server.token + "; Path=/; HttpOnly")
        self.end_headers()
        self.wfile.write(b'{"code":0,"data":null}')

    def do_GET(self):
        if self.path != "/seek/api/v1/project/list":
            return self.send_error(404)
        if self.headers.get("Cookie") != "seek-token=" + self.server.token:
            return self.send_error(401)
        projects = [{"name": "长夜 · 演示项目", "guid": "demo-night"}]
        body = json.dumps({"code": 0, "msg": "", "data": {"list": projects}}).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


if __name__ == "__main__":
    with ThreadingHTTPServer(("127.0.0.1", 0), DemoNAS) as server:
        server.token = secrets.token_urlsafe(24)
        print(f"http://127.0.0.1:{server.server_port}/seek", flush=True)
        server.serve_forever()
