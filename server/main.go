package main

import (
	"log"
	"net/http"

	"github.com/Tacoman44444/killiardsgame/server/tools"
	"github.com/gorilla/websocket"
)

var up = websocket.Upgrader{ // reuse a single upgrader
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err := ws.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	/*
		http.HandleFunc("/ws", wsHandler)
		log.Println("Listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	*/

	arena := tools.GenerateMap(100, 100, false, "")
	arena.DebugLog()
}
