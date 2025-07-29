// internal/core/types.go - Types de base utilisés dans tout le jeu
package core

import (
	"fmt"
	"image/color"
	"math"
	"time"
)

// ===============================
// MATH TYPES
// ===============================

// Vector2 représente un vecteur 2D
type Vector2 struct {
	X, Y float64
}

// NewVector2 crée un nouveau Vector2
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// Add additionne deux vecteurs
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub soustrait un vecteur
func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul multiplie par un scalaire
func (v Vector2) Mul(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Length calcule la longueur du vecteur
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize normalise le vecteur
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{0, 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// Distance calcule la distance entre deux points
func (v Vector2) Distance(other Vector2) float64 {
	return v.Sub(other).Length()
}

// Dot produit scalaire
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Rectangle représente un rectangle
type Rectangle struct {
	X, Y, Width, Height float64
}

// NewRectangle crée un nouveau Rectangle
func NewRectangle(x, y, width, height float64) Rectangle {
	return Rectangle{X: x, Y: y, Width: width, Height: height}
}

// Contains vérifie si un point est dans le rectangle
func (r Rectangle) Contains(point Vector2) bool {
	return point.X >= r.X && point.X <= r.X+r.Width &&
		point.Y >= r.Y && point.Y <= r.Y+r.Height
}

// Intersects vérifie si deux rectangles se chevauchent
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Center retourne le centre du rectangle
func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// ===============================
// GAME TYPES
// ===============================

// Direction représente une direction
type Direction int

const (
	DirectionNone Direction = iota
	DirectionUp
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionUpLeft
	DirectionUpRight
	DirectionDownLeft
	DirectionDownRight
)

// String retourne la représentation string de la direction
func (d Direction) String() string {
	switch d {
	case DirectionUp:
		return "up"
	case DirectionDown:
		return "down"
	case DirectionLeft:
		return "left"
	case DirectionRight:
		return "right"
	case DirectionUpLeft:
		return "up-left"
	case DirectionUpRight:
		return "up-right"
	case DirectionDownLeft:
		return "down-left"
	case DirectionDownRight:
		return "down-right"
	default:
		return "none"
	}
}

// ToVector2 convertit la direction en Vector2 normalisé
func (d Direction) ToVector2() Vector2 {
	switch d {
	case DirectionUp:
		return Vector2{0, -1}
	case DirectionDown:
		return Vector2{0, 1}
	case DirectionLeft:
		return Vector2{-1, 0}
	case DirectionRight:
		return Vector2{1, 0}
	case DirectionUpLeft:
		return Vector2{-1, -1}.Normalize()
	case DirectionUpRight:
		return Vector2{1, -1}.Normalize()
	case DirectionDownLeft:
		return Vector2{-1, 1}.Normalize()
	case DirectionDownRight:
		return Vector2{1, 1}.Normalize()
	default:
		return Vector2{0, 0}
	}
}

// ===============================
// INPUT TYPES
// ===============================

// Key représente une touche du clavier
type Key int

// InputAction représente une action de jeu
type InputAction int

const (
	// Actions de mouvement
	ActionMoveUp InputAction = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight

	// Actions de combat
	ActionAttack
	ActionBlock
	ActionRoll
	ActionParry

	// Actions d'interaction
	ActionInteract
	ActionPickup
	ActionUse

	// Actions de menu/UI
	ActionInventory
	ActionMap
	ActionPause
	ActionMenu
	ActionConfirm
	ActionCancel

	// Actions spéciales
	ActionCastSpell
	ActionQuickSlot1
	ActionQuickSlot2
	ActionQuickSlot3
	ActionQuickSlot4

	// Actions de camera
	ActionCameraReset
	ActionCameraZoomIn
	ActionCameraZoomOut
)

// MouseButton représente un bouton de souris
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

// GamepadButton représente un bouton de manette
type GamepadButton int

// GamepadAxis représente un axe de manette
type GamepadAxis int

// ===============================
// ENTITY TYPES
// ===============================

// EntityID représente l'identifiant unique d'une entité
type EntityID uint32

// ComponentType représente le type d'un composant
type ComponentType string

// SystemType représente le type d'un système
type SystemType string

// ===============================
// GAME STATE TYPES
// ===============================

// GameStateType représente le type d'un état de jeu
type GameStateType string

const (
	StateMenu      GameStateType = "menu"
	StateGameplay  GameStateType = "gameplay"
	StatePause     GameStateType = "pause"
	StateInventory GameStateType = "inventory"
	StateDialog    GameStateType = "dialog"
	StateLoading   GameStateType = "loading"
	StateSettings  GameStateType = "settings"
	StateGameOver  GameStateType = "gameover"
)

// ===============================
// COLLISION TYPES
// ===============================

// CollisionLayer représente les couches de collision
type CollisionLayer int

const (
	LayerPlayer CollisionLayer = iota
	LayerEnemy
	LayerNPC
	LayerEnvironment
	LayerProjectile
	LayerTrigger
	LayerItem
	LayerWall
)

// CollisionMask représente le masque de collision
type CollisionMask uint32

// ===============================
// RESOURCE TYPES
// ===============================

// TextureID représente l'identifiant d'une texture
type TextureID string

// SoundID représente l'identifiant d'un son
type SoundID string

// MusicID représente l'identifiant d'une musique
type MusicID string

// MapID représente l'identifiant d'une carte
type MapID string

// ===============================
// TIMER TYPES
// ===============================

// Timer représente un minuteur
type Timer struct {
	Duration   time.Duration
	Elapsed    time.Duration
	Running    bool
	Loop       bool
	OnComplete func()
}

// NewTimer crée un nouveau Timer
func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		Duration: duration,
		Running:  false,
		Loop:     false,
	}
}

// Start démarre le timer
func (t *Timer) Start() {
	t.Running = true
	t.Elapsed = 0
}

// Stop arrête le timer
func (t *Timer) Stop() {
	t.Running = false
}

// Reset remet le timer à zéro
func (t *Timer) Reset() {
	t.Elapsed = 0
}

// Update met à jour le timer
func (t *Timer) Update(dt time.Duration) {
	if !t.Running {
		return
	}

	t.Elapsed += dt

	if t.Elapsed >= t.Duration {
		if t.OnComplete != nil {
			t.OnComplete()
		}

		if t.Loop {
			t.Elapsed = 0
		} else {
			t.Running = false
		}
	}
}

// IsComplete retourne true si le timer est terminé
func (t *Timer) IsComplete() bool {
	return t.Elapsed >= t.Duration
}

// Progress retourne le progrès du timer (0.0 à 1.0)
func (t *Timer) Progress() float64 {
	if t.Duration == 0 {
		return 1.0
	}
	progress := float64(t.Elapsed) / float64(t.Duration)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// ===============================
// UTILITY TYPES
// ===============================

// Color représente une couleur RGBA
type Color struct {
	R, G, B, A uint8
}

// NewColor crée une nouvelle couleur
func NewColor(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// ToEbitenColor convertit vers une couleur Ebiten
func (c Color) ToEbitenColor() color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

// Predefined colors
var (
	ColorWhite   = Color{255, 255, 255, 255}
	ColorBlack   = Color{0, 0, 0, 255}
	ColorRed     = Color{255, 0, 0, 255}
	ColorGreen   = Color{0, 255, 0, 255}
	ColorBlue    = Color{0, 0, 255, 255}
	ColorYellow  = Color{255, 255, 0, 255}
	ColorMagenta = Color{255, 0, 255, 255}
	ColorCyan    = Color{0, 255, 255, 255}
	ColorGray    = Color{128, 128, 128, 255}
)

// ChunkCoord représente les coordonnées d'un chunk de carte
type ChunkCoord struct {
	X, Y int
}

// NewChunkCoord crée une nouvelle coordonnée de chunk
func NewChunkCoord(x, y int) ChunkCoord {
	return ChunkCoord{X: x, Y: y}
}

// ===============================
// CONSTANTS
// ===============================

const (
	// Constantes de jeu
	TileSize      = 32
	ChunkSize     = 16 // 16x16 tiles par chunk
	MaxEntities   = 10000
	MaxComponents = 64

	// Constantes physiques
	Gravity         = 9.81
	DefaultFriction = 0.8

	// Constantes de gameplay
	DefaultPlayerSpeed   = 100.0 // pixels par seconde
	DefaultPlayerHealth  = 100
	DefaultPlayerStamina = 100.0

	// Constantes d'animation
	DefaultAnimationFPS = 12
)

// ===============================
// ERROR TYPES
// ===============================

// GameError représente une erreur de jeu
type GameError struct {
	Code    string
	Message string
	Cause   error
}

// Error implémente l'interface error
func (e *GameError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewGameError crée une nouvelle erreur de jeu
func NewGameError(code, message string, cause error) *GameError {
	return &GameError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// ===============================
// MATH UTILITIES
// ===============================

// Cos retourne le cosinus de l'angle en radians
func Cos(angle float64) float64 {
	return math.Cos(angle)
}

// Sin retourne le sinus de l'angle en radians
func Sin(angle float64) float64 {
	return math.Sin(angle)
}

// DegreesToRadians convertit des degrés en radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// RadiansToDegrees convertit des radians en degrés
func RadiansToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// Lerp effectue une interpolation linéaire entre a et b
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Clamp limite une valeur entre min et max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
