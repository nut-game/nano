# nano-chat-demo

chat room demo base on [nano](https://github.com/topfreegames/nano) in 100 lines

refs: https://github.com/topfreegames/nano

## Required

- golang
- websocket

## Run

```bash
docker compose -f ../../testing/docker-compose.yml up -d etcd nats
go run main.go
```

open browser => http://localhost:3251/web/
