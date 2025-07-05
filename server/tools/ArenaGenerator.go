package tools

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"
)

type MapState struct {
	Arena         [][]int `json:"arena"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	currentWidth  int
	currentHeight int
	topLeft       Vector2Int
}

func ShrinkArena(mapState *MapState) *MapState {
	if mapState.currentWidth <= 4 || mapState.currentHeight <= 4 {
		fmt.Println("Map too small to shrink further.")
		return mapState
	}

	newWidth := mapState.currentWidth / 2
	newHeight := mapState.currentHeight / 2

	// Calculate the max starting point within current bounds
	maxStartX := mapState.topLeft.X + (mapState.currentWidth - newWidth)
	maxStartY := mapState.topLeft.Y + (mapState.currentHeight - newHeight)

	// Seed randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick a random top-left corner for the new chunk
	startX := r.Intn(maxStartX-mapState.topLeft.X+1) + mapState.topLeft.X
	startY := r.Intn(maxStartY-mapState.topLeft.Y+1) + mapState.topLeft.Y

	// New arena starts as all abyss
	newArena := make([][]int, mapState.Height)
	for y := 0; y < mapState.Height; y++ {
		newArena[y] = make([]int, mapState.Width)
		for x := 0; x < mapState.Width; x++ {
			newArena[y][x] = TILETYPE_ABYSS
		}
	}

	// Copy the shrunken area from the old map
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newArena[startY+y][startX+x] = mapState.Arena[startY+y][startX+x]
		}
	}

	return &MapState{
		Arena:         newArena,
		Width:         mapState.Width,
		Height:        mapState.Height,
		currentWidth:  newWidth,
		currentHeight: newHeight,
		topLeft:       Vector2Int{X: startX, Y: startY},
	}
}

func (m *MapState) DebugLog() {
	fmt.Println("[")
	for i, row := range m.Arena {
		fmt.Print("  [")
		for j, val := range row {
			fmt.Print(val)
			if j < len(row)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print("]")
		if i < len(m.Arena)-1 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("]")
}

const (
	TILETYPE_WALKABLE       = 0
	TILETYPE_ABYSS          = 1
	TILE_SIZE               = 32
	mapArraySize            = 200
	randomFillPercent       = 50
	safeSpawnDistanceSquare = 36
)

func TileToWorldCoords(tileCoords Vector2Int) Vector2 {
	return Vector2{
		X: float64(tileCoords.X*TILE_SIZE + TILE_SIZE/2),
		Y: float64(tileCoords.Y*TILE_SIZE + TILE_SIZE/2),
	}
}

func WorldToTileCoords(worldCoords Vector2) Vector2Int {
	return Vector2Int{
		X: int(worldCoords.X) / TILE_SIZE,
		Y: int(worldCoords.Y) / TILE_SIZE,
	}
}

func GenerateSafeSpawns(number int, mapState *MapState) []Vector2 {
	spawnTiles := make([]Vector2Int, 0, number)
	for i := 0; i < number; i++ {
		spawnTiles = append(spawnTiles, GetSafeTile(spawnTiles, mapState))
	}
	spawns := make([]Vector2, 0, number)
	for i := range spawnTiles {
		spawns = append(spawns, TileToWorldCoords(spawnTiles[i]))
	}
	fmt.Println(spawns)
	return spawns
}

func GetSafeTile(tiles []Vector2Int, mapState *MapState) Vector2Int {
	iteration := 0
	for {
		tile := GetWalkableTile(mapState)
		if iteration > 100 {
			fmt.Println("GetSafeTile() is not returning a safe tile -- 1/10 ragebait")
			return tile
		}
		if len(tiles) == 0 {
			return tile
		} else {
			ok := true
			for i := range tiles {
				if tile.DistanceSquared(tiles[i]) < safeSpawnDistanceSquare {
					ok = false
				}
			}
			if ok {
				return tile
			}
		}
		iteration++
	}
}

func GetWalkableTile(mapState *MapState) Vector2Int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	iteration := 0
	for {
		if iteration > 100 {
			fmt.Println("GetWalkableTile() is not returning a walkable tile -- 0/10 ragebait")
			return Vector2Int{-1, -1}
		}
		row := r.Intn(mapState.Height)
		col := r.Intn(mapState.Width)
		if mapState.Arena[row][col] == TILETYPE_WALKABLE {
			return Vector2Int{X: col, Y: row}
		}
		iteration++
	}
}

func IsPlayerEliminated(mapState *MapState, playerCenter Vector2) bool {
	tileCoords := WorldToTileCoords(playerCenter)
	if tileCoords.X >= mapState.Height || tileCoords.Y >= mapState.Width {
		return true
	}
	tile := mapState.Arena[tileCoords.Y][tileCoords.X]
	if tile == TILETYPE_ABYSS {
		return true
	} else {
		return false
	}
}

func GenerateMap(width, height int, useCustomSeed bool, seedString string) *MapState {
	newArena := RandomFillMap(width, height, useCustomSeed, seedString)
	for i := 0; i < 5; i++ {
		SmoothMap(newArena)
	}
	return newArena
}

func RandomFillMap(width int, height int, useCustomSeed bool, seedString string) *MapState {
	arena := make([][]int, height)
	for r := range arena {
		arena[r] = make([]int, width)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if useCustomSeed {
		h := fnv.New32a()
		h.Write([]byte(seedString))
		seed := int64(h.Sum32())

		r = rand.New(rand.NewSource(seed))
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			if r.Intn(100) < randomFillPercent {
				arena[x][y] = TILETYPE_ABYSS
			} else {
				arena[x][y] = TILETYPE_WALKABLE
			}
		}
	}

	return &MapState{Arena: arena, Width: width, Height: height, currentWidth: width, currentHeight: height}

}

func SmoothMap(arena *MapState) {
	for x := 0; x < arena.Height; x++ {
		for y := 0; y < arena.Width; y++ {
			neighborAbyssTiles := GetSurroundingAbyssCount(*arena, x, y)

			if neighborAbyssTiles > 4 {
				arena.Arena[x][y] = TILETYPE_ABYSS
			} else if neighborAbyssTiles < 4 {
				arena.Arena[x][y] = TILETYPE_WALKABLE
			}
		}
	}
}

func GetSurroundingAbyssCount(arena MapState, gridX int, gridY int) int {
	wallCount := 0
	for neighborRow := gridX - 1; neighborRow <= gridX+1; neighborRow++ {
		for neighborCol := gridY - 1; neighborCol <= gridY+1; neighborCol++ {
			if neighborRow >= 0 && neighborRow < arena.Height && neighborCol >= 0 && neighborCol < arena.Width {
				if neighborRow != gridX || neighborCol != gridY {
					wallCount += arena.Arena[neighborRow][neighborCol]
				}
			} else {
				wallCount++
			}
		}
	}
	return wallCount
}
