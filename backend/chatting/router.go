package chatting

import (
	"errors"
	"log"
	"main/chatting/chattingredis"
	"main/models"
	"net/http"

	redisInstance "main/redis"
	userManager "main/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func Init(r *mux.Router) {
	redisInstance.Init()
	setupChattingSocket(r)
}

func setupChattingSocket(r *mux.Router) {
	hub := newHub()
	go hub.run()
	r.HandleFunc("/chat/ws", func(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Unauthorized user, bye bye!")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("A user with this token could not be found"))
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(user.Username + " connected!")

	client := &Client{user: user, hub: hub, conn: ws, send: make(chan []byte, 256)}
	client.Register()

	messages, _ := chattingredis.GetLastN(10)
	for _, message := range messages {
		client.send <- message
	}

	go client.writePump()
	go client.readPump()
}

func validateUser(r *http.Request) (user models.User, err error) {
	if header := r.Header["Token"]; len(header) != 0 {
		return userManager.GetUserWithToken(header[0])
	}
	cokie, err := r.Cookie("Token")
	if err != nil {
		return user, errors.New("no user provided")
	}
	return userManager.GetUserWithToken(cokie.Value)
}
