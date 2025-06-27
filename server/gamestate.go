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
	id     string
	circle *tools.Circle
}

type PlayerAction struct {
	data tools.ShotData
}

type WallState struct {
	rect      *tools.Rect
	turnsLeft int
}

func WallStateRefToWallState(wallState []*WallState) []WallState {
	w := make([]WallState, len(wallState))
	for i := range wallState {
		w = append(w, *wallState[i])
	}
	return w
}

type GameState struct {
	players   map[string]*PlayerIdentity
	mapState  *tools.MapState
	walls     []*WallState
	turnTimer time.Duration
}

func GetNewGame(players []string) *GameState {

	mapState := tools.GenerateMap(200, 200, false, "")
	playerMap := make(map[string]*PlayerIdentity, len(players))
	safeSpawns := tools.GenerateSafeSpawns(len(players), mapState)
	for i := range players {
		playerMap[players[i]] = &PlayerIdentity{players[i], &tools.Circle{Center: safeSpawns[i], Radius: PUCK_RADIUS, Velocity: tools.Vector2{X: 0, Y: 0}}}
	}
	walls := make([]*WallState, 10)

	gamestate := GameState{
		players:   playerMap,
		mapState:  mapState,
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
	circleSlice := make([]*tools.Circle, len(playerSlice))
	for i := range playerSlice {
		circleSlice = append(circleSlice, playerSlice[i].circle)
	}
	return circleSlice
}

func GetWallRectRefs(walls []*WallState) []*tools.Rect {
	rects := make([]*tools.Rect, len(walls))
	for i := range walls {
		rects = append(rects, walls[i].rect)
	}
	return rects
}
