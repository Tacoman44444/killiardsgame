package tools

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"
)

type MapState struct {
	Arena  [][]int `json:"arena"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
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
	TILE_SIZE               = 64
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

	return &MapState{Arena: arena, Width: width, Height: height}

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
