// internal/world/world.go - Module monde
package world

import (
	"zelda-souls-game/internal/assets"
	"zelda-souls-game/internal/save"
)

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	WindowWidth() int
	WindowHeight() int
	TileSize() int
}

type World struct {
	entities []interface{}
}

type PlayerData struct {
	Name       string
	Difficulty string
	CreatedAt  interface{}
}

func NewWorld(config GameConfig, assetManager *assets.AssetManager) (*World, error) {
	return &World{entities: make([]interface{}, 0)}, nil
}

func (w *World) GetEntities() []interface{}                     { return w.entities }
func (w *World) Cleanup()                                       {}
func (w *World) Reset()                                         {}
func (w *World) InitializeNewGame(playerData *PlayerData) error { return nil }
func (w *World) LoadFromSave(saveData *save.SaveData) error     { return nil }
func (w *World) CreateSaveData() (interface{}, error)           { return nil, nil }
