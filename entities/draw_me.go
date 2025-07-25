package entities

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/samix73/game/components"
	"github.com/samix73/game/ecs"
)

func NewDrawMeEntity(em *ecs.EntityManager) (ecs.EntityID, error) {
	img, _, err := image.Decode(bytes.NewReader(images.Ebiten_png))
	if err != nil {
		return ecs.UndefinedID, fmt.Errorf("error decoding image: %v", err)
	}

	entity := em.NewEntity()
	ecs.AddComponent[components.Transform](em, entity)
	renderable := ecs.AddComponent[components.Renderable](em, entity)

	renderable.Sprite = ebiten.NewImageFromImage(img)

	return entity, nil
}
