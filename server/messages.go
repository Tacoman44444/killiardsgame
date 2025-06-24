package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

type ServerMessageType int

const (
	ServerRoomCreated  ServerMessageType = 1
	ServerRoomJoined   ServerMessageType = 2
	ServerGameStart    ServerMessageType = 3
	ServerTurnStart    ServerMessageType = 4
	ServerTurnTimeout  ServerMessageType = 5
	ServerEntityUpdate ServerMessageType = 6
	ServerWallUpdate   ServerMessageType = 7
	ServerMapUpdate    ServerMessageType = 8
	ServerGameFinished ServerMessageType = 9
)

type ServerMessage interface {
	isServerMessage()
}

type RoomCreatedMessage struct {
	Type ServerMessageType `json:"type"`
	Code int               `json:"code"`
}

func (m RoomCreatedMessage) isServerMessage() {}

type RoomJoinedMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m RoomJoinedMessage) isServerMessage() {}

type GameStartMessage struct {
	Type         ServerMessageType `json:"type"`
	Map          tools.MapState    `json:"map"`
	Player       PlayerState       `json:"player"`
	OtherPlayers []PlayerState     `json:"other_players"`
}

func (m GameStartMessage) isServerMessage() {}

type TurnStartMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m TurnStartMessage) isServerMessage() {}

type TurnTimeoutMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m TurnTimeoutMessage) isServerMessage() {}

type EntityUpdateMessage struct {
	Type         ServerMessageType `json:"type"`
	Player       PlayerState       `json:"player_state"`
	OtherPlayers []PlayerState     `json:"other_players"`
	Walls        []WallState       `json:"walls"`
}

func (m EntityUpdateMessage) isServerMessage() {}

type WallUpdateMessage struct {
	Type  ServerMessageType `json:"type"`
	Walls []WallState       `json:"walls"`
}

func (m WallUpdateMessage) isServerMessage() {}

type MapUpdateMessage struct {
	Type ServerMessageType `json:"type"`
	Map  tools.MapState    `json:"map"`
}

func (m MapUpdateMessage) isServerMessage() {}

type GameFinishedMessage struct {
	Type       ServerMessageType `json:"type"`
	WinnerName string            `json:"winner_name"`
}

// CREATING NEW MESSAGES

func newRoomCreatedMessage(code int) RoomCreatedMessage {
	return RoomCreatedMessage{ServerRoomCreated, code}
}

func newRoomJoinedMessage() RoomJoinedMessage {
	return RoomJoinedMessage{ServerRoomJoined}
}

func newGameStartMessage(mapState tools.MapState, player PlayerState, otherPlayers []PlayerState) GameStartMessage {
	return GameStartMessage{ServerGameStart, mapState, player, otherPlayers}
}

func newTurnStartMessage() TurnStartMessage {
	return TurnStartMessage{ServerTurnStart}
}

func newTurnTimeoutMessage() TurnTimeoutMessage {
	return TurnTimeoutMessage{ServerTurnTimeout}
}

func newEntityUpdateMessage(player PlayerState, otherPlayers []PlayerState, walls []WallState) EntityUpdateMessage {
	return EntityUpdateMessage{ServerEntityUpdate, player, otherPlayers, walls}
}

func newWallUpdateMessage(walls []WallState) WallUpdateMessage {
	return WallUpdateMessage{ServerWallUpdate, walls}
}

func newMapUpdateMessage(mapState tools.MapState) MapUpdateMessage {
	return MapUpdateMessage{ServerMapUpdate, mapState}
}

func newGameFinishedMessage(winnerName string) GameFinishedMessage {
	return GameFinishedMessage{ServerGameFinished, winnerName}
}
