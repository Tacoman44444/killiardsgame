package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

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

type Lobby struct {
	Inbound chan PlayerMessage
}
