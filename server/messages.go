package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

type ServerMessageType string

const (
	ServerRoomCreated   ServerMessageType = "room_created"
	ServerRoomJoined    ServerMessageType = "room_joined"
	ServerInvalidCode   ServerMessageType = "invalid-code"
	ServerGameStart     ServerMessageType = "game_started"
	ServerTurnStart     ServerMessageType = "turn_started"
	ServerTurnTimeout   ServerMessageType = "turn_timeout"
	ServerBroadcastTurn ServerMessageType = "broadcast-turn"
	ServerEntityUpdate  ServerMessageType = "entity_update"
	ServerWallUpdate    ServerMessageType = "wall_update"
	ServerMapUpdate     ServerMessageType = "map_update"
	ServerGameFinished  ServerMessageType = "game_finished"
)

type ServerMessage interface {
	isServerMessage()
}

type RoomCreatedMessage struct {
	Type ServerMessageType `json:"type"`
	Code string            `json:"code"`
}

func (m RoomCreatedMessage) isServerMessage() {}

type RoomJoinedMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m RoomJoinedMessage) isServerMessage() {}

type InvalidCodeMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m InvalidCodeMessage) isServerMessage() {}

type GameStartMessage struct {
	Type         ServerMessageType `json:"type"`
	Map          tools.MapState    `json:"map"`
	Player       PlayerIdentity    `json:"player"`
	OtherPlayers []PlayerIdentity  `json:"other_players"`
}

func (m GameStartMessage) isServerMessage() {}

type TurnStartMessage struct {
	Type ServerMessageType `json:"type"`
	Id   string            `json:"id"`
}

func (m TurnStartMessage) isServerMessage() {}

type TurnTimeoutMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m TurnTimeoutMessage) isServerMessage() {}

type BroadcastTurnMessage struct {
	Type   ServerMessageType `json:"type"`
	Player PlayerIdentity    `json:"player"`
	Action PlayerAction      `json:"action"`
}

func (m BroadcastTurnMessage) isServerMessage() {}

type EntityUpdateMessage struct {
	Type         ServerMessageType `json:"type"`
	Player       PlayerIdentity    `json:"player_state"`
	OtherPlayers []PlayerIdentity  `json:"other_players"`
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

func (m GameFinishedMessage) isServerMessage() {}

// CREATING NEW MESSAGES

func newRoomCreatedMessage(code string) RoomCreatedMessage {
	return RoomCreatedMessage{ServerRoomCreated, code}
}

func newRoomJoinedMessage() RoomJoinedMessage {
	return RoomJoinedMessage{ServerRoomJoined}
}

func newInvalidCodeMessage() InvalidCodeMessage {
	return InvalidCodeMessage{ServerInvalidCode}
}

func newGameStartMessage(mapState tools.MapState, player PlayerIdentity, otherPlayers []PlayerIdentity) GameStartMessage {
	return GameStartMessage{ServerGameStart, mapState, player, otherPlayers}
}

func newTurnStartMessage(playerID string) TurnStartMessage {
	return TurnStartMessage{ServerTurnStart, playerID}
}

func newTurnTimeoutMessage() TurnTimeoutMessage {
	return TurnTimeoutMessage{ServerTurnTimeout}
}

func newBroadcastTurnMessage(player PlayerIdentity, action PlayerAction) BroadcastTurnMessage {
	return BroadcastTurnMessage{ServerBroadcastTurn, player, action}
}

func newEntityUpdateMessage(player PlayerIdentity, otherPlayers []PlayerIdentity, walls []WallState) EntityUpdateMessage {
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

type ClientMessageType string

const (
	ClientCreateRoom     ClientMessageType = "create-room"
	ClientStartGame      ClientMessageType = "start-game"
	ClientJoinRoom       ClientMessageType = "join-room"
	ClientSendWall       ClientMessageType = "send-wall"
	ClientSendTurn       ClientMessageType = "send-turn"
	ClientSendId         ClientMessageType = "send-id"
	ClientSimulationDone ClientMessageType = "simulation-done"
)

type ClientMessage struct {
	Type   ClientMessageType `json:"type"`
	Id     string            `json:"id"`
	Code   string            `json:"code"`
	Wall   WallState         `json:"wall_state"`
	Action PlayerAction      `json:"player_action"`
}
