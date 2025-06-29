package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("here we fucking go")
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", hub.ServeWs)

	http.ListenAndServe("localhost:8000", nil)
}
