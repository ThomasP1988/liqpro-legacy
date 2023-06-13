package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	config "liqpro/config"
	"liqpro/shared/libs/order"
	"log"
	"strconv"
	"time"

	"github.com/fasthttp/websocket"
)

var OrderPayloadBeginning []byte = []byte("{\"event\":\"order\"")
var EventBeginning []byte = []byte("{\"event\"")

func (c *Client) readPump() {
	fmt.Println("test")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		HandleMessagesIn(c, message)
	}
}

// HandleMessagesIn handler socket messages in
func HandleMessagesIn(c *Client, message []byte) {
	// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	// doesn't look the best but we need to be fast on order
	if bytes.HasPrefix(message, OrderPayloadBeginning) {
		requestArgs := &order.RequestArgs{}
		err := json.Unmarshal(message, requestArgs)
		if err != nil {
			fmt.Println("error unmarshalling order msg", err)
		}
		bson, _ := order.Handle(requestArgs, c.userDataCache)
		(*c).Send([]byte(bson))
	} else if bytes.HasPrefix(message, EventBeginning) {
		msgIn := &MessageIn{}
		err := json.Unmarshal(message, msgIn)
		if err != nil {
			fmt.Println("error unmarshalling subscribe msg", err)
		}

		if instr, ok := config.AuthorisedInstrumentsAndLevels[msgIn.Instrument]; ok {
			for i := 0; i < len(msgIn.Levels); i++ {
				if _, ok := instr[msgIn.Levels[i]]; ok {
					sub := Subscribe{
						client:     c,
						instrument: msgIn.Instrument,
						level:      msgIn.Levels[i],
					}
					if msgIn.Event == "subscribe" {
						c.hub.subscribe <- sub
					} else if msgIn.Event == "unsubscribe" {
						c.hub.unsubscribe <- sub
					}
				} else {
					fmt.Println("error")
					(*c).Send([]byte("This level : " + strconv.FormatFloat(float64(msgIn.Levels[i]), 'f', 6, 32) + ", is not available for this instrument. (" + msgIn.Instrument + ")"))
				}
			}
		} else {
			fmt.Println("error")
			(*c).Send([]byte("This instrument is not available. (" + msgIn.Instrument + ")"))
		}
	}
}
