package chatting

import (
	"encoding/json"
	"log"
	"main/chatting/models"
	chattingModels "main/chatting/models"
	mainModels "main/models"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	user mainModels.User
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

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

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			var arr []string = make([]string, 0)
			arr = append(arr, string(message[:]))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				arr = append(arr, string(<-c.send))
			}

			data, _ := json.Marshal(arr)
			w.Write(data)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Unregister()
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
		rawMessage, err := models.Parse[models.RawMessage](message)

		if err != nil {
			c.send <- []byte(err.Error())
			continue
		}

		switch rawMessage.Type {
		case models.Text:
			var textMessage models.TextMessage = *models.NewTextMessage(rawMessage)
			textMessage.User = c.user
			data, _ := json.Marshal(textMessage)
			c.hub.textMessageBroadcast <- data
		}
	}
}

func (c *Client) Register() {
	c.hub.infoStream <- chattingModels.NewLoginMessage(c.user).AsData()
	c.hub.register <- c
}

func (c *Client) Unregister() {
	c.hub.unregister <- c
	c.hub.infoStream <- chattingModels.NewLogoutMessage(c.user).AsData()
	c.conn.Close()
}
