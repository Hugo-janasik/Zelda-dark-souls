// internal/ui/ui_manager.go - Gestionnaire UI
package ui

import (
	"time"
	"zelda-souls-game/internal/rendering"
)

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	WindowWidth() int
	WindowHeight() int
}

type UIManager struct {
	config   GameConfig
	renderer *rendering.Renderer
}

func NewUIManager(config GameConfig, renderer *rendering.Renderer) *UIManager {
	return &UIManager{
		config:   config,
		renderer: renderer,
	}
}

func (ui *UIManager) Update(deltaTime time.Duration) {
	// TODO: Mettre à jour les éléments UI
}

func (ui *UIManager) Render(renderer *rendering.Renderer) {
	// TODO: Rendre l'interface utilisateur
}
