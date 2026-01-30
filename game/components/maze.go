package components

import (
	"errors"
	"math/rand/v2"

	"github.com/jakecoffman/cp"
	ecs "github.com/samix73/ebiten-ecs"
)

func init() {
	ecs.RegisterComponent[Maze]()
}

type Maze struct {
	Width    int
	Height   int
	Cells    [][]uint8
	Seed     uint64
	CellSize float64
}

func (m *Maze) Init() {
	m.Width = 0
	m.Height = 0
	m.Cells = nil
	m.CellSize = 0
}

func (m *Maze) Reset() {
	m.Width = 0
	m.Height = 0
	m.Cells = nil
	m.CellSize = 0
}

func (m *Maze) Clone() ecs.Component {
	clone := *m
	return &clone
}

func GenerateMaze(width, height int, entry, exit cp.Vector, seed uint64, cellSize float64) (Maze, error) {
	if entry == exit {
		return Maze{}, errors.New("entry and exit must be different")
	}

	if entry.X < 0 || entry.Y < 0 || entry.X >= float64(width) || entry.Y >= float64(height) {
		return Maze{}, errors.New("entry point is out of bounds")
	}

	if exit.X < 0 || exit.Y < 0 || exit.X >= float64(width) || exit.Y >= float64(height) {
		return Maze{}, errors.New("exit point is out of bounds")
	}

	// Initialize all cells as walls (false)
	cells := make([][]uint8, height)
	for i := range cells {
		cells[i] = make([]uint8, width)
	}

	r := rand.NewPCG(seed, 0)

	// Create a guaranteed path from entry to exit using random walk
	entryX, entryY := int(entry.X), int(entry.Y)
	exitX, exitY := int(exit.X), int(exit.Y)

	// Mark entry and exit as open
	cells[entryY][entryX] = 1
	cells[exitY][exitX] = 1

	// Random walk from entry to exit to ensure connectivity
	currentX, currentY := entryX, entryY

	for currentX != exitX || currentY != exitY {
		// Mark current cell as open
		cells[currentY][currentX] = 1

		// Calculate direction bias towards exit
		directions := make([][2]int, 0, 4)

		// Prefer moving towards exit
		if currentX < exitX && currentX+1 < width {
			directions = append(directions, [2]int{1, 0}) // right
		}
		if currentX > exitX && currentX-1 >= 0 {
			directions = append(directions, [2]int{-1, 0}) // left
		}
		if currentY < exitY && currentY+1 < height {
			directions = append(directions, [2]int{0, 1}) // down
		}
		if currentY > exitY && currentY-1 >= 0 {
			directions = append(directions, [2]int{0, -1}) // up
		}

		// Add other valid directions
		allDirections := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, dir := range allDirections {
			newX, newY := currentX+dir[0], currentY+dir[1]
			if newX >= 0 && newX < width && newY >= 0 && newY < height {
				// Add non-biased directions with lower priority
				alreadyAdded := false
				for _, d := range directions {
					if d[0] == dir[0] && d[1] == dir[1] {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					directions = append(directions, dir)
				}
			}
		}

		if len(directions) == 0 {
			// Should not happen with valid entry/exit, but handle gracefully
			break
		}

		// Choose a direction (biased towards exit due to ordering)
		chosenDir := directions[r.Uint64()%uint64(len(directions))]
		currentX += chosenDir[0]
		currentY += chosenDir[1]
	}

	// Fill remaining cells randomly, keeping the guaranteed path
	for y := range cells {
		for x := range cells[y] {
			// Skip cells that are already marked as open (part of the path)
			if cells[y][x] == 0 {
				cells[y][x] = uint8(r.Uint64() % 2)
			}
		}
	}

	return Maze{
		Width:    width,
		Height:   height,
		Cells:    cells,
		Seed:     seed,
		CellSize: cellSize,
	}, nil
}
