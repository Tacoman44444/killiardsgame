package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type PlayerMessageType int

const (
	PlayerSendAction PlayerMessageType = iota
	PlayerSendWall
	PlayerStartGame
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

type PlayerInHub struct {
}

func (p *PlayerInHub) Enter(player *Player) {}

func (p *PlayerInHub) HandleClientMessage(cm ClientMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	switch cm.Type {
	case ClientSendId:
		player.id = cm.Id
	case ClientCreateRoom:
		if player.id != "" {
			msg := PlayerMessage{
				msgType:  PlayerCreateRoom,
				player:   player,
				sender:   player.conn,
				senderID: player.id,
				msg:      cm,
			}

			player.hub.readPlayer <- msg
		}
	case ClientJoinRoom:
		if player.id != "" {
			msg := PlayerMessage{
				msgType:  PlayerJoinRoom,
				player:   player,
				sender:   player.conn,
				senderID: player.id,
				msg:      cm,
			}

			player.hub.readPlayer <- msg
		}
	}
}

func (p *PlayerInHub) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	//lobby should not send messages to the player in this state
}

func (p *PlayerInHub) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	//hub does not need to send messages to the player in this state, except maybe a timeout message
}

func (p *PlayerInHub) Exit() {}

type PlayerRequestedForLobby struct {
	hub *Hub
}

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
		player.SetState(&PlayerInLobby{})
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
	}
}

func (p *PlayerInLobby) HandleLobbyMessage(lm LobbyMessage, channelOpen bool, player *Player) {
	if !channelOpen {
		//
	}

	switch lm.msgType {
	case LobbySendGameStart:
		msg := newGameStartMessage(lm.mapState, lm.player, lm.allPlayers)
		player.WriteToClient(msg, player.id)
		player.SetState(&PlayerInGame{})
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
			msgType: PlayerSendAction,
			msg:     cm,
		}

		player.lobby.Inbound <- playerMsg
	case ClientSendWall:
		playerMsg := PlayerMessage{
			msgType: PlayerSendWall,
			msg:     cm,
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
		serverMsg := newMapUpdateMessage(lm.mapState)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendTurnStart:
		serverMsg := newTurnStartMessage(lm.player.id)
		player.WriteToClient(serverMsg, player.id)
	case LobbySendTurnTimeout:
		serverMsg := newTurnTimeoutMessage()
		player.WriteToClient(serverMsg, player.id)
	}
}

func (p *PlayerInGame) HandleHubMessage(hm HubMessage, channelOpen bool, player *Player) {}

func (p *PlayerInGame) Exit() {}

type Player struct {
	id           string
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
	p.state.Exit()
	p.state = newState
}

func GetNewPlayer(conn *websocket.Conn, hub *Hub) *Player {
	player := &Player{
		id:           "",
		conn:         conn,
		socketClosed: false,
		clientMsg:    make(chan ClientMessage),
		readLobby:    make(chan LobbyMessage),
		readHub:      make(chan HubMessage),
		hub:          hub,
	}

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
