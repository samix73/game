package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jakecoffman/cp"
	"github.com/samix73/game/game/components"
)

var (
	output   = flag.String("output", "maze.toml", "Output file")
	seed     = flag.Uint64("seed", rand.Uint64(), "Seed")
	width    = flag.Int("width", 10, "Width")
	height   = flag.Int("height", 10, "Height")
	entry    = flag.String("entry", "0,0", "Entry point")
	exit     = flag.String("exit", "9,9", "Exit point")
	cellSize = flag.Float64("cell-size", 1, "Cell size")
)

func parsePoint(point string) (cp.Vector, error) {
	var x, y float64
	if _, err := fmt.Sscanf(point, "%f,%f", &x, &y); err != nil {
		return cp.Vector{}, err
	}
	return cp.Vector{X: x, Y: y}, nil
}

func main() {
	defer func(start time.Time) {
		fmt.Printf("Overall time: %s\n", time.Since(start))
	}(time.Now())

	flag.Parse()

	if *width < 2 || *height < 2 {
		fmt.Println("Width and height must be at least 2")
		os.Exit(1)
	}

	entry, err := parsePoint(*entry)
	if err != nil {
		fmt.Println("Invalid entry point")
		os.Exit(1)
	}
	exit, err := parsePoint(*exit)
	if err != nil {
		fmt.Println("Invalid exit point")
		os.Exit(1)
	}

	m, err := components.GenerateMaze(*width, *height, entry, exit, *seed, *cellSize)
	if err != nil {
		fmt.Println("Invalid seed")
		os.Exit(1)
	}

	entityComponentsConfig := map[string]components.Maze{
		"Maze": m,
	}

	b := new(bytes.Buffer)
	enc := toml.NewEncoder(b)
	if err := enc.Encode(entityComponentsConfig); err != nil {
		panic(err)
	}

	os.WriteFile(*output, b.Bytes(), os.ModePerm)
}
