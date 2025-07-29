// internal/ecs/components/player_components.go - Composants ECS sans imports cycliques
package components

import (
	"time"
)

// ===============================
// TYPES DE BASE (copiés pour éviter les cycles)
// ===============================

// Vector2 représente un vecteur 2D
type Vector2 struct {
	X, Y float64
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
	return v.X*v.X + v.Y*v.Y // sqrt omis pour performance, calculé ailleurs si nécessaire
}

// Normalize normalise le vecteur
func (v Vector2) Normalize() Vector2 {
	length := v.X*v.X + v.Y*v.Y
	if length == 0 {
		return Vector2{0, 0}
	}
	// Approximation rapide de sqrt
	invLength := 1.0 / (length * 0.5) // Approximation
	return Vector2{X: v.X * invLength, Y: v.Y * invLength}
}

// Rectangle représente un rectangle
type Rectangle struct {
	X, Y, Width, Height float64
}

// Color représente une couleur RGBA
type Color struct {
	R, G, B, A uint8
}

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
		return Vector2{-0.707, -0.707} // Approximation de 1/sqrt(2)
	case DirectionUpRight:
		return Vector2{0.707, -0.707}
	case DirectionDownLeft:
		return Vector2{-0.707, 0.707}
	case DirectionDownRight:
		return Vector2{0.707, 0.707}
	default:
		return Vector2{0, 0}
	}
}

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

// Couleurs prédéfinies
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

// ===============================
// COMPOSANTS DE BASE
// ===============================

// PositionComponent représente la position dans le monde
type PositionComponent struct {
	Position     Vector2
	LastPosition Vector2 // Pour calculer la vélocité
}

// NewPositionComponent crée un nouveau composant de position
func NewPositionComponent(x, y float64) *PositionComponent {
	pos := Vector2{X: x, Y: y}
	return &PositionComponent{
		Position:     pos,
		LastPosition: pos,
	}
}

// MovementComponent gère le mouvement
type MovementComponent struct {
	Velocity      Vector2
	Speed         float64
	MaxSpeed      float64
	Acceleration  float64
	Friction      float64
	IsMoving      bool
	Direction     Direction
	FacingDir     Direction // Direction vers laquelle regarde l'entité
}

// NewMovementComponent crée un nouveau composant de mouvement
func NewMovementComponent(speed, maxSpeed float64) *MovementComponent {
	return &MovementComponent{
		Velocity:     Vector2{X: 0, Y: 0},
		Speed:        speed,
		MaxSpeed:     maxSpeed,
		Acceleration: speed * 4, // Accélération rapide
		Friction:     speed * 6, // Friction plus élevée pour un arrêt rapide
		IsMoving:     false,
		Direction:    DirectionNone,
		FacingDir:    DirectionDown, // Direction par défaut
	}
}

// SpriteComponent gère l'affichage graphique
type SpriteComponent struct {
	TextureID    string
	SourceRect   Rectangle // Rectangle source dans la texture
	Size         Vector2   // Taille d'affichage
	Scale        Vector2   // Échelle de rendu
	Rotation     float64   // Rotation en radians
	Color        Color     // Teinte
	FlipX        bool      // Miroir horizontal
	FlipY        bool      // Miroir vertical
	Visible      bool      // Visibilité
	Layer        int       // Couche de rendu (plus haut = devant)
	Offset       Vector2   // Décalage par rapport à la position
}

// NewSpriteComponent crée un nouveau composant de sprite
func NewSpriteComponent(textureID string, width, height float64) *SpriteComponent {
	return &SpriteComponent{
		TextureID:  textureID,
		SourceRect: Rectangle{X: 0, Y: 0, Width: width, Height: height},
		Size:       Vector2{X: width, Y: height},
		Scale:      Vector2{X: 1.0, Y: 1.0},
		Rotation:   0,
		Color:      ColorWhite,
		FlipX:      false,
		FlipY:      false,
		Visible:    true,
		Layer:      0,
		Offset:     Vector2{X: 0, Y: 0},
	}
}

// ===============================
// COMPOSANTS D'ANIMATION
// ===============================

// AnimationFrame représente une frame d'animation
type AnimationFrame struct {
	SourceRect Rectangle // Rectangle dans la texture
	Duration   time.Duration  // Durée de la frame
}

// Animation représente une séquence d'animation
type Animation struct {
	Name       string
	Frames     []AnimationFrame
	Loop       bool
	PlayRate   float64 // Multiplicateur de vitesse (1.0 = normal)
}

// AnimationComponent gère les animations de sprites
type AnimationComponent struct {
	Animations      map[string]*Animation
	CurrentAnim     string
	CurrentFrame    int
	ElapsedTime     time.Duration
	Playing         bool
	PlayRate        float64
	OnComplete      func() // Callback à la fin de l'animation
}

// NewAnimationComponent crée un nouveau composant d'animation
func NewAnimationComponent() *AnimationComponent {
	return &AnimationComponent{
		Animations:   make(map[string]*Animation),
		CurrentAnim:  "",
		CurrentFrame: 0,
		ElapsedTime:  0,
		Playing:      false,
		PlayRate:     1.0,
	}
}

// AddAnimation ajoute une animation
func (ac *AnimationComponent) AddAnimation(name string, animation *Animation) {
	ac.Animations[name] = animation
}

// Play joue une animation
func (ac *AnimationComponent) Play(name string) {
	if animation, exists := ac.Animations[name]; exists {
		ac.CurrentAnim = name
		ac.CurrentFrame = 0
		ac.ElapsedTime = 0
		ac.Playing = true
		ac.PlayRate = animation.PlayRate
	}
}

// Stop arrête l'animation
func (ac *AnimationComponent) Stop() {
	ac.Playing = false
}

// IsPlaying retourne si une animation est en cours
func (ac *AnimationComponent) IsPlaying(name string) bool {
	return ac.Playing && ac.CurrentAnim == name
}

// GetCurrentFrame retourne la frame actuelle
func (ac *AnimationComponent) GetCurrentFrame() *AnimationFrame {
	if ac.CurrentAnim == "" {
		return nil
	}
	
	animation := ac.Animations[ac.CurrentAnim]
	if animation == nil || ac.CurrentFrame >= len(animation.Frames) {
		return nil
	}
	
	return &animation.Frames[ac.CurrentFrame]
}

// ===============================
// COMPOSANTS DE COLLISION
// ===============================

// ColliderComponent représente la zone de collision
type ColliderComponent struct {
	Bounds       Rectangle
	Offset       Vector2
	Layer        CollisionLayer
	Mask         CollisionMask
	IsTrigger    bool
	Enabled      bool
}

// NewColliderComponent crée un nouveau composant de collision
func NewColliderComponent(width, height float64, layer CollisionLayer) *ColliderComponent {
	return &ColliderComponent{
		Bounds:    Rectangle{X: 0, Y: 0, Width: width, Height: height},
		Offset:    Vector2{X: 0, Y: 0},
		Layer:     layer,
		Mask:      0xFFFFFFFF, // Collisionne avec tout par défaut
		IsTrigger: false,
		Enabled:   true,
	}
}

// GetWorldBounds retourne les limites dans le monde
func (cc *ColliderComponent) GetWorldBounds(position Vector2) Rectangle {
	return Rectangle{
		X:      position.X + cc.Offset.X - cc.Bounds.Width/2,
		Y:      position.Y + cc.Offset.Y - cc.Bounds.Height/2,
		Width:  cc.Bounds.Width,
		Height: cc.Bounds.Height,
	}
}

// ===============================
// COMPOSANTS SPÉCIFIQUES AU JOUEUR
// ===============================

// PlayerComponent marque une entité comme étant le joueur
type PlayerComponent struct {
	// Stats de base
	Health          int
	MaxHealth       int
	Stamina         float64
	MaxStamina      float64
	StaminaRegen    float64
	
	// Combat
	AttackPower     int
	Defense         int
	CriticalChance  float64
	
	// Progression
	Level           int
	Experience      int
	ExperienceToNext int
	
	// États
	InvulnTime      time.Duration // Temps d'invulnérabilité restant
	Stunned         bool
	StunTime        time.Duration
	
	// Statistiques de jeu
	PlayTime        time.Duration
	EnemiesKilled   int
	ItemsCollected  int
}

// NewPlayerComponent crée un nouveau composant joueur
func NewPlayerComponent() *PlayerComponent {
	return &PlayerComponent{
		Health:           100,
		MaxHealth:        100,
		Stamina:          100.0,
		MaxStamina:       100.0,
		StaminaRegen:     25.0, // Points par seconde
		AttackPower:      10,
		Defense:          5,
		CriticalChance:   0.05, // 5%
		Level:            1,
		Experience:       0,
		ExperienceToNext: 100,
		InvulnTime:       0,
		Stunned:          false,
		StunTime:         0,
		PlayTime:         0,
		EnemiesKilled:    0,
		ItemsCollected:   0,
	}
}

// IsAlive retourne si le joueur est vivant
func (pc *PlayerComponent) IsAlive() bool {
	return pc.Health > 0
}

// TakeDamage inflige des dégâts au joueur
func (pc *PlayerComponent) TakeDamage(damage int) bool {
	if pc.InvulnTime > 0 {
		return false // Invulnérable
	}
	
	actualDamage := damage - pc.Defense
	if actualDamage < 1 {
		actualDamage = 1 // Au minimum 1 dégât
	}
	
	pc.Health -= actualDamage
	if pc.Health < 0 {
		pc.Health = 0
	}
	
	// Temps d'invulnérabilité après dégâts
	pc.InvulnTime = time.Millisecond * 1000 // 1 seconde
	
	return true
}

// Heal soigne le joueur
func (pc *PlayerComponent) Heal(amount int) {
	pc.Health += amount
	if pc.Health > pc.MaxHealth {
		pc.Health = pc.MaxHealth
	}
}

// UseStamina consomme de la stamina
func (pc *PlayerComponent) UseStamina(amount float64) bool {
	if pc.Stamina >= amount {
		pc.Stamina -= amount
		return true
	}
	return false
}

// RegenerateStamina régénère la stamina
func (pc *PlayerComponent) RegenerateStamina(deltaTime time.Duration) {
	if pc.Stamina < pc.MaxStamina {
		regen := pc.StaminaRegen * deltaTime.Seconds()
		pc.Stamina += regen
		if pc.Stamina > pc.MaxStamina {
			pc.Stamina = pc.MaxStamina
		}
	}
}

// Update met à jour les timers du joueur
func (pc *PlayerComponent) Update(deltaTime time.Duration) {
	// Réduire le temps d'invulnérabilité
	if pc.InvulnTime > 0 {
		pc.InvulnTime -= deltaTime
		if pc.InvulnTime < 0 {
			pc.InvulnTime = 0
		}
	}
	
	// Réduire le temps de stun
	if pc.Stunned && pc.StunTime > 0 {
		pc.StunTime -= deltaTime
		if pc.StunTime <= 0 {
			pc.Stunned = false
			pc.StunTime = 0
		}
	}
	
	// Régénération de stamina
	pc.RegenerateStamina(deltaTime)
	
	// Compteur de temps de jeu
	pc.PlayTime += deltaTime
}

// ===============================
// COMPOSANT DE CONTRÔLE INPUT
// ===============================

// InputComponent gère les entrées pour cette entité
type InputComponent struct {
	Enabled         bool
	ControllerID    int  // Pour le multijoueur futur
	
	// Actions actuelles
	MoveUp          bool
	MoveDown        bool
	MoveLeft        bool
	MoveRight       bool
	Attack          bool
	Block           bool
	Roll            bool
	Interact        bool
	UseItem         bool
	
	// Actions "just pressed" (frame unique)
	AttackJustPressed   bool
	BlockJustPressed    bool
	RollJustPressed     bool
	InteractJustPressed bool
	UseItemJustPressed  bool
}

// NewInputComponent crée un nouveau composant d'entrée
func NewInputComponent() *InputComponent {
	return &InputComponent{
		Enabled:      true,
		ControllerID: 0,
	}
}

// Reset remet toutes les actions à false (appelé chaque frame)
func (ic *InputComponent) Reset() {
	ic.MoveUp = false
	ic.MoveDown = false
	ic.MoveLeft = false
	ic.MoveRight = false
	ic.Attack = false
	ic.Block = false
	ic.Roll = false
	ic.Interact = false
	ic.UseItem = false
	
	// Reset "just pressed"
	ic.AttackJustPressed = false
	ic.BlockJustPressed = false
	ic.RollJustPressed = false
	ic.InteractJustPressed = false
	ic.UseItemJustPressed = false
}

// GetMovementVector retourne le vecteur de mouvement normalisé
func (ic *InputComponent) GetMovementVector() Vector2 {
	movement := Vector2{X: 0, Y: 0}
	
	if ic.MoveLeft {
		movement.X -= 1
	}
	if ic.MoveRight {
		movement.X += 1
	}
	if ic.MoveUp {
		movement.Y -= 1
	}
	if ic.MoveDown {
		movement.Y += 1
	}
	
	// Normaliser pour éviter que la diagonale soit plus rapide
	if movement.X != 0 || movement.Y != 0 {
		return movement.Normalize()
	}
	
	return movement
}

// IsMoving retourne si le joueur essaie de bouger
func (ic *InputComponent) IsMoving() bool {
	return ic.MoveUp || ic.MoveDown || ic.MoveLeft || ic.MoveRight
}