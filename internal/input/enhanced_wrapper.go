// internal/input/enhanced_wrapper.go - Wrapper étendu avec actions
package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// CoreGame interface minimale pour éviter les cycles
type CoreGame interface {
	GetBuiltinStateManager() interface{}
}

// EnhancedInputWrapper wrapper avec logique d'actions
type EnhancedInputWrapper struct {
	inputManager  *InputManager
	coreGame      CoreGame
	lastFrameKeys map[ebiten.Key]bool
}

// NewEnhancedInputWrapper crée un wrapper étendu
func NewEnhancedInputWrapper(im *InputManager) *EnhancedInputWrapper {
	return &EnhancedInputWrapper{
		inputManager:  im,
		lastFrameKeys: make(map[ebiten.Key]bool),
	}
}

// SetCoreGame injecte le jeu core
func (w *EnhancedInputWrapper) SetCoreGame(cg CoreGame) {
	w.coreGame = cg
	fmt.Println("CoreGame injecté dans InputWrapper")
}

// Update met à jour et traite les actions
func (w *EnhancedInputWrapper) Update() {
	w.inputManager.Update()

	// Mettre à jour les entrées souris pour le StateManager
	if w.coreGame != nil {
		stateManager := w.coreGame.GetBuiltinStateManager()
		if sm, ok := stateManager.(interface {
			UpdateMouseInput(mouseX, mouseY int, mousePressed bool)
		}); ok {
			mouseX, mouseY := ebiten.CursorPosition()
			mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			sm.UpdateMouseInput(mouseX, mouseY, mousePressed)
		}
	}

	w.handleActions()
	w.updateLastFrameKeys()
}

// handleActions traite les actions spéciales
func (w *EnhancedInputWrapper) handleActions() {
	// Obtenir le StateManager depuis le core
	var stateManager interface{}
	if w.coreGame != nil {
		stateManager = w.coreGame.GetBuiltinStateManager()
	}

	// ESC - Retour au menu ou toggle
	if w.wasKeyJustPressed(ebiten.KeyEscape) {
		fmt.Println("ESC pressé - traitement...")
		if stateManager != nil {
			if sm, ok := stateManager.(interface {
				GetCurrentStateType() interface{}
				ChangeState(interface{})
			}); ok {
				currentState := sm.GetCurrentStateType()
				fmt.Printf("État actuel lors d'ESC: %v\n", currentState)
				if fmt.Sprintf("%v", currentState) == "gameplay" {
					sm.ChangeState("menu")
					fmt.Println("Retour au menu")
				}
			}
		}
	}

	// I - Toggle instructions (seulement en gameplay)
	if w.wasKeyJustPressed(ebiten.KeyI) {
		fmt.Println("I pressé - traitement...")
		if stateManager != nil {
			if sm, ok := stateManager.(interface {
				GetCurrentStateType() interface{}
				ToggleInstructions()
			}); ok {
				currentState := sm.GetCurrentStateType()
				if fmt.Sprintf("%v", currentState) == "gameplay" {
					sm.ToggleInstructions()
					fmt.Println("I pressé - Toggle instructions")
				} else {
					fmt.Printf("I ignoré car état = %v\n", currentState)
				}
			}
		}
	}

	// Test des touches de mouvement
	w.testMovementKeys()
}

// testMovementKeys teste et affiche les touches de mouvement
func (w *EnhancedInputWrapper) testMovementKeys() {
	movements := []struct {
		keys []ebiten.Key
		name string
	}{
		{[]ebiten.Key{ebiten.KeyZ, ebiten.KeyW}, "Haut"},
		{[]ebiten.Key{ebiten.KeyS}, "Bas"},
		{[]ebiten.Key{ebiten.KeyQ, ebiten.KeyA}, "Gauche"},
		{[]ebiten.Key{ebiten.KeyD}, "Droite"},
	}

	for _, movement := range movements {
		for _, key := range movement.keys {
			if w.wasKeyJustPressed(key) {
				fmt.Printf("Mouvement: %s (touche %v)\n", movement.name, key)
			}
		}
	}
}

// wasKeyJustPressed vérifie si une touche vient d'être pressée
func (w *EnhancedInputWrapper) wasKeyJustPressed(key ebiten.Key) bool {
	currentlyPressed := ebiten.IsKeyPressed(key)
	wasPressed := w.lastFrameKeys[key]
	return currentlyPressed && !wasPressed
}

// updateLastFrameKeys met à jour l'état des touches de la frame précédente
func (w *EnhancedInputWrapper) updateLastFrameKeys() {
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		w.lastFrameKeys[key] = ebiten.IsKeyPressed(key)
	}
}

// Interface core.InputManager
func (w *EnhancedInputWrapper) IsKeyJustPressed(key int) bool {
	return w.wasKeyJustPressed(ebiten.Key(key))
}

func (w *EnhancedInputWrapper) IsActionPressed(action int) bool {
	// Utiliser directement l'InputManager pour les actions de mouvement
	return w.inputManager.IsMovementActionPressed(action)
}

func (w *EnhancedInputWrapper) IsWindowCloseRequested() bool {
	return w.inputManager.IsWindowCloseRequested()
}
