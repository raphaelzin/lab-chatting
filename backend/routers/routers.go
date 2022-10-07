package routers

import (
	"main/chatting"

	"github.com/gorilla/mux"
)

func Init() *mux.Router {
	route := mux.NewRouter()
	chatting.Init(route)
	return route
}
