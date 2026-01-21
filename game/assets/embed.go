package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SpritesDir = "Sprites"
	WorldsDir  = "Worlds"
)

//go:embed Sprites/*
var sprites embed.FS

var (
	spriteCache = make(map[string]*ebiten.Image, 10)
)

func GetSprite(name string) (*ebiten.Image, error) {
	if v, ok := spriteCache[name]; ok {
		return v, nil
	}

	path := filepath.Join(SpritesDir, name)

	data, err := sprites.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("assets.GetSprite: %w", err)
	}

	eImg := ebiten.NewImageFromImage(img)
	spriteCache[name] = eImg

	return eImg, nil
}

//go:embed Worlds/*.hcl
var worlds embed.FS

func GetWorld(name string) ([]byte, error) {
	return worlds.ReadFile(filepath.Join(WorldsDir, name))
}
