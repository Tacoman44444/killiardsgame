package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerMessageType int

const (
	PlayerSendAction PlayerMessageType = iota
	PlayerSendWall
	PlayerStartGame
	PlayerLeaveRoom
	PlayerSimulationDone
	PlayerJoinRoom
	PlayerCreateRoom
)

type PlayerMessage struct {
	msgType  PlayerMessageType
	player   *Player
	sender   *websocket.Conn
	senderID string
	msg      ClientMessage
}

/*
POSSIBLE PLAYER STATES:
1. PLAYER CONNECTED BUT NOT HAS NOT ENTERED ANY CODE (OR HAS JUST EXITED A LOBBY) -- IN HUB BUT NOT IN LOBBY
2. PLAYER HAS REQUESTED FOR A LOBBY VIA A CODE
3. PLAYER HAS BEEN GRANTED THE LOBBY, PLAYER IS NOW IN LOBBY BUT NOT IN GAME
4. GAME HAS BEEN STARTED BY LOBBY OWNER, PLAYER IS NOW INGAME
*/

type PlayerState interface {
	Enter(player *Player)
	HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player)
	HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player)
	HandleHubMessage(hm HubMessage, channelOpen bool, player *Player)
	Exit()
}

type PlayerInHub struct{}

func (p *PlayerInHub) Enter(player *Player) {}

func (p *PlayerInHub) HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		// the player <- client channel has been closed.
	}

	switch cm.Type {
	case ClientCreateRoom:
		msg := PlayerMessage{
			msgType:  PlayerCreateRoom,
			player:   player,
			sender:   player.conn,
			senderID: player.id,
			msg:      cm,
		}
		player.username = cm.Username
		player.hub.readPlayer <- msg
		player.SetState(&PlayerRequestedForLobby{})
	case ClientJoinRoom:
		msg := PlayerMessage{
			msgType:  PlayerJoinRoom,
			player:   player,
			sender:   player.conn,
			senderID: player.id,
			msg:      cm,
		}
		player.username = cm.Username
		player.hub.readPlayer <- msg
		player.SetState(&PlayerRequestedForLobby{})
	}
}

func (p *PlayerInHub) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {}

func (p *PlayerInHub) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {}

func (p *PlayerInHub) Exit() {}

type PlayerRequestedForLobby struct{}

func (p *PlayerRequestedForLobby) Enter(player *Player) {}

func (p *PlayerRequestedForLobby) HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player) {
}

func (p *PlayerRequestedForLobby) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {
}

func (p *PlayerRequestedForLobby) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	switch hm.msgType {
	case HubSendPlayerToLobby:
		player.WriteToClient(newRoomJoinedMessage(), player.id)
		player.SetState(&PlayerInLobby{l: hm.lobby})
	case HubRoomCreated:
		fmt.Println("writeing to player rn that room has been created")
		player.WriteToClient(newRoomCreatedMessage(hm.code), player.id)
		player.SetState(&PlayerInLobby{l: hm.lobby})
	case HubPlayerInvalidCode:
		player.WriteToClient(newInvalidCodeMessage(), player.id)
		player.SetState(&PlayerInHub{})
	}
}

func (p *PlayerRequestedForLobby) Exit() {}

type PlayerInLobby struct {
	l *Lobby
}

func (p *PlayerInLobby) Enter(player *Player) {
	player.lobby = p.l
}

func (p *PlayerInLobby) HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	switch cm.Type {
	case ClientStartGame:
		msg := PlayerMessage{
			msgType:  PlayerStartGame,
			sender:   player.conn,
			senderID: player.id,
			msg:      cm,
		}
		player.lobby.Inbound <- msg
	case ClientLeaveRoom:
		msg := PlayerMessage{
			msgType:  PlayerLeaveRoom,
			player:   player,
			sender:   player.conn,
			senderID: player.id,
			msg:      cm,
		}
		player.lobby.Inbound <- msg
		player.SetState(&PlayerInHub{})
	}
}

func (p *PlayerInLobby) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}
	switch lm.msgType {
	case LobbySendGameStart:
		fmt.Println("[PLAYER IN LOBBY] sending player a game-start message")
		otherPlayers := make([]PlayerIdentity, 0)
		for i := range lm.allPlayers {
			if lm.allPlayers[i].id != lm.player.id {
				otherPlayers = append(otherPlayers, lm.allPlayers[i])
			}
		}
		msg := newGameStartMessage(lm.currentMap, lm.nextMap, lm.player, otherPlayers)
		player.WriteToClient(msg, player.id)
		player.SetState(&PlayerInGame{})
	case LobbySendMakeOwner:
		fmt.Println("[PLAYER IN LOBBY] sending player a make owner message")
		msg := newMakeOwnerMessage()
		player.WriteToClient(msg, player.id)
	}
}

func (p *PlayerInLobby) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {}

func (p *PlayerInLobby) Exit() {}

type PlayerInGame struct{}

func (p *PlayerInGame) Enter(player *Player) {}

func (p *PlayerInGame) HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//the channel that recieves messages from the client has been closed, handle this accordingly
	}

	switch cm.Type {
	case ClientCreateRoom:
		fmt.Println("cannot create room while player is already ingame")
	case ClientJoinRoom:
		fmt.Println("cannot join room while player is already ingame")
	case ClientStartGame:
		fmt.Println("cannot start game while player is already ingame")
	case ClientSendTurn:
		playerMsg := PlayerMessage{
			msgType:  PlayerSendAction,
			player:   player,
			sender:   player.conn,
			senderID: player.id,
			msg:      cm,
		}

		player.lobby.Inbound <- playerMsg
	case ClientSendWall:
		playerMsg := PlayerMessage{
			msgType: PlayerSendWall,
			msg:     cm,
		}

		player.lobby.Inbound <- playerMsg
	case ClientSimulationDone:
		playerMsg := PlayerMessage{
			msgType: PlayerSimulationDone,
		}
		player.lobby.Inbound <- playerMsg
	}
}

func (p *PlayerInGame) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//that channel that recieves messages from the lobby has been closed, handle this
	}

	switch lm.msgType {
	case LobbySendEntityUpdate:
		serverMsg := newEntityUpdateMessage(lm.player, lm.allPlayers, lm.walls)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendWallUpdate:
		serverMsg := newWallUpdateMessage(lm.walls)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendMapUpdate:
		serverMsg := newMapUpdateMessage(lm.currentMap, lm.nextMap)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendTurnStart:
		serverMsg := newTurnStartMessage(lm.player.id)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendTurnTimeout:
		serverMsg := newTurnTimeoutMessage()
		player.WriteToClient(serverMsg, player.id)
	case LobbyBroadcastMove:
		serverMsg := newBroadcastTurnMessage(lm.player, lm.action)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendEliminations:
		serverMsg := newEliminationMessage(lm.eliminatedPlayers)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendGameOver:
		serverMsg := newGameFinishedMessage(lm.result, lm.winnerName)
		player.WriteToClient(serverMsg, player.id)
	}
}

func (p *PlayerInGame) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {}

func (p *PlayerInGame) Exit() {}

type Player struct {
	id           string
	username     string
	turnsPlayed  int
	conn         *websocket.Conn
	socketClosed bool
	state        PlayerState
	clientMsg    chan ClientMessage
	readLobby    chan LobbyMessage //lobby will write into this
	readHub      chan HubMessage
	lobby        *Lobby
	hub          *Hub
}

func (p *Player) SetState(newState PlayerState) {
	if p.state != nil {
		p.state.Exit()
	}
	fmt.Printf("[STATE] Entering %T for player %s\n", newState, p.id)
	p.state = newState
	p.state.Enter(p)
}

func GetNewPlayer(conn *websocket.Conn, hub *Hub) *Player {
	player := &Player{
		id:           randomAlphanumericString(),
		turnsPlayed:  0,
		conn:         conn,
		socketClosed: false,
		clientMsg:    make(chan ClientMessage),
		readLobby:    make(chan LobbyMessage),
		readHub:      make(chan HubMessage),
		hub:          hub,
	}
	player.SetState(&PlayerInHub{})

	return player
}

func (p *Player) Run() {
	go p.ReadClientMessage()
	defer fmt.Println("player goroutine exited")

	for {
		select {
		case cm, ok := <-p.clientMsg:
			p.HandleClientMessage(cm, ok)
		case rm, ok := <-p.readLobby:
			p.HandleLobbyMessage(rm, ok)
		case hm, ok := <-p.readHub:
			p.HandleHubMessage(hm, ok)
		}
	}
}

func (p *Player) ReadClientMessage() {
	defer func() {
		p.conn.Close()
		close(p.clientMsg)
	}()

	for {
		wsMsgType, data, err := p.conn.ReadMessage()
		if err != nil {
			return
		}
		if wsMsgType != websocket.TextMessage {
			return
		}

		msg := ClientMessage{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			return
		}

		p.clientMsg <- msg
	}
}

func (p *Player) WriteToClient(msg ServerMessage, playerID string) error {
	if p.socketClosed {
		return errors.New("connection closed, player ID: " + playerID)
	}

	err := p.conn.WriteJSON(msg)
	if err != nil {
		p.socketClosed = true
		p.conn.Close()
		return err
	}
	return nil
}

func (p *Player) HandleClientMessage(cm ClientMessage, channelOpen bool) {
	p.state.HandleClientMessage(cm, channelOpen, p)
}

func (p *Player) HandleLobbyMessage(lm LobbyMessage, channelOpen bool) {
	p.state.HandleLobbyMessage(lm, channelOpen, p)
}

func (p *Player) HandleHubMessage(hm HubMessage, channelOpen bool) {
	p.state.HandleHubMessage(hm, channelOpen, p)
}

func randomAlphanumericString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 10)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
