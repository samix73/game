package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const SpritesDir = "Sprites/"

//go:embed Sprites/*
var sprites embed.FS

var (
	spriteCache = make(map[string]*ebiten.Image, 10)
)

func GetSprite(name string) (*ebiten.Image, error) {
	if v, ok := spriteCache[name]; ok {
		return v, nil
	}

	data, err := sprites.ReadFile(SpritesDir + name)
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
