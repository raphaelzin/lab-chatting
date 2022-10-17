package chatting

import (
	"log"
	"main/chatting/chattingredis"
)

type Hub struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	redisBroadcast chan []byte
	register       chan *Client
	unregister     chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		redisBroadcast: make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	go h.subscribe()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.redisBroadcast:
			log.Println("New message: " + string(message))

			// Adds to stream for replay
			chattingredis.AddMessage(message)

			// Send to pub sub
			chattingredis.PublishMessage(message)
		}
	}
}

func (h *Hub) subscribe() {
	pubsub := chattingredis.GetSubscriptionToChannel("message-channel")
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			log.Println(err)
		}
		h.broadcast <- []byte(msg.Payload)
	}
}
