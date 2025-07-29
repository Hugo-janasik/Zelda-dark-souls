// internal/input/input_manager.go - Gestionnaire d'entrées corrigé
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	WindowWidth() int
	WindowHeight() int
}

// InputAction représente une action de jeu (local pour éviter cycles)
type InputAction int

const (
	ActionMoveUp InputAction = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionAttack
	ActionBlock
	ActionRoll
	ActionParry
	ActionInteract
	ActionPickup
	ActionUse
	ActionInventory
	ActionMap
	ActionPause
	ActionMenu
	ActionConfirm
	ActionCancel
	ActionCastSpell
	ActionQuickSlot1
	ActionQuickSlot2
	ActionQuickSlot3
	ActionQuickSlot4
	ActionCameraReset
	ActionCameraZoomIn
	ActionCameraZoomOut
)

// InputManagerImpl implémentation concrète du gestionnaire d'entrées
type InputManagerImpl struct {
	config               GameConfig
	keyPressed           map[ebiten.Key]bool
	keyJustPressed       map[ebiten.Key]bool
	keyJustReleased      map[ebiten.Key]bool
	mouseX, mouseY       int
	mousePressed         map[int]bool
	windowCloseRequested bool
}

// NewInputManager crée un nouveau gestionnaire d'entrées
func NewInputManager(config GameConfig) *InputManagerImpl {
	return &InputManagerImpl{
		config:          config,
		keyPressed:      make(map[ebiten.Key]bool),
		keyJustPressed:  make(map[ebiten.Key]bool),
		keyJustReleased: make(map[ebiten.Key]bool),
		mousePressed:    make(map[int]bool),
	}
}

// Update met à jour les entrées
func (im *InputManagerImpl) Update() {
	// Mise à jour des touches
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		pressed := ebiten.IsKeyPressed(key)
		wasPressed := im.keyPressed[key]

		im.keyJustPressed[key] = pressed && !wasPressed
		im.keyJustReleased[key] = !pressed && wasPressed
		im.keyPressed[key] = pressed
	}

	// Mise à jour de la souris
	im.mouseX, im.mouseY = ebiten.CursorPosition()

	// Vérifier si la fenêtre doit se fermer (stub)
	im.windowCloseRequested = false
}

// IsKeyPressed vérifie si une touche est pressée
func (im *InputManagerImpl) IsKeyPressed(key ebiten.Key) bool {
	return im.keyPressed[key]
}

// IsKeyJustPressed vérifie si une touche vient d'être pressée
func (im *InputManagerImpl) IsKeyJustPressed(key ebiten.Key) bool {
	return im.keyJustPressed[key]
}

// Méthodes pour l'interface core (avec int au lieu d'ebiten.Key)
func (im *InputManagerImpl) IsKeyCorePressed(key int) bool {
	return im.keyJustPressed[ebiten.Key(key)]
}

func (im *InputManagerImpl) IsActionCorePressed(action int) bool {
	return im.IsActionPressed(InputAction(action))
}

// IsActionPressed vérifie si une action est pressée
func (im *InputManagerImpl) IsActionPressed(action InputAction) bool {
	// Mapping pour clavier français AZERTY et international
	switch action {
	case ActionPause:
		return im.IsKeyJustPressed(ebiten.KeyEscape)
	case ActionMoveUp:
		return im.IsKeyPressed(ebiten.KeyW) || im.IsKeyPressed(ebiten.KeyZ) // W ou Z
	case ActionMoveDown:
		return im.IsKeyPressed(ebiten.KeyS) // S
	case ActionMoveLeft:
		return im.IsKeyPressed(ebiten.KeyA) || im.IsKeyPressed(ebiten.KeyQ) // A ou Q
	case ActionMoveRight:
		return im.IsKeyPressed(ebiten.KeyD) // D
	case ActionAttack:
		return im.IsKeyPressed(ebiten.KeySpace) // Espace pour attaquer
	case ActionBlock:
		return im.IsKeyPressed(ebiten.KeyShiftLeft) // Shift pour bloquer
	case ActionRoll:
		return im.IsKeyPressed(ebiten.KeyControlLeft) // Ctrl pour rouler
	default:
		return false
	}
}

// IsMovementActionPressed vérifie si une action de mouvement spécifique est pressée
func (im *InputManagerImpl) IsMovementActionPressed(action int) bool {
	switch action {
	case 0: // ActionMoveUp
		return im.IsKeyPressed(ebiten.KeyW) || im.IsKeyPressed(ebiten.KeyZ)
	case 1: // ActionMoveDown
		return im.IsKeyPressed(ebiten.KeyS)
	case 2: // ActionMoveLeft
		return im.IsKeyPressed(ebiten.KeyA) || im.IsKeyPressed(ebiten.KeyQ)
	case 3: // ActionMoveRight
		return im.IsKeyPressed(ebiten.KeyD)
	default:
		return false
	}
}

// IsWindowCloseRequested vérifie si la fenêtre doit se fermer
func (im *InputManagerImpl) IsWindowCloseRequested() bool {
	return im.windowCloseRequested
}

// Interface pour systems.InputManager
func (im *InputManagerImpl) IsActionPressedSystems(action int) bool {
	switch action {
	case 0: // ActionMoveUp
		return im.IsKeyPressed(ebiten.KeyW) || im.IsKeyPressed(ebiten.KeyZ)
	case 1: // ActionMoveDown
		return im.IsKeyPressed(ebiten.KeyS)
	case 2: // ActionMoveLeft
		return im.IsKeyPressed(ebiten.KeyA) || im.IsKeyPressed(ebiten.KeyQ)
	case 3: // ActionMoveRight
		return im.IsKeyPressed(ebiten.KeyD)
	case 4: // ActionAttack
		return im.IsKeyPressed(ebiten.KeySpace)
	case 5: // ActionBlock
		return im.IsKeyPressed(ebiten.KeyShiftLeft) || im.IsKeyPressed(ebiten.KeyShiftRight)
	case 6: // ActionRoll
		return im.IsKeyPressed(ebiten.KeyC)
	case 8: // ActionInteract
		return im.IsKeyPressed(ebiten.KeyE)
	default:
		return false
	}
}

// IsKeyJustPressedSystems pour l'interface systems
func (im *InputManagerImpl) IsKeyJustPressedSystems(key int) bool {
	return im.IsKeyJustPressed(ebiten.Key(key))
}

// Constantes manquantes
const (
	KeyF12       = ebiten.KeyF12
	KeyBackQuote = ebiten.KeyBackquote
	KeyAltF4     = ebiten.KeyF4 // Simplification pour Alt+F4
)