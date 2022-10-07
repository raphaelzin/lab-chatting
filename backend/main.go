package main

import (
	"main/routers"
	"net/http"
)

func setupRouters() {
	http.ListenAndServe(":8082", routers.Init())
}

func main() {
	setupRouters()
}
