# Trim Seek for Codex 和 DeepSeek Harness

在 Codex 或 DeepSeek Harness 中连接 NAS 上的 Trim Seek，查询项目并使用 Seek OpenAPI。插件已包含运行程序，支持 Windows、macOS 和 Linux 的 x64 和 ARM64 设备。用户无需另装 Python 或 Go。

## 在 DeepSeek Harness 安装

打开「插件」，选择从 Git 仓库安装，粘贴 `https://github.com/zcj93638-del/zcj-plugin-test`。安装并启用后，打开「设置 → Trim Seek」，输入 NAS 地址、Seek 账号和密码。看到「登录成功」后，就可以在聊天里问「我有哪些项目？」。

同一系统用户的 Harness 和 Codex 共享 Seek 登录会话。在一边切换 NAS 或账号，另一边下次查询也会使用新连接。Harness 会自动选择对应平台的程序。

## 在 Codex 安装

在 Codex 添加本仓库作为插件市场，按运行 Codex 的系统安装一个插件。原生 Windows 选 **Trim Seek (Windows)**，Mac 和 Linux（包括 WSL 内运行的 Codex）选 **Trim Seek (macOS / Linux)**。

也可以在终端添加市场：

```sh
codex plugin marketplace add https://github.com/zcj93638-del/zcj-plugin-test.git
```

Windows 在 PowerShell 中运行：

```powershell
codex plugin add seek-windows@trim-local
```

macOS / Linux 运行：

```sh
codex plugin add seek@trim-local
```

已装 0.1.7 的用户先刷新市场 `codex plugin marketplace upgrade trim-local`，移除旧插件 `codex plugin remove seek@trim-local`，再安装对应平台的版本并重开 Codex。若市场来自本地解压目录，请先更新目录中的文件。升版可沿用已有的 NAS 会话。

## 连接 NAS

打开插件的「管理 → Seek 连接」，填入 NAS 地址、Seek 账号，然后输入密码登录。只填 IP 时默认使用 HTTP 和 5666 端口；完整 URL 使用指定协议和端口。成功后可以问「我有哪些项目？」。密码仅用于本次登录，会话保存在当前用户的 `.config/trim-seek-codex/connection.json` 中。HTTP 连接不会加密登录摘要和 Cookie。

若客户端没有显示表单，也可以从仓库解压目录运行 `plugins\seek-windows\bin\start.cmd setup`（Windows）或 `sh plugins/seek/bin/start.sh setup`（macOS / Linux），打开终端打印的本机地址填写连接。密码不用发到对话里。

## 维护者构建和测试

在 macOS / Linux 安装 Go 后运行 `./scripts/build.sh` 可为两个安装包生成原生程序并同步技能和接口文档。`go test ./...` 会通过当前平台包的实际启动命令验证首次读取设置、登录、重启后查询项目和切换账号。发布时将生成的两个插件目录一并提交。根目录的 `demo_nas.py` 是开发者手动测试用的模拟 NAS。
