package main

import (
	"time"

	"github.com/fasthttp/websocket"
)

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send send message to client
func (c *Client) Send(message []byte) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	// if !ok {
	// 	// The hub closed the channel.
	// 	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 	return
	// }

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write(message)

	// Add queued chat messages to the current websocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		w.Write(newline)
		w.Write(<-c.send)
	}

	if err := w.Close(); err != nil {
		return
	}
}
