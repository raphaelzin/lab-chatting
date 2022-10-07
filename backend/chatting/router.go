package chatting

import (
	"errors"
	"log"
	"main/chatting/chattingredis"
	"main/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func Init(r *mux.Router) {
	chattingredis.Init()
	setupChattingSocket(r)
}

func setupChattingSocket(r *mux.Router) {
	hub := newHub()
	go hub.run()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(w, r, hub)
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request, hub *Hub) {
	user, err := validateUser(r)
	if err != nil {
		log.Println(r.Cookie("Username"))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing user"))
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(user.Username + " connected!")

	client := &Client{user: user, hub: hub, conn: ws, send: make(chan []byte, 256)}
	client.hub.register <- client

	messages, _ := chattingredis.GetLastN(10)
	for _, message := range messages {
		log.Println("Message: " + string(message))
		client.send <- message
	}

	go client.writePump()
	go client.readPump()
}

func validateUser(r *http.Request) (user models.User, err error) {
	if header := r.Header["Username"]; len(header) != 0 {
		return models.User{Username: header[0]}, nil
	}

	cokie, err := r.Cookie("Username")
	if err != nil {
		return user, errors.New("no user provided")
	}

	return models.User{Username: cokie.Value}, nil
}
