package main

type PlayerState struct {
	position_x int
	position_y int
	velocity_x int
	velocity_y int
}

type MapState struct {
	arena [][]int
}

type WallState struct {
	position_x int
	position_y int
	turnsLeft  int
}

type GameState struct {
	players    map[int]PlayerState
	mapState   MapState
	walls      []WallState
	activeTurn int
}
