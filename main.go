package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketProxy struct {
	targetAddr string // 目标 WebSocket 服务器的地址
	upgrader   *websocket.Upgrader
}

func NewProxy(targetAddr string) *WebsocketProxy {
	return &WebsocketProxy{
		targetAddr: targetAddr,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}}
}

func (wp *WebsocketProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 请求升级为 WebSocket 连接
	conn, err := wp.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// 连接到目标 WebSocket 服务器
	targetConn, _, err := websocket.DefaultDialer.Dial(wp.targetAddr, nil)
	if err != nil {
		log.Println("Error connecting to target WebSocket:", err)
		return
	}
	defer targetConn.Close()

	log.Println("Start to copy messages")
	// 在代理之前将消息从客户端复制到目标服务器，以及从目标服务器复制到客户端
	go copyMessages(conn, targetConn)
	copyMessages(targetConn, conn)
}

func copyMessages(dst, src *websocket.Conn) {
	for {
		messageType, p, err := src.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		err = dst.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

// usage
var (
	flagBackend = flag.String("backend", "ws://198.18.32.19:23333", "Backend URL for proxying")
	flagAddr    = flag.String("addr", ":7788", "Proxy server listen address")
)

func main() {
	flag.Parse()
	proxy := NewProxy(*flagBackend)
	log.Println("Starting websocket proxy server on", *flagAddr)
	log.Fatal(http.ListenAndServe(*flagAddr, proxy))
}
