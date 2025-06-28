package main

import (
	"github.com/Tacoman44444/killiardsgame/server/tools"
)

func main() {
	/*
		http.HandleFunc("/ws", wsHandler)
		log.Println("Listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/

	arena := tools.GenerateMap(100, 100, false, "")
	arena.DebugLog()
}
