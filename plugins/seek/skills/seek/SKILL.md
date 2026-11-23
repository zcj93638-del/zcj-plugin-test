---
name: seek
description: 配置 Trim Seek 连接、查询项目，以及根据 Seek OpenAPI 查找素材搜索、文件夹、标签、评论和分享等接口的调用方法。
---

- 用户尚未配置时只提示「请在插件设置中填好 NAS 地址、账号和密码」。密码由用户在设置页面输入，勿通过对话索取。
- 在处理 Seek 业务请求前，调用 seek_status 检查 configured，若未就绪，返回提示并停止。设置页可用 seek_settings_read 查看地址和账号。避免读取原始会话配置文件。
- 查询项目数量或列表时，调用 seek_projects。其他 JSON 业务接口可通过 seek_api 传入 method、path 和可选的 payload，使用已保存的会话。依据实际响应作答。

- 用户要求上传素材时使用 `seek_upload`。仅传用户指定的文件在插件运行机器上可读的绝对路径，先解析目标项目和项目内的文件夹 GUID。聊天附件没有可用路径时请用户提供路径。一次上传一个文件，同名默认 `keepBoth`；只有用户要求覆盖时传 `overwrite`。`skip` 返回的 GUID 也可能属于已有文件。

## 按任务查找 API

用户询问 Seek 的能力、请求示例或需要调用接口时，先检索本技能的 [接口索引](references/api-index.md)，再按行号读取[随包 OpenAPI 文档](references/trim-seek.openapi.yaml)中的路径条目。例如，按关键字跨项目找素材可查“全局标准搜索”，按画面描述找素材可查“执行 AI 搜索”，给素材打标签可查“修改标签”。

从所选操作和路径级 parameters 确认请求方法、参数位置、必填项、请求体和媒体类型。遇到 `$ref` 时在同一文档中查对应的 components 定义。按操作级或全局 security 说明认证方式。接口文档可离线查阅；返回示例中的项目名、数量和 GUID 仅用于解释格式。

实际请求使用用户配置的服务地址和文档定义的路径。普通 Seek API 在配置好的 /seek 基址后追加 /api/v1/...；分享页按对应的分享 URL 和网关路径构造。使用设置中已登录的目标。按实际 HTTP 状态、业务 code 和响应数据报告结果。
