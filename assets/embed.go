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

func GetSprite(name string) (*ebiten.Image, error) {
	data, err := sprites.ReadFile(SpritesDir + name)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("assets.GetSprite: %w", err)
	}

	return ebiten.NewImageFromImage(img), nil
}
