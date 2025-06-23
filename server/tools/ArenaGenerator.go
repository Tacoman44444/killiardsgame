package tools

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"
)

type MapState struct {
	Arena  [][]int
	Width  int
	Height int
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
	TILETYPE_WALKABLE = 0
	TILETYPE_ABYSS    = 1
	mapArraySize      = 200
	randomFillPercent = 50
)

func GenerateMap(width, height int, playerPositions []Circle, useCustomSeed bool, seedString string) MapState {
	newArena := RandomFillMap(width, height, useCustomSeed, seedString)
	for i := 0; i < 5; i++ {
		SmoothMap(&newArena)
	}
	// NOT YET IMPLEMENTED A FUNCTION TO MAKE PLAYER TILES WALKABLES
	return newArena
}

func RandomFillMap(width int, height int, useCustomSeed bool, seedString string) MapState {
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

	return MapState{Arena: arena, Width: width, Height: height}

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
