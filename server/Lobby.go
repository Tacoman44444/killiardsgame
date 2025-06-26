package main

import (
	"fmt"

	"github.com/Tacoman44444/killiardsgame/server/tools"
	"github.com/gorilla/websocket"
)

type LobbyMessageType int

const (
	LobbySendEntityUpdate LobbyMessageType = iota
	LobbySendWallUpdate
	LobbySendMapUpdate
	LobbySendTurnStart
	LobbySendTurnTimeout
)

type LobbyMessage struct {
	msgType      LobbyMessageType
	player       PlayerState
	otherPlayers []PlayerState
	walls        []WallState
	mapState     tools.MapState
}

type LobbyState interface {
	HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby)
	HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby)
}

type Lobby struct {
	Inbound   chan PlayerMessage
	state     LobbyState
	gameState *GameState
	players   map[*websocket.Conn]Player
	owner     *Player
}

func (l *Lobby) SetState(newState LobbyState) {
	l.state = newState
}

type LobbyWaitingForPlayers struct {
}

func (l *LobbyWaitingForPlayers) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle player channel not being open
	}

	switch pm.msgType {
	case PlayerStartGame:
		if pm.sender == lobby.owner.conn {
			//start game and send all clients a message that the game has begun
			lobby.SetState(&LobbyWaitingForTurn{})
			players := make([]string, 10)
			for _, value := range lobby.players {
				players = append(players, value.id)
			}
			lobby.gameState = GetNewGame(players)

		} else if pm.sender != lobby.owner.conn {
			fmt.Println("only the party owner can start the match")
			return
		}
	}
}

func (l *LobbyWaitingForPlayers) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {
	if !channelOpen {
		//handle player channel not being open
	}

	switch hm.msgType {
	case SendPlayerToLobby:

	}
}

type LobbyWaitingForTurn struct {
}

func (l *LobbyWaitingForTurn) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {

}

func (l *LobbyWaitingForTurn) HandleHubMessage(hm HubMessage, channelOpen bool, lobby *Lobby) {

}

type LobbyBetweenTurns struct {
}

func (l *LobbyBetweenTurns) HandlePlayerMessage(pm PlayerMessage, channelOpen bool, lobby *Lobby) {

}

func (l *LobbyBetweenTurns) HandleJoinRequest(hm HubMessage, channelOpen bool, lobby *Lobby) {

}
