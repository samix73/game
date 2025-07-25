package assets

import (
	"embed"
	"fmt"
)

//go:embed sprites/*
var sprites embed.FS

func GetSprite(name string) ([]byte, error) {
	extensions := []string{
		".png",
		".jpg",
		".jpeg",
	}

	data, err := sprites.ReadFile("sprites/" + name)
	if err == nil {
		return data, nil
	}

	for _, ext := range extensions {
		data, err := sprites.ReadFile("sprites/" + name + ext)
		if err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("sprite %s not found", name)
}
