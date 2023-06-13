package binance

import (
	// "encoding/json"

	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
)

var (
	baseWsURL = "wss://stream.binance.com:9443"
	err       error
)

// WsDepthServe connect to binance stream
func WsDepthServe(symbols []string, handler WsDepthHandler, onConnected *func()) (doneC, stopC chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	subs := []string{}
	for _, symbol := range symbols {
		subs = append(subs, strings.ToLower(symbol)+"@depth20@100ms")
	}
	url := baseWsURL + "/stream?streams=" + strings.Join(subs, "/")

	c, _, err := websocket.DefaultDialer.Dial(url, nil)

	msg := WsSubscribeMessage{
		Method: "SUBSCRIBE",
		Params: subs,
		ID:     1,
	}

	log.Println("msg:", msg)
	byteJSONMsg, _ := json.Marshal(msg)
	errSub := c.WriteMessage(websocket.TextMessage, byteJSONMsg)

	if errSub != nil {
		log.Println("error subscription:", errSub)
		return
	}

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
	(*onConnected)()
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
