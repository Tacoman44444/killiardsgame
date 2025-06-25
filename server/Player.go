package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Tacoman44444/killiardsgame/server/tools"
	"github.com/gorilla/websocket"
)

type PlayerMessageType int

const (
	PlayerSendAction PlayerMessageType = iota
	PlayerSendWall
)

type PlayerMessage struct {
	msgType PlayerMessageType
	msg     ClientMessage
}

type Player struct {
	id           string
	conn         *websocket.Conn
	socketClosed bool
	state        PlayerState
	clientMsg    chan ClientMessage
	readLobby    chan LobbyMessage //lobby will write into this
	lobby        *Lobby
}

func GetNewPlayer(id string, conn *websocket.Conn, playerPosition tools.Vector2) *Player {
	player := &Player{
		id:           id,
		conn:         conn,
		socketClosed: false,
		state:        PlayerState{playerPosition.X, playerPosition.Y, 0, 0},
		clientMsg:    make(chan ClientMessage),
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

		p.lobby.Inbound <- playerMsg
	case ClientSendWall:
		playerMsg := PlayerMessage{
			msgType: PlayerSendWall,
			msg:     cm,
		}

		p.lobby.Inbound <- playerMsg
	}
}

func (p *Player) HandleLobbyMessage(lm LobbyMessage, channelOpen bool) {
	if !channelOpen {
		//that channel that recieves messages from the lobby has been closed, handle this
	}

	switch lm.msgType {
	case LobbySendEntityUpdate:
		serverMsg := newEntityUpdateMessage(lm.player, lm.otherPlayers, lm.walls)
		p.WriteToClient(serverMsg, p.id)
	case LobbySendWallUpdate:
		serverMsg := newWallUpdateMessage(lm.walls)
		p.WriteToClient(serverMsg, p.id)
	case LobbySendMapUpdate:
		serverMsg := newMapUpdateMessage(lm.mapState)
		p.WriteToClient(serverMsg, p.id)
	case LobbySendTurnStart:
		serverMsg := newTurnStartMessage()
		p.WriteToClient(serverMsg, p.id)
	case LobbySendTurnTimeout:
		serverMsg := newTurnTimeoutMessage()
		p.WriteToClient(serverMsg, p.id)
	}
}
