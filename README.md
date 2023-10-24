# WebSocket Proxy

A simple websocket proxy server.

## Usage

```bash
go run main.go -backend ws://example.com:3000 -addr :7788
```

- `backend`：要代理的目标地址

- `addr`：代理程序监听的端口号（注意别漏了 `:`）

现在，所有传入此服务器的 WebSocket 请求都将被代理到 `ws://example.com:3000`

## 奇怪的报错

```bash
Error connecting to target WebSocket: read tcp 127.0.0.1:51510->127.0.0.1:7890: read: connection reset by peer
```

解决：重启终端就好了，原因不清楚...

