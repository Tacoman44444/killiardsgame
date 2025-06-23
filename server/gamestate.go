package main

import "github.com/Tacoman44444/killiardsgame/server/tools"

type PlayerState struct {
	position_x int
	position_y int
	velocity_x int
	velocity_y int
}

type WallState struct {
	position_x int
	position_y int
	turnsLeft  int
}

type GameState struct {
	players    map[int]PlayerState
	mapState   tools.MapState
	walls      []WallState
	activeTurn int
}
