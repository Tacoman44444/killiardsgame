package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

type ServerMessageType string

const (
	ServerRoomCreated   ServerMessageType = "room-created"
	ServerRoomJoined    ServerMessageType = "room-joined"
	ServerInvalidCode   ServerMessageType = "invalid-code"
	ServerMakeOwner     ServerMessageType = "make-owner"
	ServerGameStart     ServerMessageType = "game-start"
	ServerTurnStart     ServerMessageType = "turn-started"
	ServerTurnTimeout   ServerMessageType = "turn-timeout"
	ServerBroadcastTurn ServerMessageType = "broadcast-turn"
	ServerEntityUpdate  ServerMessageType = "entity-update"
	ServerEliminations  ServerMessageType = "e"
	ServerWallUpdate    ServerMessageType = "wall-update"
	ServerMapUpdate     ServerMessageType = "map-update"
	ServerGameFinished  ServerMessageType = "game-finished"
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

type MakeOwnerMessage struct {
	Type ServerMessageType `json:"type"`
}

func (m MakeOwnerMessage) isServerMessage() {}

type GameStartMessage struct {
	Type         ServerMessageType      `json:"type"`
	CurrentMap   tools.MapState         `json:"current_map"`
	NextMap      tools.MapState         `json:"next_map"`
	Player       ClientPlayerIdentity   `json:"player"`
	OtherPlayers []ClientPlayerIdentity `json:"other_players"`
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
	Type   ServerMessageType    `json:"type"`
	Player ClientPlayerIdentity `json:"player"`
	Action PlayerAction         `json:"action"`
}

func (m BroadcastTurnMessage) isServerMessage() {}

type EntityUpdateMessage struct {
	Type         ServerMessageType      `json:"type"`
	Player       ClientPlayerIdentity   `json:"player"`
	OtherPlayers []ClientPlayerIdentity `json:"all_players"`
	Walls        []WallState            `json:"walls"`
}

func (m EntityUpdateMessage) isServerMessage() {}

type EliminationMessage struct {
	Type         ServerMessageType      `json:"type"`
	Eliminations []ClientPlayerIdentity `json:"eliminated_players"`
}

func (m EliminationMessage) isServerMessage() {}

type WallUpdateMessage struct {
	Type  ServerMessageType `json:"type"`
	Walls []WallState       `json:"walls"`
}

func (m WallUpdateMessage) isServerMessage() {}

type MapUpdateMessage struct {
	Type       ServerMessageType `json:"type"`
	CurrentMap tools.MapState    `json:"current_map"`
	NextMap    tools.MapState    `json:"next_map"`
}

func (m MapUpdateMessage) isServerMessage() {}

type GameFinishedMessage struct {
	Type       ServerMessageType `json:"type"`
	Result     string            `json:"result"` //win or draw
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

func newMakeOwnerMessage() MakeOwnerMessage {
	return MakeOwnerMessage{ServerMakeOwner}
}

func newGameStartMessage(currentMap tools.MapState, nextMap tools.MapState, player PlayerIdentity, otherPlayers []PlayerIdentity) GameStartMessage {
	return GameStartMessage{ServerGameStart, currentMap, nextMap, FormatPlayerId(player), FormatPlayerIds(otherPlayers)}
}

func newTurnStartMessage(playerID string) TurnStartMessage {
	return TurnStartMessage{ServerTurnStart, playerID}
}

func newTurnTimeoutMessage() TurnTimeoutMessage {
	return TurnTimeoutMessage{ServerTurnTimeout}
}

func newBroadcastTurnMessage(player PlayerIdentity, action PlayerAction) BroadcastTurnMessage {
	return BroadcastTurnMessage{ServerBroadcastTurn, FormatPlayerId(player), action}
}

func newEntityUpdateMessage(player PlayerIdentity, allPlayers []PlayerIdentity, walls []WallState) EntityUpdateMessage {
	return EntityUpdateMessage{ServerEntityUpdate, FormatPlayerId(player), FormatPlayerIds(allPlayers), walls}
}

func newEliminationMessage(elims []PlayerIdentity) EliminationMessage {
	return EliminationMessage{ServerEliminations, FormatPlayerIds(elims)}
}

func newWallUpdateMessage(walls []WallState) WallUpdateMessage {
	return WallUpdateMessage{ServerWallUpdate, walls}
}

func newMapUpdateMessage(currentMap tools.MapState, nextMap tools.MapState) MapUpdateMessage {
	return MapUpdateMessage{ServerMapUpdate, currentMap, nextMap}
}

func newGameFinishedMessage(result string, winnerName string) GameFinishedMessage {
	return GameFinishedMessage{ServerGameFinished, result, winnerName}
}

type ClientMessageType string

const (
	ClientCreateRoom     ClientMessageType = "create-room"
	ClientStartGame      ClientMessageType = "start-game"
	ClientJoinRoom       ClientMessageType = "join-room"
	ClientLeaveRoom      ClientMessageType = "leave-room"
	ClientSendWall       ClientMessageType = "send-wall"
	ClientSendTurn       ClientMessageType = "send-turn"
	ClientSendId         ClientMessageType = "send-id"
	ClientSimulationDone ClientMessageType = "simulation-done"
)

type ClientMessage struct {
	Type     ClientMessageType `json:"type"`
	Id       string            `json:"id"`
	Username string            `json:"username"`
	JoinData PlayerJoinData    `json:"data"`
	Wall     WallState         `json:"wall_state"`
	Action   PlayerAction      `json:"player_action"`
}
