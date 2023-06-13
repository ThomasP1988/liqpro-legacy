package socketredis

import (
	socket "github.com/fasthttp/websocket"
	"github.com/go-redis/redis/v8"
)

// Member to add to Redis set
type Member struct {
	Price    string
	Volume   string
	Platform string
}

// WsHandler function to be send to the websocket function
type WsHandler func(c *socket.Conn)

// WsHandlerBytes function to be send to the websocket function
type WsHandlerBytes func(msg *[]byte)

type ConsumerBase struct {
	Market       string
	Action       string
	Price        string
	Volume       string
	Platform     string
	Key          string
	Score        string
	Err          error
	Prefix       string
	RedisResult  []redis.Cmder
	RedisMemberZ redis.Z
}
