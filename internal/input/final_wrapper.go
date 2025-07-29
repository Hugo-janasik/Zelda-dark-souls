// internal/input/final_wrapper.go - Wrapper final sans conflits de types
package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// FinalInputWrapper wrapper final sans imports cycliques ni conflits
type FinalInputWrapper struct {
	inputManager  *InputManagerImpl
	coreGame      interface{}
	lastFrameKeys map[ebiten.Key]bool
	
	// État des actions pour éviter les répétitions
	lastPauseState     bool
	lastInstructState  bool
}

// NewFinalInputWrapper crée un wrapper final
func NewFinalInputWrapper(im *InputManagerImpl) *FinalInputWrapper {
	return &FinalInputWrapper{
		inputManager:  im,
		lastFrameKeys: make(map[ebiten.Key]bool),
	}
}

// SetCoreGame injecte le jeu core
func (w *FinalInputWrapper) SetCoreGame(cg interface{}) {
	w.coreGame = cg
	fmt.Println("CoreGame injecté dans FinalInputWrapper")
}

// Update met à jour et traite les actions
func (w *FinalInputWrapper) Update() {
	w.inputManager.Update()
	w.updateMouseInput()
	w.handleGlobalActions()
	w.updateLastFrameKeys()
}

// updateMouseInput met à jour les entrées souris - SOLUTION SIMPLE
func (w *FinalInputWrapper) updateMouseInput() {
	if w.coreGame == nil {
		return
	}

	// Interface pour obtenir le StateManager
	type StateManagerProvider interface {
		GetBuiltinStateManager() interface{}
	}

	if provider, ok := w.coreGame.(StateManagerProvider); ok {
		stateManager := provider.GetBuiltinStateManager()
		
		// Utiliser reflection pour appeler UpdateMouseInput
		if sm, ok := stateManager.(interface {
			UpdateMouseInput(int, int, bool)
		}); ok {
			mouseX, mouseY := ebiten.CursorPosition()
			mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
			sm.UpdateMouseInput(mouseX, mouseY, mousePressed)
			
			// Debug pour vérifier que la souris est bien détectée
			if mousePressed {
				fmt.Printf("Souris cliquée à (%d, %d)\n", mouseX, mouseY)
			}
		} else {
			// Debug pour voir le type réel
			fmt.Printf("Type StateManager reçu: %T\n", stateManager)
			fmt.Println("Impossible de caster vers interface UpdateMouseInput")
		}
	} else {
		fmt.Println("StateManagerProvider non trouvé")
	}
}

// handleGlobalActions traite les actions globales (ESC, I, etc.)
func (w *FinalInputWrapper) handleGlobalActions() {
	if w.coreGame == nil {
		return
	}

	// Interface pour obtenir le StateManager
	type StateManagerProvider interface {
		GetBuiltinStateManager() interface{}
	}

	var stateManager interface{}
	if provider, ok := w.coreGame.(StateManagerProvider); ok {
		stateManager = provider.GetBuiltinStateManager()
	}

	if stateManager == nil {
		return
	}

	// ESC - Menu/Pause
	escPressed := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if escPressed && !w.lastPauseState {
		fmt.Println("ESC pressé - traitement global...")
		
		if sm, ok := stateManager.(interface {
			GetCurrentStateType() interface{}
			ChangeState(interface{})
		}); ok {
			currentState := sm.GetCurrentStateType()
			stateStr := fmt.Sprintf("%v", currentState)
			
			switch stateStr {
			case "gameplay":
				sm.ChangeState("menu")
				fmt.Println("Retour au menu depuis le jeu")
			case "pause":
				sm.ChangeState("gameplay")
				fmt.Println("Reprise du jeu")
			case "menu":
				// Déjà dans le menu, ne rien faire
			}
		}
	}
	w.lastPauseState = escPressed

	// I - Toggle instructions (seulement en gameplay)
	iPressed := ebiten.IsKeyPressed(ebiten.KeyI)
	if iPressed && !w.lastInstructState {
		if sm, ok := stateManager.(interface {
			GetCurrentStateType() interface{}
			ToggleInstructions()
		}); ok {
			currentState := sm.GetCurrentStateType()
			if fmt.Sprintf("%v", currentState) == "gameplay" {
				sm.ToggleInstructions()
				fmt.Println("Toggle instructions")
			}
		}
	}
	w.lastInstructState = iPressed
}

// ===============================
// INTERFACE CORE.INPUTMANAGER
// ===============================

func (w *FinalInputWrapper) IsKeyJustPressed(key int) bool {
	return w.wasKeyJustPressed(ebiten.Key(key))
}

func (w *FinalInputWrapper) IsActionPressed(action int) bool {
	return w.inputManager.IsActionPressedSystems(action)
}

func (w *FinalInputWrapper) IsWindowCloseRequested() bool {
	return w.inputManager.IsWindowCloseRequested()
}

// ===============================
// INTERFACE SYSTEMS.INPUTMANAGER
// ===============================

func (w *FinalInputWrapper) IsActionPressedSystems(action int) bool {
	return w.inputManager.IsActionPressedSystems(action)
}

func (w *FinalInputWrapper) IsKeyJustPressedSystems(key int) bool {
	return w.wasKeyJustPressed(ebiten.Key(key))
}

// ===============================
// MÉTHODES UTILITAIRES
// ===============================

// wasKeyJustPressed vérifie si une touche vient d'être pressée cette frame
func (w *FinalInputWrapper) wasKeyJustPressed(key ebiten.Key) bool {
	currentlyPressed := ebiten.IsKeyPressed(key)
	wasPressed := w.lastFrameKeys[key]
	return currentlyPressed && !wasPressed
}

// updateLastFrameKeys met à jour l'état des touches de la frame précédente
func (w *FinalInputWrapper) updateLastFrameKeys() {
	// Sauvegarder l'état de toutes les touches importantes
	keysToTrack := []ebiten.Key{
		ebiten.KeyEscape,
		ebiten.KeyI,
		ebiten.KeySpace,
		ebiten.KeyC,
		ebiten.KeyE,
		ebiten.KeyW, ebiten.KeyZ,
		ebiten.KeyS,
		ebiten.KeyA, ebiten.KeyQ,
		ebiten.KeyD,
		ebiten.KeyShiftLeft, ebiten.KeyShiftRight,
	}

	for _, key := range keysToTrack {
		w.lastFrameKeys[key] = ebiten.IsKeyPressed(key)
	}
}

// GetMovementVector retourne le vecteur de mouvement actuel
func (w *FinalInputWrapper) GetMovementVector() (float64, float64) {
	var x, y float64
	
	if w.IsActionPressed(2) { // Left
		x -= 1
	}
	if w.IsActionPressed(3) { // Right
		x += 1
	}
	if w.IsActionPressed(0) { // Up
		y -= 1
	}
	if w.IsActionPressed(1) { // Down
		y += 1
	}
	
	return x, y
}

// IsMoving retourne si le joueur bouge actuellement
func (w *FinalInputWrapper) IsMoving() bool {
	return w.IsActionPressed(0) || w.IsActionPressed(1) || 
		   w.IsActionPressed(2) || w.IsActionPressed(3)
}