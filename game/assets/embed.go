package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SpritesDir  = "Sprites"
	WorldsDir   = "Worlds"
	EntitiesDir = "Entities"
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

	data, err := sprites.ReadFile(SpritesDir + "/" + name)
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

//go:embed Worlds/*.toml
var worlds embed.FS

func GetWorld(name string) ([]byte, error) {
	f, err := worlds.ReadFile(WorldsDir + "/" + name + ".toml")
	if err != nil {
		return nil, fmt.Errorf("assets.GetWorld: %w", err)
	}

	return f, nil
}

//go:embed Entities/*.toml
var entities embed.FS

func GetEntity(name string) ([]byte, error) {
	f, err := entities.ReadFile(EntitiesDir + "/" + name + ".toml")
	if err != nil {
		return nil, fmt.Errorf("assets.GetEntity: %w", err)
	}

	return f, nil
}
