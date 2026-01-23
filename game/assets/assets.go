package assets

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SpritesDir  = "game/assets/Sprites"
	WorldsDir   = "game/assets/Worlds"
	EntitiesDir = "game/assets/Entities"
)

func GetSprite(name string) (*ebiten.Image, error) {
	path := path.Join(SpritesDir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("assets.GetSprite: %w", err)
	}

	return ebiten.NewImageFromImage(img), nil
}

func GetWorld(name string) ([]byte, error) {
	path := path.Join(WorldsDir, name+".toml")
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("assets.GetWorld: %w", err)
	}

	return f, nil
}

func GetEntity(name string) ([]byte, error) {
	path := path.Join(EntitiesDir, name+".toml")
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("assets.GetEntity: %w", err)
	}

	return f, nil
}
