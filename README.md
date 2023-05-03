# via-chat

前后端分离的多人在线聊天室。
基于 [go-gin-chat](https://github.com/hezhizheng/go-gin-chat) 进行修改。

- 前后端通过 http 通信，使用 session 记录状态实现登录验证
- 启用 websocket 实现双向通信及消息实时更新
- 使用 goroutine 处理用户连接、离线、消息发送等各种事件
- 通过 channel 实现并发情况下消息的处理
- 支持多房间聊天，支持文字和图片上传，支持加载历史消息
- ~~实现心跳检测机制，通过定时任务清理没有心跳的连接~~
- 接入 ChatGPT，用户可以与 GPT 进行单轮互动
- 实现前端和后端的 docker 部署，通过 docker-compose 自动编排



## Quick Start

- 自行安装docker

- 进入项目根目录

- 根据个人情况配置server的数据库、GPT

  ```
  cp server/configs/config.go.env server/configs/config.go
  cp server/configs/openai_config.ini.env server/configs/openai_config.ini
  ```

- docker compose 编译和启动（如需更改docker-compose 配置，请查看 `docker-compose.yml`

  ```
  docker-compose build
  docker-compose up -d
  
  # 停止
  # docker-compose stop
  # 移除
  # docker-compose down
  ```



## TODO
- [x] 数据库事务
- [ ] 引入Redis（在线用户列表、缓存聊天消息等）
- [ ] websocket HTTPS
- [x] certificates （access GPT）
- [x] JWT 替换 session （以去除同源访问）
- [ ] 心跳检测
- [x] bcrypt 替换 md5
- [x] Gorm add 重构（map to model）
- [ ] 高并发测试
- [ ] 撤回消息
- [ ] 私聊
- [ ] 分布式部署
  - [ ] nginx
  - [x] 远程数据库
  - [ ] GRPC
- [ ] kafka消息队列
- [ ] 音频、图片、文件等多模态
- [ ] langchain
