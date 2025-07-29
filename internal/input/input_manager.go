// internal/input/input_manager.go - Gestionnaire d'entrées
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	WindowWidth() int
	WindowHeight() int
}

// InputAction représente une action de jeu (copié de core pour éviter cycle)
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

type InputManager struct {
	config               GameConfig
	keyPressed           map[ebiten.Key]bool
	keyJustPressed       map[ebiten.Key]bool
	keyJustReleased      map[ebiten.Key]bool
	mouseX, mouseY       int
	mousePressed         map[int]bool
	windowCloseRequested bool
}

func NewInputManager(config GameConfig) *InputManager {
	return &InputManager{
		config:          config,
		keyPressed:      make(map[ebiten.Key]bool),
		keyJustPressed:  make(map[ebiten.Key]bool),
		keyJustReleased: make(map[ebiten.Key]bool),
		mousePressed:    make(map[int]bool),
	}
}

func (im *InputManager) Update() {
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

func (im *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return im.keyPressed[key]
}

func (im *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return im.keyJustPressed[key]
}

// Méthodes pour l'interface core (avec int au lieu d'ebiten.Key)
func (im *InputManager) IsKeyCorePressed(key int) bool {
	return im.keyJustPressed[ebiten.Key(key)]
}

func (im *InputManager) IsActionCorePressed(action int) bool {
	return im.IsActionPressed(InputAction(action))
}

func (im *InputManager) IsActionPressed(action InputAction) bool {
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
func (im *InputManager) IsMovementActionPressed(action int) bool {
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

func (im *InputManager) IsWindowCloseRequested() bool {
	return im.windowCloseRequested
}

// Constantes manquantes
const (
	KeyF12       = ebiten.KeyF12
	KeyBackQuote = ebiten.KeyBackquote
	KeyAltF4     = ebiten.KeyF4 // Simplification pour Alt+F4
)
