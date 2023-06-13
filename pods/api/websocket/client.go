package main

import (
	"liqpro/shared/repositories/cache"
	"log"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// save all subscription
	subscriptions []Subscribe

	sendMutex *sync.Mutex

	userDataCache *cache.UserDataCache
}

// HandleNewConnection handle new connection, register client
func HandleNewConnection(ctx *fasthttp.RequestCtx, hub *Hub, userDataCache *cache.UserDataCache) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256),
			subscriptions: []Subscribe{}, sendMutex: &sync.Mutex{}, userDataCache: userDataCache}
		client.hub.register <- client

		go client.writePump()
		client.readPump()
	})

	if err != nil {
		log.Println(err)
	}
}
