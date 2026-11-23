# Seek 接口索引

按业务名称、接口说明或路径检索，行号对应同目录的 trim-seek.openapi.yaml。

## 登录

- 密码登录 · `POST /api/v1/auth/user/loginByPassword` · YAML 第 28 行
- 授权登录 · `POST /api/v1/auth/user/loginByAuth` · YAML 第 79 行
- 退出登陆 · `POST /api/v1/auth/user/logout` · YAML 第 132 行

## 用户管理

- 管理员创建用户 · `POST /api/v1/manager/user/create` · YAML 第 169 行
- 管理员看用户列表 · `GET /api/v1/manager/user/list` · YAML 第 226 行
- 管理更新用户 · `POST /api/v1/manager/user/update` · YAML 第 276 行
- 管理删除用户 · `POST /api/v1/manager/user/delete` · YAML 第 328 行
- 获取用户信息 · `GET /api/v1/user/user/info` · YAML 第 366 行
- 用户更新信息 · `POST /api/v1/user/user/update` · YAML 第 432 行
- 上传用户头像 · `POST /upload/api/v1/user/avatar` · YAML 第 476 行
- 内部页面读取用户或访客头像 · `GET /api/v1/avatar/{ownerType}/{ownerGuid}/{avatar}` · YAML 第 601 行
- 用户项目列表 · `GET /api/v1/manager/user/projects` · YAML 第 626 行
- 管理解封用户 · `POST /api/v1/manager/user/unlock` · YAML 第 679 行
- 设置用户配置 · `POST /api/v1/user/user/setUserConfig` · YAML 第 10471 行

## 分享页

- 获取外链访客资料 · `GET /s/{shareId}/busapi/seek/api/v1/guest/info` · YAML 第 496 行
- 更新外链访客昵称和头像 · `POST /s/{shareId}/busapi/seek/api/v1/guest/update` · YAML 第 512 行
- 上传外链访客头像 · `POST /s/{shareId}/busapi/seek/api/v1/guest/avatar` · YAML 第 541 行
- 分享页读取用户或访客头像 · `GET /s/{shareId}/busapi/seek/api/v1/avatar/{ownerType}/{ownerGuid}/{avatar}` · YAML 第 569 行

## 播放器

- 单个视频proxy 生成 · `POST /api/v1/player/genProxy` · YAML 第 717 行
- 查询视频proxy状况 · `GET /api/v1/player/queryProxyStatus` · YAML 第 775 行
- 生成实时转码链接 · `POST /api/v1/player/getTranscodeLink` · YAML 第 872 行
- 获取实时转码统计信息 · `POST /api/v1/player/getTranscodeStatis` · YAML 第 997 行
- 实时转码退出播放 · `POST /api/v1/player/quit` · YAML 第 1103 行
- 实时转码 record · `POST /api/v1/player/record` · YAML 第 1184 行
- 原片range播 · `GET /api/v1/player/original` · YAML 第 1265 行
- proxy range播 · `GET /api/v1/player/proxy` · YAML 第 1320 行
- 访问m3u8 & ts文件 · `GET /media/{path}` · YAML 第 1374 行
- 进度条小视频播 · `GET /api/v1/player/preview` · YAML 第 1408 行
- 获取音乐实时转码链接 · `POST /api/v1/player/getMusicPlayLink` · YAML 第 10564 行

## 项目

- 项目信息 · `GET /api/v1/project/info` · YAML 第 1437 行
- 上传logo · `POST /api/v1/project/uploadLogo` · YAML 第 1532 行
- 获取项目logo(文件流) · `GET /api/v1/project/logo/{projectGuid}/{logoGuid}` · YAML 第 1578 行
- 创建项目 · `POST /api/v1/project/create` · YAML 第 1607 行
- 更新项目 · `POST /api/v1/project/update` · YAML 第 1684 行
- 项目列表 · `GET /api/v1/project/list` · YAML 第 1760 行
- 将用户加入项目 · `POST /api/v1/project/addUser` · YAML 第 1834 行
- 将用户移出项目 · `POST /api/v1/project/removeUser` · YAML 第 1881 行
- 获取项目用户列表 · `GET /api/v1/project/users` · YAML 第 1928 行
- 删除项目 · `POST /api/v1/project/delete` · YAML 第 1977 行
- 删除项目proxy文件 · `POST /api/v1/project/deleteProxyFiles` · YAML 第 2020 行
- 删除项目垃圾文件 · `POST /api/v1/project/deleteUselessFiles` · YAML 第 2059 行
- 统计项目垃圾文件 · `GET /api/v1/project/statUselessFiles` · YAML 第 2098 行
- 重连项目 · `POST /api/v1/project/reconnect` · YAML 第 2150 行
- 获取有项目权限的用户列表 · `GET /api/v1/project/permissionUsers` · YAML 第 7686 行
- 项目详情 · `GET /api/v1/project/detail` · YAML 第 10357 行
- 导出项目 · `GET /api/v1/project/export` · YAML 第 10384 行

## 文件夹

- 新建 · `POST /api/v1/folder/create` · YAML 第 2201 行
- 批量新建 · `POST /api/v1/folder/batchCreate` · YAML 第 2255 行
- 重命名 · `POST /api/v1/folder/rename` · YAML 第 2329 行
- 扫描、刷新 · `POST /api/v1/folder/scan` · YAML 第 2375 行
- 文件夹基础信息 · `GET /api/v1/folder/info` · YAML 第 2409 行
- 获取文件夹封面信息（批量） · `POST /api/v1/folder/thumbnails` · YAML 第 2480 行
- 文件夹父目录信息 · `GET /api/v1/folder/parent` · YAML 第 2644 行
- 智能文件夹列表 · `GET /api/v1/collection/list` · YAML 第 11728 行
- 智能文件夹详情 · `GET /api/v1/collection/info` · YAML 第 11755 行
- 创建智能文件夹 · `POST /api/v1/collection/create` · YAML 第 11782 行
- 更新智能文件夹 · `POST /api/v1/collection/update` · YAML 第 11815 行
- 智能文件夹关联外链数量 · `POST /api/v1/collection/relatedShareLinkCount` · YAML 第 11848 行
- 删除智能文件夹 · `POST /api/v1/collection/delete` · YAML 第 11881 行

## asset

- 文件夹详情 · `GET /api/v1/asset/folderDetail` · YAML 第 2786 行
- 批量复制 · `POST /api/v1/asset/copy` · YAML 第 2820 行
- 批量移动 · `POST /api/v1/asset/move` · YAML 第 2863 行
- 上传 · `POST /upload/api/v1/asset` · YAML 第 2905 行
- 上传文件封面 · `POST /upload/api/v1/assetCover` · YAML 第 2968 行
- 重命名 · `POST /api/v1/asset/rename` · YAML 第 3019 行
- 列表 · `POST /api/v1/asset/list` · YAML 第 3067 行
- 文件详情 · `GET /api/v1/asset/detail` · YAML 第 3154 行
- 文件缩略图 · `GET /api/v1/asset/thumbnail/{size}/{assetVersionGuid}` · YAML 第 3856 行
- 获取文件封面 · `GET /api/v1/asset/cover/{assetVersionGuid}` · YAML 第 3896 行
- 删除文件封面 · `POST /api/v1/asset/removeCover` · YAML 第 3921 行
- 文件雪碧图 · `GET /api/v1/asset/sprite/{assetVersionGuid}` · YAML 第 3961 行
- 合并 · `POST /api/v1/asset/merge` · YAML 第 3995 行
- 上传并合并 · `POST /upload/api/v1/asset/merge` · YAML 第 4043 行
- 移除版本 · `POST /api/v1/asset/separate` · YAML 第 4093 行
- 调整版本排序 · `POST /api/v1/asset/reindexVersions` · YAML 第 4141 行
- 获取下载id · `POST /api/v1/asset/getDownloadId` · YAML 第 4187 行
- 自动下载检测 · `POST /api/v1/asset/autoDownloadCheck` · YAML 第 4247 行
- 下载 · `GET /api/v1/asset/download` · YAML 第 4295 行
- 单个文件获取下载id · `POST /api/v1/asset/getDownloadIdForSingle` · YAML 第 4328 行
- 单个文件下载 · `GET /api/v1/asset/downloadForSingle` · YAML 第 4377 行
- 通用流播放 · `GET /api/v1/asset/stream` · YAML 第 4411 行
- 清理文件的代理文件 · `POST /api/v1/asset/cleanProxyFile` · YAML 第 4451 行
- 重试失败的预处理 · `POST /api/v1/asset/retryPrepareFailed` · YAML 第 4500 行
- 修改状态描述评分 · `POST /api/v1/asset/updateAttributes` · YAML 第 4547 行
- 修改标签 · `POST /api/v1/asset/updateTags` · YAML 第 4600 行
- 批量设置素材自定义字段 · `POST /api/v1/asset/updateCustomField` · YAML 第 4654 行
- 解压eaglepack文件 · `POST /api/v1/asset/unpackEagleFile` · YAML 第 4727 行
- 按真实格式解压素材压缩包 · `POST /api/v1/asset/unpackZipFile` · YAML 第 4770 行
- 素材精简列表 · `GET /api/v1/asset/liteList` · YAML 第 10411 行
- 批量素材详情 · `POST /api/v1/asset/multiDetail` · YAML 第 10438 行
- 全局标准搜索 · `POST /api/v1/standard-search/global-search` · YAML 第 11483 行
- 标准搜索 · `POST /api/v1/standard-search/search` · YAML 第 11516 行
- 搜索定位 · `POST /api/v1/standard-search/position` · YAML 第 11553 行
- 搜索邻近素材 · `POST /api/v1/standard-search/neighbors` · YAML 第 11586 行
- 文件夹搜索 · `POST /api/v1/standard-search/folder-search` · YAML 第 11619 行
- 搜索聚合 · `POST /api/v1/standard-search/facets` · YAML 第 11660 行
- 获取搜索足迹 · `GET /api/v1/user-track/getSearchTrack` · YAML 第 12061 行
- 添加素材访问记录 · `POST /api/v1/user-track/addAssetVisit` · YAML 第 12088 行
- 删除用户足迹 · `POST /api/v1/user-track/delete` · YAML 第 12121 行
- 按类型删除用户足迹 · `POST /api/v1/user-track/deleteByType` · YAML 第 12154 行

## 任务

- 任务列表 · `GET /api/v1/task/list` · YAML 第 4835 行
- 任务列表 by ids · `POST /api/v1/task/listByTaskIds` · YAML 第 5048 行
- 任务详情 · `GET /api/v1/task/detail` · YAML 第 5089 行
- 删除任务 · `POST /api/v1/task/delete` · YAML 第 5186 行
- 取消任务 · `POST /api/v1/task/cancel` · YAML 第 5264 行
- 重试任务 · `POST /api/v1/task/retry` · YAML 第 5342 行
- 清理已完成的任务 · `POST /api/v1/task/cleanDone` · YAML 第 5420 行

## 回收站

- 列表 · `GET /api/v1/trash/list` · YAML 第 5498 行
- 移入 · `POST /api/v1/trash/moveIn` · YAML 第 5582 行
- 移出 · `POST /api/v1/trash/moveOut` · YAML 第 5629 行
- 彻底删除 · `POST /api/v1/trash/clean` · YAML 第 5680 行

## 系统

- 应用初始化 · `POST /api/v1/sys/init` · YAML 第 5735 行
- 获取系统配置 · `GET /api/v1/sys/config` · YAML 第 5782 行
- 获取授权路径 · `GET /api/v1/sys/authPaths` · YAML 第 5925 行
- 获取授权目录的子项 · `POST /api/v1/sys/authPathEntries` · YAML 第 6011 行
- 新建目录（基于授权目录） · `POST /api/v1/sys/makeAuthPathDir` · YAML 第 7230 行
- 获取转码信息 · `GET /api/v1/sys/getTranscodeInfo` · YAML 第 7273 行
- 修改转码信息 · `POST /api/v1/sys/updateTranscodeInfo` · YAML 第 7370 行
- 获取系统卷列表 · `GET /api/v1/sys/vols` · YAML 第 10504 行
- 获取系统用户信息 · `POST /api/v1/sys/sysUser` · YAML 第 10531 行
- 查询试用权益 · `GET /api/v1/benefits/trialAvailable` · YAML 第 12187 行
- 申请试用授权 · `POST /api/v1/benefits/applyTrialLicense` · YAML 第 12214 行
- 更新品牌设置 · `POST /api/v1/brand/update` · YAML 第 12247 行
- 获取品牌 Logo · `GET /api/v1/brand/logo/{logoGuid}` · YAML 第 12280 行
- 上传品牌 Logo · `POST /upload/api/v1/brand/logo` · YAML 第 12314 行

## 评论

- 查看评论列表 · `GET /api/v1/comment/list` · YAML 第 7420 行
- 添加评论 · `POST /api/v1/comment/create` · YAML 第 7745 行
- 删除评论 · `POST /api/v1/comment/delete` · YAML 第 7809 行
- 修改评论 · `POST /api/v1/comment/update` · YAML 第 7847 行
- 修改评论status · `POST /api/v1/comment/update/status` · YAML 第 7999 行
- 添加评论 reaction · `POST /api/v1/comment/reaction/add` · YAML 第 8041 行
- 移除评论 reaction · `POST /api/v1/comment/reaction/remove` · YAML 第 8082 行
- 附件上传 · `POST /upload/api/v1/comment/attachment` · YAML 第 8123 行
- 附件删除 · `POST /api/v1/comment/attachment/delete` · YAML 第 8172 行
- 附件预览 · `GET /api/v1/comment/attachment/preview` · YAML 第 8212 行
- 附件下载 · `GET /api/v1/comment/attachment/download` · YAML 第 8241 行
- 附件缩略图预览 · `GET /api/v1/comment/attachment/thumbnail` · YAML 第 8281 行
- 评论设为已读 · `POST /api/v1/comment/makeItRead` · YAML 第 8321 行

## ws

- websocket建立 · `GET /api/v1/ws` · YAML 第 8361 行
- 查询ws连接信息(开发用) · `GET /api/v1/ws/stats` · YAML 第 8388 行

## 外链

- 获取外链列表 · `GET /api/v1/sharelink/list` · YAML 第 8426 行
- 生成分享链接的微信小程序码 · `POST /api/v1/sharelink/miniAppCode` · YAML 第 8657 行
- 获取外链信息 · `GET /api/v1/sharelink/info` · YAML 第 8695 行
- 创建外链 · `POST /api/v1/sharelink/create` · YAML 第 8796 行
- 更新外链 · `POST /api/v1/sharelink/update` · YAML 第 8955 行
- 删除外链 · `POST /api/v1/sharelink/delete` · YAML 第 9017 行
- 外链新增 asset · `POST /api/v1/sharelink/add-assets` · YAML 第 9057 行
- 外链移除 asset · `POST /api/v1/sharelink/remove-assets` · YAML 第 9100 行
- 获取外链 assets · `GET /api/v1/sharelink/asset/list` · YAML 第 9143 行
- 获取导航栏信息 · `GET /api/v1/sharelink/getNavigationPath` · YAML 第 9242 行
- 升级外链资产水印 · `POST /api/v1/sharelink/upgradeAssetWatermark` · YAML 第 10826 行
- 重试外链任务 · `POST /api/v1/sharelink/retryTask` · YAML 第 10859 行
- 外链素材精简列表 · `GET /api/v1/sharelink/liteList` · YAML 第 10892 行
- 外链智能文件夹素材列表 · `POST /api/v1/sharelink/collectionAsset/list` · YAML 第 10919 行

## 通知

- 获取通知列表 · `GET /api/v1/notice/list` · YAML 第 9341 行
- 设为已读 · `POST /api/v1/notice/makeItRead` · YAML 第 9478 行
- 全部已读 · `POST /api/v1/notice/makeAllRead` · YAML 第 9516 行
- 获取通知设置 · `GET /api/v1/notice/getSetting` · YAML 第 9543 行
- 修改通知设置 · `POST /api/v1/notice/updateSetting` · YAML 第 9646 行
- 获取未读数量 · `GET /api/v1/notice/unreadCount` · YAML 第 9703 行
- 清除通知 · `POST /api/v1/notice/cleanNotice` · YAML 第 9806 行

## 标签、状态管理

- 查看状态列表 · `GET /api/v1/status/list` · YAML 第 9850 行
- 新建状态 · `POST /api/v1/status/create` · YAML 第 9899 行
- 删除状态 · `POST /api/v1/status/delete` · YAML 第 9942 行
- 移动状态 · `POST /api/v1/status/sort` · YAML 第 9984 行
- 查看标签列表 · `GET /api/v1/tag/list` · YAML 第 10029 行
- 新建标签 · `POST /api/v1/tag/create` · YAML 第 10063 行
- 删除标签 · `POST /api/v1/tag/delete` · YAML 第 10103 行
- 标签管理列表 · `GET /api/v1/tag/manage/list` · YAML 第 10597 行
- 更新标签 · `POST /api/v1/tag/update` · YAML 第 10624 行
- 统计标签素材数量 · `POST /api/v1/tag/countEntries` · YAML 第 10657 行
- 合并标签 · `POST /api/v1/tag/merge` · YAML 第 10690 行
- 更新状态 · `POST /api/v1/status/update` · YAML 第 10723 行
- 更新状态前缀配置 · `POST /api/v1/status/updatePrefix` · YAML 第 10758 行
- 设置默认 LUT · `POST /api/v1/lut/setDefault` · YAML 第 10793 行

## 预览管理

- lut上传 · `POST /upload/api/v1/lut` · YAML 第 10143 行
- 查看lut · `GET /api/v1/lut/list` · YAML 第 10181 行
- 移动lut · `POST /api/v1/lut/sort` · YAML 第 10208 行
- 获取lut缩略图 · `GET /api/v1/lut/thumbnail/{lutGuid}` · YAML 第 10251 行
- 查看lut原图 · `GET /api/v1/lut/stream/{lutGuid}` · YAML 第 10285 行
- 删除lut · `POST /api/v1/lut/delete` · YAML 第 10319 行
- AI 设置状态 · `GET /api/v1/ai-setting/status` · YAML 第 10952 行
- 检查 AI 安装要求 · `GET /api/v1/ai-setting/check-requirement` · YAML 第 10979 行
- 设置 AI 数据卷 · `POST /api/v1/ai-setting/data-vol` · YAML 第 11006 行
- 设置 AI GPU · `POST /api/v1/ai-setting/gpu` · YAML 第 11039 行
- 安装 AI 模型 · `POST /api/v1/ai-setting/install` · YAML 第 11072 行
- 取消安装 AI 模型 · `POST /api/v1/ai-setting/cancel-install` · YAML 第 11105 行
- 卸载 AI 模型 · `POST /api/v1/ai-setting/uninstall` · YAML 第 11138 行
- 获取 AI 安装进度 · `GET /api/v1/ai-setting/get-install-process` · YAML 第 11171 行
- 获取 AI 搜索索引状态 · `GET /api/v1/ai-setting/search-index/status` · YAML 第 11198 行
- 设置 AI 搜索索引开关 · `POST /api/v1/ai-setting/search-index/set-enabled` · YAML 第 11225 行
- 设置 AI 搜索排除项目 · `POST /api/v1/ai-setting/search-index/set-exclude-projects` · YAML 第 11258 行
- 设置 AI 搜索索引策略 · `POST /api/v1/ai-setting/search-index/set-strategy` · YAML 第 11291 行
- 构建 AI 搜索索引 · `POST /api/v1/ai-setting/search-index/build` · YAML 第 11324 行
- AI 搜索准备 · `POST /api/v1/ai-search/prepare` · YAML 第 11357 行
- 上传 AI 搜索临时文件 · `POST /api/v1/ai-search/upload-tmp-file` · YAML 第 11390 行
- AI 搜索临时文件缩略图 · `GET /api/v1/ai-search/tmp-file-thumbnail/{tmpFileId}` · YAML 第 11417 行
- 执行 AI 搜索 · `POST /api/v1/ai-search/do-search` · YAML 第 11450 行
- 获取水印设置 · `GET /api/v1/watermark/getSetting` · YAML 第 11914 行
- 获取全部资产水印统计 · `GET /api/v1/watermark/getAllAssetStat` · YAML 第 11941 行
- 获取水印图片 · `GET /api/v1/watermark/image/{imageGuid}` · YAML 第 11968 行
- 上传水印图片 · `POST /api/v1/watermark/uploadImage` · YAML 第 12001 行
- 更新水印设置 · `POST /api/v1/watermark/updateSetting` · YAML 第 12028 行

## 维护工具

- 重建索引 · `POST /api/v1/maintenance/rebuildIndex` · YAML 第 11695 行
- 数据升级 · `POST /api/v1/maintenance/dataUpgrade` · YAML 第 12341 行

## 其他

- a · `POST /` · YAML 第 12368 行

## 云空间

- 更新云缓存设置 · `POST /api/v1/cloudSpace/updateSetting` · YAML 第 12401 行
- 查询云空间详情 · `GET /api/v1/cloudSpace/detail` · YAML 第 12443 行
- 查询云空间文件列表 · `GET /api/v1/cloudSpace/list` · YAML 第 12504 行
- 批量删除云空间文件 · `POST /api/v1/cloudSpace/deleteFiles` · YAML 第 12593 行
