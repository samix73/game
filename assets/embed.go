package assets

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"runtime/trace"

	"github.com/hajimehoshi/ebiten/v2"
)

const SpritesDir = "Sprites/"

//go:embed Sprites/*
var sprites embed.FS

func GetSprite(ctx context.Context, name string) (*ebiten.Image, error) {
	region := trace.StartRegion(ctx, "assets.GetSprite")
	defer region.End()

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
