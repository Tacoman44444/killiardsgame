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
	players    map[int]PlayerState
	mapState   tools.MapState
	walls      []WallState
	activeTurn int
}
