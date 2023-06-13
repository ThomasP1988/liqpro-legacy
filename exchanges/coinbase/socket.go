package coinbase

import (
	// "encoding/json"

	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fasthttp/websocket"
)

var (
	baseWsURL = "wss://ws-feed.pro.coinbase.com"
)

// WsDepthServe connect to bitstamp stream
func WsDepthServe(symbols []string, handler WsDepthHandler) (doneC, stopC chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(baseWsURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		handler(c)
		defer close(done)
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var channels []interface{}
	var productIds []string
	for _, symbol := range symbols {
		productIds = append(productIds, symbol)
	}

	channels = append(channels, "level2")
	channels = append(channels, "heartbeat")

	msg := WsSubscribeMessage{
		Type:       "subscribe",
		ProductIds: productIds,
		Channels:   channels,
	}

	log.Println("msg:", msg)
	byteJSONMsg, _ := json.Marshal(msg)
	errSub := c.WriteMessage(websocket.TextMessage, byteJSONMsg)

	if errSub != nil {
		log.Println("error subscription:", errSub)
		return
	}

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
