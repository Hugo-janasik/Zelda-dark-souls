// internal/states/adapter.go - Adaptateur pour l'interface core
package states

import (
	"time"
)

// Types copiés de core pour éviter les cycles d'import
type GameStateType string

const (
	StateDemo     GameStateType = "demo"
	StateMenu     GameStateType = "menu"
	StatePause    GameStateType = "pause"
	StateGameplay GameStateType = "gameplay"
)

// Renderer interface pour éviter les cycles (copié de core)
type CoreRenderer interface {
	DrawText(text string, pos CoreVector2, color CoreColor)
	DrawRectangle(rect CoreRectangle, color CoreColor, filled bool)
}

// Types géométriques pour l'interface core
type CoreVector2 struct {
	X, Y float64
}

type CoreColor struct {
	R, G, B, A uint8
}

type CoreRectangle struct {
	X, Y, Width, Height float64
}

// StateManagerAdapter adapte SimpleStateManager pour l'interface core
type StateManagerAdapter struct {
	stateManager *SimpleStateManager
}

// NewStateManagerAdapter crée un adaptateur
func NewStateManagerAdapter(sm *SimpleStateManager) *StateManagerAdapter {
	return &StateManagerAdapter{
		stateManager: sm,
	}
}

// Update met à jour l'état (interface core)
func (a *StateManagerAdapter) Update(deltaTime time.Duration) error {
	return a.stateManager.Update(deltaTime)
}

// Render rend l'état avec adaptation des types
func (a *StateManagerAdapter) Render(renderer CoreRenderer) error {
	// Créer un adaptateur de renderer
	rendererAdapter := &RendererAdapter{coreRenderer: renderer}
	return a.stateManager.Render(rendererAdapter)
}

// GetCurrentStateType retourne le type d'état actuel (interface core)
func (a *StateManagerAdapter) GetCurrentStateType() GameStateType {
	currentState := a.stateManager.GetCurrentStateType()
	return GameStateType(currentState)
}

// ChangeState change l'état (interface core)
func (a *StateManagerAdapter) ChangeState(stateType GameStateType) {
	a.stateManager.ChangeState(string(stateType))
}

// RendererAdapter adapte le renderer core vers notre interface locale
type RendererAdapter struct {
	coreRenderer CoreRenderer
}

// DrawText adapte l'appel de rendu de texte
func (r *RendererAdapter) DrawText(text string, pos Vector2, color Color) {
	corePos := CoreVector2{X: pos.X, Y: pos.Y}
	coreColor := CoreColor{R: color.R, G: color.G, B: color.B, A: color.A}
	r.coreRenderer.DrawText(text, corePos, coreColor)
}

// DrawRectangle adapte l'appel de rendu de rectangle
func (r *RendererAdapter) DrawRectangle(rect Rectangle, color Color, filled bool) {
	coreRect := CoreRectangle{X: rect.X, Y: rect.Y, Width: rect.Width, Height: rect.Height}
	coreColor := CoreColor{R: color.R, G: color.G, B: color.B, A: color.A}
	r.coreRenderer.DrawRectangle(coreRect, coreColor, filled)
}

// Méthodes supplémentaires pour SimpleStateManager
func (a *StateManagerAdapter) ToggleInstructions() {
	a.stateManager.ToggleInstructions()
}

func (a *StateManagerAdapter) GetStateManager() *SimpleStateManager {
	return a.stateManager
}
