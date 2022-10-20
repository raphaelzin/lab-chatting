package routers

import (
	"main/chatting"
	"main/user"

	"github.com/gorilla/mux"
)

func Init() *mux.Router {
	route := mux.NewRouter()
	chatting.Init(route)
	user.Init(route)
	return route
}
