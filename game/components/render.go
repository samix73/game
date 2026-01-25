package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samix73/game/ecs"
	"github.com/samix73/game/game/assets"
)

func init() {
	if err := ecs.RegisterComponent[Renderable](); err != nil {
		panic(err)
	}
	if err := ecs.RegisterComponent[Render](); err != nil {
		panic(err)
	}
}

// Renderable represents a 2D entity that can be rendered on the screen.
type Renderable struct {
	Order      int // Rendering order; lower values are rendered first
	SpritePath string
	Sprite     *ebiten.Image `toml:"-"`
	GeoM       ebiten.GeoM   `toml:"-"`
}

func (r *Renderable) UnmarshalTOML(data any) error {
	d, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("failed to decode renderable: expected map[string]any, got %T", data)
	}

	r.SpritePath = d["SpritePath"].(string)
	r.Order = int(d["Order"].(int64))

	var err error
	r.Sprite, err = assets.GetSprite(r.SpritePath)
	if err != nil {
		return fmt.Errorf("failed to decode renderable: %w", err)
	}

	return nil
}

func (r *Renderable) Reset() {
	if r.Sprite != nil {
		r.Sprite.Deallocate()
	}
	r.GeoM.Reset()
	r.Order = 0
}

// Render marks entities to be rendered.
type Render struct{}
