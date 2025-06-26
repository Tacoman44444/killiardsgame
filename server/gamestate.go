package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

type PlayerState struct {
	Position_x float64
	Position_y float64
	Velocity_x float64
	Velocity_y float64
}

type PlayerAction struct {
	power              int
	DirectionHorizonal float64
	DirectionVertical  float64
}

type WallState struct {
	position_x float64
	position_y float64
	turnsLeft  int
}

type GameState struct {
	players    map[string]PlayerState
	mapState   tools.MapState
	walls      []WallState
	activeTurn string
}

func GetNewGame(players []string) *GameState {

	mapState := tools.GenerateMap(200, 200, false, "")
	playerMap := make(map[string]PlayerState, len(players))
	safeSpawns := tools.GenerateSafeSpawns(len(players), &mapState)
	for i := range players {
		playerMap[players[i]] = PlayerState{safeSpawns[i].X, safeSpawns[i].Y, 0, 0}
	}
	walls := make([]WallState, 10)
	turn := players[0]

	gamestate := GameState{
		players:    playerMap,
		mapState:   mapState,
		walls:      walls,
		activeTurn: turn,
	}

	return &gamestate

}
