// internal/states/simple_state_manager.go - StateManager minimal
package states

import (
	"fmt"
	"time"
)

// Renderer interface minimale pour éviter les cycles
type Renderer interface {
	DrawText(text string, pos Vector2, color Color)
	DrawRectangle(rect Rectangle, color Color, filled bool)
}

// Structures minimales pour éviter les imports
type Vector2 struct {
	X, Y float64
}

type Color struct {
	R, G, B, A uint8
}

type Rectangle struct {
	X, Y, Width, Height float64
}

// Couleurs prédéfinies
var (
	ColorWhite  = Color{255, 255, 255, 255}
	ColorGreen  = Color{0, 255, 0, 255}
	ColorYellow = Color{255, 255, 0, 255}
	ColorRed    = Color{255, 0, 0, 255}
)

// SimpleStateManager gestionnaire d'états minimal
type SimpleStateManager struct {
	currentState     string
	frameCount       int
	showInstructions bool
}

// NewSimpleStateManager crée un gestionnaire d'états minimal
func NewSimpleStateManager() *SimpleStateManager {
	return &SimpleStateManager{
		currentState:     "demo",
		frameCount:       0,
		showInstructions: true,
	}
}

// Update met à jour l'état
func (sm *SimpleStateManager) Update(deltaTime time.Duration) error {
	sm.frameCount++
	return nil
}

// Render rend l'état actuel
func (sm *SimpleStateManager) Render(renderer Renderer) error {
	switch sm.currentState {
	case "demo":
		sm.renderDemoState(renderer)
	case "menu":
		sm.renderMenuState(renderer)
	default:
		sm.renderDemoState(renderer)
	}
	return nil
}

// renderDemoState rend l'état de démonstration
func (sm *SimpleStateManager) renderDemoState(renderer Renderer) {
	// Titre principal
	renderer.DrawText("Zelda Souls Game", Vector2{100, 100}, ColorWhite)
	renderer.DrawText("Systèmes de base opérationnels !", Vector2{100, 130}, ColorGreen)

	// Instructions
	if sm.showInstructions {
		renderer.DrawText("Contrôles:", Vector2{100, 180}, ColorYellow)
		renderer.DrawText("ESC - Changer d'état", Vector2{120, 200}, ColorWhite)
		renderer.DrawText("ZQSD ou WASD - Test mouvement", Vector2{120, 220}, ColorWhite)
		renderer.DrawText("I - Toggle instructions", Vector2{120, 240}, ColorWhite)
	}

	// Compteur de frames pour montrer que ça tourne
	frameText := fmt.Sprintf("Frames: %d", sm.frameCount)
	renderer.DrawText(frameText, Vector2{100, 300}, ColorWhite)

	// État du jeu
	stateText := fmt.Sprintf("État: %s", sm.currentState)
	renderer.DrawText(stateText, Vector2{100, 320}, ColorWhite)
}

// renderMenuState rend l'état menu
func (sm *SimpleStateManager) renderMenuState(renderer Renderer) {
	renderer.DrawText("=== MENU PRINCIPAL ===", Vector2{100, 100}, ColorYellow)
	renderer.DrawText("1. Nouvelle partie", Vector2{100, 150}, ColorWhite)
	renderer.DrawText("2. Charger partie", Vector2{100, 170}, ColorWhite)
	renderer.DrawText("3. Options", Vector2{100, 190}, ColorWhite)
	renderer.DrawText("4. Quitter", Vector2{100, 210}, ColorWhite)

	renderer.DrawText("ESC - Retour démo", Vector2{100, 250}, ColorGreen)
}

// GetCurrentStateType retourne le type d'état actuel
func (sm *SimpleStateManager) GetCurrentStateType() string {
	return sm.currentState
}

// ChangeState change l'état
func (sm *SimpleStateManager) ChangeState(stateType string) {
	fmt.Printf("Changement d'état: %s -> %s\n", sm.currentState, stateType)
	sm.currentState = stateType
}

// ToggleInstructions active/désactive les instructions
func (sm *SimpleStateManager) ToggleInstructions() {
	sm.showInstructions = !sm.showInstructions
	fmt.Printf("Instructions: %t\n", sm.showInstructions)
}
