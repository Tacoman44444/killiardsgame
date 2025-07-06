package main

import (
	"time"

	"github.com/Tacoman44444/killiardsgame/server/tools"
)

const (
	TURN_TIMER_IN_SECONDS                = 3
	CLIENT_AFFIRMATON_TIMEOUT_IN_SECONDS = 10
	PUCK_RADIUS                          = 16.0
)

type PlayerIdentity struct {
	id       string
	circle   *tools.Circle
	username string
}

type ClientPlayerIdentity struct {
	Id        string  `json:"id"`
	PositionX float64 `json:"position_x"`
	PositionY float64 `json:"position_y"`
	VelocityX float64 `json:"velocity_x"`
	VelocityY float64 `json:"velocity_y"`
	Username  string  `json:"username"`
}

func FormatPlayerId(playerId PlayerIdentity) ClientPlayerIdentity {
	return ClientPlayerIdentity{
		Id:        playerId.id,
		PositionX: playerId.circle.Center.X,
		PositionY: playerId.circle.Center.Y,
		VelocityX: playerId.circle.Velocity.X,
		VelocityY: playerId.circle.Velocity.Y,
	}
}

func FormatPlayerIds(playerIds []PlayerIdentity) []ClientPlayerIdentity {
	formattedIds := make([]ClientPlayerIdentity, 0)
	for i := range playerIds {
		formattedIds = append(formattedIds, FormatPlayerId(playerIds[i]))
	}
	return formattedIds
}

type PlayerAction struct {
	Power               int     `json:"power"`
	DirectionHorizontal float64 `json:"direction_horizontal"`
	DirectionVertical   float64 `json:"direction_vertical"`
}

func PlayerActionToShotData(action PlayerAction) tools.ShotData {
	return tools.ShotData{
		Power:     action.Power,
		Direction: tools.Vector2{X: action.DirectionHorizontal, Y: action.DirectionVertical},
	}
}

type WallState struct {
	rect      *tools.Rect
	turnsLeft int
}

func WallStateRefToWallState(wallState []*WallState) []WallState {
	if wallState == nil {
		return []WallState{}
	}
	w := make([]WallState, 0, len(wallState))
	for _, ws := range wallState {
		if ws != nil {
			w = append(w, *ws)
		}
	}
	return w
}

type GameState struct {
	players   map[string]*PlayerIdentity
	mapState  *tools.MapState
	nextMap   *tools.MapState
	walls     []*WallState
	turnTimer time.Duration
}

func GetNewGame(playerIDs []string, playerUsernames []string) *GameState {

	mapState := tools.GenerateMap(200, 200, false, "")
	nextMap := tools.ShrinkArena(mapState)
	playerMap := make(map[string]*PlayerIdentity, len(playerIDs))
	safeSpawns := tools.GenerateSafeSpawns(len(playerIDs), mapState)
	for i := range playerIDs {
		playerMap[playerIDs[i]] = &PlayerIdentity{playerIDs[i], &tools.Circle{Center: safeSpawns[i], Radius: PUCK_RADIUS, Velocity: tools.Vector2{X: 0, Y: 0}}, playerUsernames[i]}
	}
	walls := make([]*WallState, 0, 10)

	gamestate := GameState{
		players:   playerMap,
		mapState:  mapState,
		nextMap:   nextMap,
		walls:     walls,
		turnTimer: time.Duration(TURN_TIMER_IN_SECONDS) * time.Second,
	}

	return &gamestate
}

func PlayerMapToSlice(playerMap map[string]*PlayerIdentity) []PlayerIdentity {
	players := make([]PlayerIdentity, 0, len(playerMap))
	for _, player := range playerMap {
		players = append(players, *player)
	}
	return players
}

func PlayerMapToSliceRef(playerMap map[string]*PlayerIdentity) []*PlayerIdentity {
	players := make([]*PlayerIdentity, 0, len(playerMap))
	for _, player := range playerMap {
		players = append(players, player)
	}
	return players
}

func PlayerMapToCircles(playerMap map[string]*PlayerIdentity) []*tools.Circle {
	playerSlice := PlayerMapToSliceRef(playerMap)
	circleSlice := make([]*tools.Circle, 0, len(playerSlice))
	for i := range playerSlice {
		circleSlice = append(circleSlice, playerSlice[i].circle)
	}
	return circleSlice
}

func GetWallRectRefs(walls []*WallState) []*tools.Rect {
	rects := make([]*tools.Rect, 0, len(walls))
	if len(walls) == 0 {
		return rects
	}
	for i := range walls {
		rects = append(rects, walls[i].rect)
	}
	return rects
}
