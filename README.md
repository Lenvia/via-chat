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



## TODO
- [ ] 数据库事务
- [ ] 图片上传图床
- [ ] 引入Redis（在线用户列表、缓存聊天消息等）
- [ ] grpc 远程调用（用在langchain）
- [ ] websocket HTTPS
- [x] certificates （access GPTt）
- [ ] JWT 替换 session （以去除同源访问）
- [ ] 心跳检测
