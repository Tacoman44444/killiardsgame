package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type HubMessageType int

const (
	HubSendPlayerToLobby HubMessageType = iota
	HubPlayerInvalidCode
	HubRoomCreated
)

type HubMessage struct {
	msgType HubMessageType
	code    string
	player  *Player
	lobby   *Lobby
}

type Hub struct {
	readPlayer chan PlayerMessage
	readLobby  chan LobbyMessage
}

func NewHub() *Hub {
	return &Hub{make(chan PlayerMessage), make(chan LobbyMessage)}
}

func (h *Hub) Run() {
	lobbies := make(map[string]*Lobby)

	for {
		select {
		case plrmsg := <-h.readPlayer:
			if plrmsg.msgType == PlayerJoinRoom {
				lobby, ok := lobbies[plrmsg.msg.Code]
				if ok {
					hubmsg := HubMessage{
						msgType: HubSendPlayerToLobby,
						code:    plrmsg.msg.Code,
						player:  plrmsg.player,
						lobby:   lobby,
					}
					lobby.readHub <- hubmsg
					plrmsg.player.readHub <- hubmsg
				}
			} else if plrmsg.msgType == PlayerCreateRoom {

				newCode := RandomUppercaseString6()
				lobbies[plrmsg.msg.Code] = NewLobby(newCode, plrmsg.player)
				hubmsg := HubMessage{
					msgType: HubRoomCreated,
					code:    newCode,
					player:  plrmsg.player,
					lobby:   lobbies[plrmsg.msg.Code],
				}
				plrmsg.player.readHub <- hubmsg
			}
		case lbmsg := <-h.readLobby:
			if lbmsg.msgType == LobbyClose {
				lobby, ok := lobbies[lbmsg.lobbyCode]
				if !ok {
					fmt.Println("lobby does not exist, wtf?")

				}
				close(lobby.readHub)
				close(lobby.Inbound)
				delete(lobbies, lbmsg.lobbyCode)
			}
		}
	}
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("hub: conn error")
		return
	}

	player := GetNewPlayer(conn, h)

	go player.Run()
}

func RandomUppercaseString6() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
