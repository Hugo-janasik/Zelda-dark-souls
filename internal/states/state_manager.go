// internal/states/state_manager.go - Gestionnaire d'états (stub)
package states

import (
	"fmt"
	"time"

	"zelda-souls-game/internal/core"
	"zelda-souls-game/internal/ecs/systems"
	"zelda-souls-game/internal/input"
	"zelda-souls-game/internal/rendering"
	"zelda-souls-game/internal/save"
	"zelda-souls-game/internal/ui"
	"zelda-souls-game/internal/world"
)

// GameState interface pour les états de jeu
type GameState interface {
	Enter()
	Update(deltaTime time.Duration) error
	Render(renderer *rendering.Renderer) error
	Exit()
	GetType() core.GameStateType
}

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	WindowWidth() int
	WindowHeight() int
	IsDebugEnabled() bool
}

// StateManager gère les transitions entre les états
type StateManager struct {
	states       map[core.GameStateType]GameState
	currentState GameState
	nextState    core.GameStateType
	changing     bool
}

// NewStateManager crée un nouveau gestionnaire d'états
func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[core.GameStateType]GameState),
	}
}

// AddState ajoute un état au gestionnaire
func (sm *StateManager) AddState(stateType core.GameStateType, state GameState) {
	sm.states[stateType] = state
}

// ChangeState change l'état actuel
func (sm *StateManager) ChangeState(stateType core.GameStateType) {
	if state, exists := sm.states[stateType]; exists {
		if sm.currentState != nil {
			sm.currentState.Exit()
		}
		sm.currentState = state
		sm.currentState.Enter()
		fmt.Printf("État changé vers: %s\n", stateType)
	}
}

// Update met à jour l'état actuel
func (sm *StateManager) Update(deltaTime time.Duration) error {
	if sm.currentState != nil {
		return sm.currentState.Update(deltaTime)
	}
	return nil
}

// Render rend l'état actuel
func (sm *StateManager) Render(renderer *rendering.Renderer) error {
	if sm.currentState != nil {
		return sm.currentState.Render(renderer)
	}
	return nil
}

// GetCurrentStateType retourne le type de l'état actuel
func (sm *StateManager) GetCurrentStateType() core.GameStateType {
	if sm.currentState != nil {
		return sm.currentState.GetType()
	}
	return ""
}

// États spécifiques (stubs)

// MenuState état du menu principal
type MenuState struct {
	config      GameConfig
	uiManager   *ui.UIManager
	saveManager *save.SaveManager
}

func NewMenuState(config GameConfig, uiManager *ui.UIManager, saveManager *save.SaveManager) *MenuState {
	return &MenuState{
		config:      config,
		uiManager:   uiManager,
		saveManager: saveManager,
	}
}

func (ms *MenuState) Enter()                               { fmt.Println("Entrée dans MenuState") }
func (ms *MenuState) Update(deltaTime time.Duration) error { return nil }
func (ms *MenuState) Render(renderer *rendering.Renderer) error {
	// Dessiner un rectangle pour indiquer qu'on est dans le menu
	renderer.DrawRectangle(
		core.Rectangle{X: 100, Y: 100, Width: 200, Height: 50},
		core.ColorBlue,
		true,
	)
	renderer.DrawText("MENU PRINCIPAL", core.Vector2{X: 120, Y: 120}, core.ColorWhite)
	return nil
}
func (ms *MenuState) Exit()                       { fmt.Println("Sortie de MenuState") }
func (ms *MenuState) GetType() core.GameStateType { return core.StateMenu }

// GameplayState état de jeu
type GameplayState struct {
	config        GameConfig
	world         *world.World
	systemManager *systems.SystemsManager
	uiManager     *ui.UIManager
	inputManager  *input.InputManagerImpl
}

func NewGameplayState(config GameConfig, world *world.World, systemManager *systems.SystemsManager, uiManager *ui.UIManager, inputManager *input.InputManagerImpl) *GameplayState {
	return &GameplayState{
		config:        config,
		world:         world,
		systemManager: systemManager,
		uiManager:     uiManager,
		inputManager:  inputManager,
	}
}

func (gs *GameplayState) Enter()                               { fmt.Println("Entrée dans GameplayState") }
func (gs *GameplayState) Update(deltaTime time.Duration) error { return nil }
func (gs *GameplayState) Render(renderer *rendering.Renderer) error {
	// Dessiner un rectangle pour indiquer qu'on est en jeu
	renderer.DrawRectangle(
		core.Rectangle{X: 100, Y: 200, Width: 200, Height: 50},
		core.ColorGreen,
		true,
	)
	renderer.DrawText("EN JEU", core.Vector2{X: 150, Y: 220}, core.ColorWhite)
	return nil
}
func (gs *GameplayState) Exit()                       { fmt.Println("Sortie de GameplayState") }
func (gs *GameplayState) GetType() core.GameStateType { return core.StateGameplay }

// PauseState état de pause
type PauseState struct {
	config    GameConfig
	uiManager *ui.UIManager
}

func NewPauseState(config GameConfig, uiManager *ui.UIManager) *PauseState {
	return &PauseState{
		config:    config,
		uiManager: uiManager,
	}
}

func (ps *PauseState) Enter()                               { fmt.Println("Entrée dans PauseState") }
func (ps *PauseState) Update(deltaTime time.Duration) error { return nil }
func (ps *PauseState) Render(renderer *rendering.Renderer) error {
	// Dessiner un rectangle pour indiquer qu'on est en pause
	renderer.DrawRectangle(
		core.Rectangle{X: 100, Y: 300, Width: 200, Height: 50},
		core.ColorYellow,
		true,
	)
	renderer.DrawText("PAUSE", core.Vector2{X: 160, Y: 320}, core.ColorBlack)
	return nil
}
func (ps *PauseState) Exit()                       { fmt.Println("Sortie de PauseState") }
func (ps *PauseState) GetType() core.GameStateType { return core.StatePause }

// InventoryState état d'inventaire
type InventoryState struct {
	config    GameConfig
	uiManager *ui.UIManager
}

func NewInventoryState(config GameConfig, uiManager *ui.UIManager) *InventoryState {
	return &InventoryState{
		config:    config,
		uiManager: uiManager,
	}
}

func (is *InventoryState) Enter()                               { fmt.Println("Entrée dans InventoryState") }
func (is *InventoryState) Update(deltaTime time.Duration) error { return nil }
func (is *InventoryState) Render(renderer *rendering.Renderer) error {
	renderer.DrawRectangle(
		core.Rectangle{X: 100, Y: 400, Width: 200, Height: 50},
		core.ColorMagenta,
		true,
	)
	renderer.DrawText("INVENTAIRE", core.Vector2{X: 130, Y: 420}, core.ColorWhite)
	return nil
}
func (is *InventoryState) Exit()                       { fmt.Println("Sortie d'InventoryState") }
func (is *InventoryState) GetType() core.GameStateType { return core.StateInventory }
