package ws

import (
	"github.com/gorilla/websocket"
	"log/slog"
)

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}


// ReadPump listens for messages from the client (we don't use these for now,
// but we need to read to detect disconnection)
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			slog.Info("client disconnected")
			break
		}
	}
}

// WritePump sends messages from the hub to this client
func (c *Client) WritePump() {
	defer c.conn.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			slog.Error("write error", "error", err)
			break
		}
	}
}