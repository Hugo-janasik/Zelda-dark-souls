// internal/ecs/systems/player_system.go - Système de joueur avec support des sprites
package systems

import (
	"fmt"
	"math"
	"time"
	"zelda-souls-game/internal/ecs/components"
)

// ===============================
// INTERFACES MINIMALES CORRIGÉES
// ===============================

// InputManager interface minimale pour éviter les cycles
type InputManager interface {
	IsActionPressedSystems(action int) bool
	IsKeyJustPressedSystems(key int) bool
}

// Renderer interface minimale pour le rendu (utilise les types components)
type Renderer interface {
	DrawRectangle(rect components.Rectangle, color components.Color, filled bool)
	DrawText(text string, pos components.Vector2, color components.Color)
	DrawSprite(sprite interface{}, position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, rotation float64, tint components.Color) // NOUVEAU
}

// Camera interface pour la caméra
type Camera interface {
	SetTarget(position components.Vector2)
	FollowTarget(target interface{}, speed float64, offset components.Vector2)
}

// SpriteLoader interface pour charger les sprites
type SpriteLoader interface {
	LoadPlayerSprites(assetsDir string) (interface{}, error) // Retourne *PlayerSpriteSet
}

// ===============================
// ENTITÉ JOUEUR AVEC SPRITES
// ===============================

// PlayerEntity représente l'entité joueur complète avec sprites
type PlayerEntity struct {
	// Composants
	Position      *components.PositionComponent
	Movement      *components.MovementComponent
	Sprite        *components.SpriteComponent
	SpriteRenderer *components.SpriteRendererComponent // NOUVEAU
	Animation     *components.AnimationComponent
	Collider      *components.ColliderComponent
	Player        *components.PlayerComponent
	Input         *components.InputComponent
	
	// État interne
	EntityID      uint32
	Active        bool
	
	// Sprites
	PlayerSprites interface{} // *PlayerSpriteSet
}

// NewPlayerEntity crée une nouvelle entité joueur avec sprites
func NewPlayerEntity(x, y float64) *PlayerEntity {
	entity := &PlayerEntity{
		Position:       components.NewPositionComponent(x, y),
		Movement:       components.NewMovementComponent(200.0, 250.0),
		Sprite:         components.NewSpriteComponent("player", 32, 32),
		SpriteRenderer: components.NewSpriteRendererComponent(), // NOUVEAU
		Animation:      components.NewAnimationComponent(),
		Collider:       components.NewColliderComponent(24, 24, components.LayerPlayer),
		Player:         components.NewPlayerComponent(),
		Input:          components.NewInputComponent(),
		EntityID:       1,
		Active:         true,
	}
	
	// Configuration du sprite renderer
	entity.SpriteRenderer.Position = components.Vector2{X: x, Y: y}
	entity.SpriteRenderer.Scale = components.Vector2{X: 2.0, Y: 2.0} // Agrandir 2x
	entity.SpriteRenderer.Layer = 10
	
	// Configuration du sprite de fallback
	entity.Sprite.Layer = 10
	entity.Sprite.Color = components.Color{100, 150, 255, 255}
	entity.Collider.Offset = components.Vector2{X: 0, Y: 4}
	
	// Initialiser les animations de fallback
	entity.setupAnimations()
	
	return entity
}

// setupAnimations configure les animations du joueur (fallback)
func (pe *PlayerEntity) setupAnimations() {
	// Animation idle (statique)
	idleFrames := []components.AnimationFrame{
		{SourceRect: components.Rectangle{X: 0, Y: 0, Width: 32, Height: 32}, Duration: time.Second},
	}
	idleAnim := &components.Animation{
		Name:     "idle",
		Frames:   idleFrames,
		Loop:     true,
		PlayRate: 1.0,
	}
	pe.Animation.AddAnimation("idle", idleAnim)
	
	// Animation de marche
	walkFrames := []components.AnimationFrame{
		{SourceRect: components.Rectangle{X: 0, Y: 0, Width: 32, Height: 32}, Duration: time.Millisecond * 200},
		{SourceRect: components.Rectangle{X: 32, Y: 0, Width: 32, Height: 32}, Duration: time.Millisecond * 200},
		{SourceRect: components.Rectangle{X: 64, Y: 0, Width: 32, Height: 32}, Duration: time.Millisecond * 200},
		{SourceRect: components.Rectangle{X: 32, Y: 0, Width: 32, Height: 32}, Duration: time.Millisecond * 200},
	}
	walkAnim := &components.Animation{
		Name:     "walk",
		Frames:   walkFrames,
		Loop:     true,
		PlayRate: 1.5,
	}
	pe.Animation.AddAnimation("walk", walkAnim)
	
	pe.Animation.Play("idle")
}

// SetPlayerSprites définit les sprites du joueur
func (pe *PlayerEntity) SetPlayerSprites(sprites interface{}) {
	pe.PlayerSprites = sprites
	fmt.Println("Sprites du joueur configurés sur l'entité")
}

// GetPosition implémente l'interface Positionable pour PlayerEntity
func (pe *PlayerEntity) GetPosition() components.Vector2 {
	return pe.Position.Position
}

// GetVelocity implémente l'interface Moveable pour PlayerEntity
func (pe *PlayerEntity) GetVelocity() components.Vector2 {
	return pe.Movement.Velocity
}

// ===============================
// SYSTÈME JOUEUR AVEC SPRITES
// ===============================

// PlayerSystem gère la logique du joueur avec support des sprites
type PlayerSystem struct {
	player       *PlayerEntity
	inputManager InputManager
	camera       Camera
	spriteLoader SpriteLoader
	spritesLoaded bool
}

// NewPlayerSystem crée un nouveau système joueur
func NewPlayerSystem() *PlayerSystem {
	return &PlayerSystem{
		player:        nil,
		spritesLoaded: false,
	}
}

// SetInputManager injecte le gestionnaire d'entrées
func (ps *PlayerSystem) SetInputManager(inputManager interface{}) {
	if im, ok := inputManager.(InputManager); ok {
		ps.inputManager = im
	}
}

// SetCamera injecte la caméra
func (ps *PlayerSystem) SetCamera(camera interface{}) {
	if cam, ok := camera.(Camera); ok {
		ps.camera = cam
	}
}

// SetSpriteLoader injecte le chargeur de sprites
func (ps *PlayerSystem) SetSpriteLoader(loader interface{}) {
	if sl, ok := loader.(SpriteLoader); ok {
		ps.spriteLoader = sl
	}
}

// CreatePlayer crée l'entité joueur avec sprites
func (ps *PlayerSystem) CreatePlayer(x, y float64) {
	ps.player = NewPlayerEntity(x, y)
	
	// Charger les sprites si le loader est disponible
	if ps.spriteLoader != nil && !ps.spritesLoaded {
		ps.loadPlayerSprites()
	}
	
	fmt.Printf("Joueur créé à la position (%.1f, %.1f)\n", x, y)
}

// loadPlayerSprites charge les sprites du joueur
func (ps *PlayerSystem) loadPlayerSprites() {
	if ps.spriteLoader == nil {
		fmt.Println("Pas de SpriteLoader disponible, utilisation des sprites de fallback")
		return
	}
	
	sprites, err := ps.spriteLoader.LoadPlayerSprites("assets")
	if err != nil {
		fmt.Printf("Erreur chargement sprites: %v\n", err)
		return
	}
	
	if ps.player != nil {
		ps.player.SetPlayerSprites(sprites)
	}
	
	ps.spritesLoaded = true
}

// GetPlayer retourne l'entité joueur
func (ps *PlayerSystem) GetPlayer() *PlayerEntity {
	return ps.player
}

// GetPlayerPosition retourne la position du joueur
func (ps *PlayerSystem) GetPlayerPosition() components.Vector2 {
	if ps.player != nil {
		return ps.player.Position.Position
	}
	return components.Vector2{X: 0, Y: 0}
}

// Update met à jour le système joueur avec sprites
func (ps *PlayerSystem) Update(deltaTime time.Duration) {
	if ps.player == nil || !ps.player.Active {
		return
	}
	
	// Charger les sprites si pas encore fait
	if !ps.spritesLoaded && ps.spriteLoader != nil {
		ps.loadPlayerSprites()
	}
	
	// Mise à jour dans l'ordre logique
	ps.updateInput(deltaTime)
	ps.updateMovement(deltaTime)
	ps.updateSprites(deltaTime) // NOUVEAU
	ps.updateAnimation(deltaTime)
	ps.updatePlayer(deltaTime)
	ps.updateCamera()
}

// updateInput met à jour les entrées du joueur
func (ps *PlayerSystem) updateInput(deltaTime time.Duration) {
	if ps.inputManager == nil {
		return
	}
	
	input := ps.player.Input
	
	// Reset des actions de la frame précédente
	input.Reset()
	
	// Actions de mouvement
	input.MoveUp = ps.inputManager.IsActionPressedSystems(0)    // ActionMoveUp
	input.MoveDown = ps.inputManager.IsActionPressedSystems(1)  // ActionMoveDown
	input.MoveLeft = ps.inputManager.IsActionPressedSystems(2)  // ActionMoveLeft
	input.MoveRight = ps.inputManager.IsActionPressedSystems(3) // ActionMoveRight
	
	// Actions "just pressed"
	input.AttackJustPressed = ps.inputManager.IsKeyJustPressedSystems(32)  // Espace
	input.RollJustPressed = ps.inputManager.IsKeyJustPressedSystems(99)    // C
	input.InteractJustPressed = ps.inputManager.IsKeyJustPressedSystems(101) // E
	
	// Actions maintenues
	input.Block = ps.inputManager.IsActionPressedSystems(5) // ActionBlock
}

// updateSprites met à jour le système de sprites (NOUVEAU)
func (ps *PlayerSystem) updateSprites(deltaTime time.Duration) {
	if ps.player == nil || ps.player.SpriteRenderer == nil {
		return
	}
	
	spriteRenderer := ps.player.SpriteRenderer
	movement := ps.player.Movement
	
	// Mettre à jour la position du sprite
	spriteRenderer.Position = ps.player.Position.Position
	
	// Déterminer la direction
	var direction string
	switch movement.FacingDir {
	case components.DirectionUp, components.DirectionUpLeft, components.DirectionUpRight:
		direction = "up"
	case components.DirectionDown, components.DirectionDownLeft, components.DirectionDownRight:
		direction = "down"
	case components.DirectionLeft:
		direction = "left"
	case components.DirectionRight:
		direction = "right"
	default:
		direction = "down"
	}
	
	// Mettre à jour la direction et l'état
	spriteRenderer.SetDirection(direction, movement.IsMoving)
	
	// Mettre à jour l'animation du sprite
	spriteRenderer.Update(deltaTime)
	
	// Appliquer le flip horizontal pour gauche/droite si nécessaire
	if direction == "left" {
		spriteRenderer.FlipX = true
	} else if direction == "right" {
		spriteRenderer.FlipX = false
	}
}

// updateMovement met à jour le mouvement du joueur
func (ps *PlayerSystem) updateMovement(deltaTime time.Duration) {
	movement := ps.player.Movement
	position := ps.player.Position
	input := ps.player.Input
	
	dt := deltaTime.Seconds()
	
	// Calculer le vecteur de mouvement depuis les inputs
	inputVector := input.GetMovementVector()
	
	// Appliquer l'accélération ou la friction
	if inputVector.X != 0 || inputVector.Y != 0 {
		movement.IsMoving = true
		
		targetVelocity := inputVector.Mul(movement.Speed)
		velocityDiff := targetVelocity.Sub(movement.Velocity)
		acceleration := velocityDiff.Mul(movement.Acceleration * dt)
		
		movement.Velocity = movement.Velocity.Add(acceleration)
		
		// Limiter à la vitesse maximale
		velocityLengthSq := movement.Velocity.X*movement.Velocity.X + movement.Velocity.Y*movement.Velocity.Y
		if velocityLengthSq > movement.MaxSpeed*movement.MaxSpeed {
			invLength := movement.MaxSpeed / math.Sqrt(velocityLengthSq)
			movement.Velocity.X *= invLength
			movement.Velocity.Y *= invLength
		}
		
		movement.Direction = ps.vectorToDirection(inputVector)
		movement.FacingDir = movement.Direction
		
	} else {
		movement.IsMoving = false
		
		velocityLength := math.Sqrt(movement.Velocity.X*movement.Velocity.X + movement.Velocity.Y*movement.Velocity.Y)
		if velocityLength > 0 {
			frictionMagnitude := movement.Friction * dt
			if frictionMagnitude >= velocityLength {
				movement.Velocity = components.Vector2{X: 0, Y: 0}
			} else {
				invLength := frictionMagnitude / velocityLength
				movement.Velocity.X -= movement.Velocity.X * invLength
				movement.Velocity.Y -= movement.Velocity.Y * invLength
			}
		}
	}
	
	position.LastPosition = position.Position
	position.Position = position.Position.Add(movement.Velocity.Mul(dt))
	
	ps.applyScreenBounds()
}

// vectorToDirection convertit un vecteur en direction
func (ps *PlayerSystem) vectorToDirection(vector components.Vector2) components.Direction {
	if vector.X == 0 && vector.Y == 0 {
		return components.DirectionNone
	}
	
	if vector.X == 0 {
		if vector.Y < 0 {
			return components.DirectionUp
		}
		return components.DirectionDown
	}
	
	if vector.Y == 0 {
		if vector.X < 0 {
			return components.DirectionLeft
		}
		return components.DirectionRight
	}
	
	if vector.X < 0 && vector.Y < 0 {
		return components.DirectionUpLeft
	}
	if vector.X > 0 && vector.Y < 0 {
		return components.DirectionUpRight
	}
	if vector.X < 0 && vector.Y > 0 {
		return components.DirectionDownLeft
	}
	return components.DirectionDownRight
}

// applyScreenBounds limite le joueur aux bords de l'écran
func (ps *PlayerSystem) applyScreenBounds() {
	position := ps.player.Position
	size := ps.player.Sprite.Size
	
	const margin = 16
	minX := margin + size.X/2
	maxX := 1280 - margin - size.X/2
	minY := margin + size.Y/2
	maxY := 720 - margin - size.Y/2
	
	if position.Position.X < minX {
		position.Position.X = minX
		ps.player.Movement.Velocity.X = 0
	} else if position.Position.X > maxX {
		position.Position.X = maxX
		ps.player.Movement.Velocity.X = 0
	}
	
	if position.Position.Y < minY {
		position.Position.Y = minY
		ps.player.Movement.Velocity.Y = 0
	} else if position.Position.Y > maxY {
		position.Position.Y = maxY
		ps.player.Movement.Velocity.Y = 0
	}
}

// updateAnimation met à jour les animations
func (ps *PlayerSystem) updateAnimation(deltaTime time.Duration) {
	animation := ps.player.Animation
	movement := ps.player.Movement
	sprite := ps.player.Sprite
	
	// Choisir l'animation appropriée
	var targetAnim string
	if movement.IsMoving {
		targetAnim = "walk"
	} else {
		targetAnim = "idle"
	}
	
	// Changer d'animation si nécessaire
	if !animation.IsPlaying(targetAnim) {
		animation.Play(targetAnim)
	}
	
	// Mettre à jour l'animation de fallback
	if animation.Playing {
		animation.ElapsedTime += deltaTime
		
		currentAnimation := animation.Animations[animation.CurrentAnim]
		if currentAnimation != nil && len(currentAnimation.Frames) > 0 {
			currentFrame := &currentAnimation.Frames[animation.CurrentFrame]
			frameDuration := time.Duration(float64(currentFrame.Duration) / animation.PlayRate)
			
			if animation.ElapsedTime >= frameDuration {
				animation.ElapsedTime = 0
				animation.CurrentFrame++
				
				if animation.CurrentFrame >= len(currentAnimation.Frames) {
					if currentAnimation.Loop {
						animation.CurrentFrame = 0
					} else {
						animation.Playing = false
						animation.CurrentFrame = len(currentAnimation.Frames) - 1
						if animation.OnComplete != nil {
							animation.OnComplete()
						}
					}
				}
			}
			
			frame := animation.GetCurrentFrame()
			if frame != nil {
				sprite.SourceRect = frame.SourceRect
			}
		}
	}
	
	// Appliquer le flip horizontal selon la direction
	if movement.FacingDir == components.DirectionLeft || 
	   movement.FacingDir == components.DirectionUpLeft || 
	   movement.FacingDir == components.DirectionDownLeft {
		sprite.FlipX = true
	} else if movement.FacingDir == components.DirectionRight || 
	          movement.FacingDir == components.DirectionUpRight || 
	          movement.FacingDir == components.DirectionDownRight {
		sprite.FlipX = false
	}
}

// updatePlayer met à jour les stats du joueur et traite les actions
func (ps *PlayerSystem) updatePlayer(deltaTime time.Duration) {
	ps.player.Player.Update(deltaTime)
	ps.handlePlayerActions()
}

// handlePlayerActions traite les actions spéciales du joueur
func (ps *PlayerSystem) handlePlayerActions() {
	if !ps.player.Player.IsAlive() {
		return
	}
	
	input := ps.player.Input
	
	if input.AttackJustPressed {
		if ps.TryAttack() && ps.player.SpriteRenderer != nil {
			ps.player.SpriteRenderer.StartAttack() // NOUVEAU: déclencher animation d'attaque
		}
	}
	
	if input.RollJustPressed {
		ps.TryRoll()
	}
	
	if input.InteractJustPressed {
		ps.TryInteract()
	}
}

// updateCamera met à jour la caméra pour suivre le joueur
func (ps *PlayerSystem) updateCamera() {
	if ps.camera == nil {
		return
	}
	
	offset := components.Vector2{X: 0, Y: -20}
	ps.camera.FollowTarget(ps.player, 3.0, offset)
}

// Render rend le joueur avec sprites ou fallback
func (ps *PlayerSystem) Render(renderer Renderer) {
	if ps.player == nil || !ps.player.Active {
		return
	}
	
	// Essayer d'abord le rendu avec sprites
	if ps.renderWithSprites(renderer) {
		return
	}
	
	// Fallback vers le rendu rectangulaire
	ps.renderFallback(renderer)
}

// renderWithSprites tente le rendu avec les vrais sprites
func (ps *PlayerSystem) renderWithSprites(renderer Renderer) bool {
	if ps.player.PlayerSprites == nil || ps.player.SpriteRenderer == nil {
		return false
	}
	
	spriteRenderer := ps.player.SpriteRenderer
	if !spriteRenderer.Visible {
		return false
	}
	
	// Vérifier si le renderer supporte DrawSprite
	if spriteRenderer, ok := renderer.(interface {
		DrawSprite(sprite interface{}, position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, rotation float64, tint components.Color)
	}); ok {
		
		position := ps.player.SpriteRenderer.Position
		sourceRect := ps.player.SpriteRenderer.SourceRect
		scale := ps.player.SpriteRenderer.Scale
		rotation := ps.player.SpriteRenderer.Rotation
		tint := ps.player.SpriteRenderer.Tint
		
		// Utiliser le sprite principal pour l'instant
		spriteRenderer.DrawSprite(nil, position, sourceRect, scale, rotation, tint)
		
		return true
	}
	
	return false
}

// renderFallback rendu rectangulaire de fallback
func (ps *PlayerSystem) renderFallback(renderer Renderer) {
	if !ps.player.Sprite.Visible {
		return
	}
	
	position := ps.player.Position.Position
	sprite := ps.player.Sprite
	
	playerRect := components.Rectangle{
		X:      position.X - sprite.Size.X/2 + sprite.Offset.X,
		Y:      position.Y - sprite.Size.Y/2 + sprite.Offset.Y,
		Width:  sprite.Size.X,
		Height: sprite.Size.Y,
	}
	
	color := sprite.Color
	if ps.player.Movement.IsMoving {
		color = components.Color{
			R: minByte(sprite.Color.R + 30, 255),
			G: minByte(sprite.Color.G + 30, 255),
			B: minByte(sprite.Color.B + 30, 255),
			A: sprite.Color.A,
		}
	}
	
	if ps.player.Player.InvulnTime > 0 {
		if (ps.player.Player.InvulnTime.Milliseconds()/100)%2 == 0 {
			color.A = 128
		}
	}
	
	renderer.DrawRectangle(playerRect, color, true)
	
	borderColor := components.ColorWhite
	if ps.player.Player.InvulnTime > 0 {
		borderColor = components.ColorYellow
	}
	renderer.DrawRectangle(playerRect, borderColor, false)
	
	if ps.player.Movement.IsMoving {
		ps.renderDirectionIndicator(renderer, position)
	}
	
	ps.renderHealthBar(renderer, position)
	ps.renderStaminaBar(renderer, position)
}

// renderDirectionIndicator dessine un indicateur de direction
func (ps *PlayerSystem) renderDirectionIndicator(renderer Renderer, position components.Vector2) {
	direction := ps.player.Movement.Direction
	if direction == components.DirectionNone {
		return
	}
	
	arrowLength := 15.0
	dirVector := direction.ToVector2()
	arrowEnd := position.Add(dirVector.Mul(arrowLength))
	
	if dirVector.X != 0 {
		arrowRect := components.Rectangle{
			X:      math.Min(position.X, arrowEnd.X) - 1,
			Y:      position.Y - 1,
			Width:  math.Abs(arrowEnd.X - position.X) + 2,
			Height: 2,
		}
		renderer.DrawRectangle(arrowRect, components.ColorYellow, true)
	}
	if dirVector.Y != 0 {
		arrowRect := components.Rectangle{
			X:      position.X - 1,
			Y:      math.Min(position.Y, arrowEnd.Y) - 1,
			Width:  2,
			Height: math.Abs(arrowEnd.Y - position.Y) + 2,
		}
		renderer.DrawRectangle(arrowRect, components.ColorYellow, true)
	}
}

// renderHealthBar dessine une barre de vie
func (ps *PlayerSystem) renderHealthBar(renderer Renderer, position components.Vector2) {
	player := ps.player.Player
	
	barWidth := 30.0
	barHeight := 4.0
	barY := position.Y - ps.player.Sprite.Size.Y/2 - 8
	barX := position.X - barWidth/2
	
	bgRect := components.Rectangle{X: barX, Y: barY, Width: barWidth, Height: barHeight}
	renderer.DrawRectangle(bgRect, components.ColorBlack, true)
	
	healthPercent := float64(player.Health) / float64(player.MaxHealth)
	healthWidth := barWidth * healthPercent
	
	var healthColor components.Color
	if healthPercent > 0.6 {
		healthColor = components.ColorGreen
	} else if healthPercent > 0.3 {
		healthColor = components.ColorYellow
	} else {
		healthColor = components.ColorRed
	}
	
	if healthWidth > 0 {
		healthRect := components.Rectangle{X: barX, Y: barY, Width: healthWidth, Height: barHeight}
		renderer.DrawRectangle(healthRect, healthColor, true)
	}
	
	renderer.DrawRectangle(bgRect, components.ColorWhite, false)
}

// renderStaminaBar dessine une barre de stamina
func (ps *PlayerSystem) renderStaminaBar(renderer Renderer, position components.Vector2) {
	player := ps.player.Player
	
	barWidth := 30.0
	barHeight := 3.0
	barY := position.Y - ps.player.Sprite.Size.Y/2 - 14
	barX := position.X - barWidth/2
	
	bgRect := components.Rectangle{X: barX, Y: barY, Width: barWidth, Height: barHeight}
	renderer.DrawRectangle(bgRect, components.ColorBlack, true)
	
	staminaPercent := player.Stamina / player.MaxStamina
	staminaWidth := barWidth * staminaPercent
	
	if staminaWidth > 0 {
		staminaRect := components.Rectangle{X: barX, Y: barY, Width: staminaWidth, Height: barHeight}
		renderer.DrawRectangle(staminaRect, components.ColorCyan, true)
	}
	
	renderer.DrawRectangle(bgRect, components.ColorGray, false)
}

// ===============================
// ACTIONS DU JOUEUR
// ===============================

// TryAttack tente une attaque
func (ps *PlayerSystem) TryAttack() bool {
	if ps.player == nil || !ps.player.Player.IsAlive() {
		return false
	}
	
	staminaCost := 15.0
	if !ps.player.Player.UseStamina(staminaCost) {
		fmt.Println("Pas assez de stamina pour attaquer!")
		return false
	}
	
	fmt.Println("Attaque réussie!")
	return true
}

// TryRoll tente une roulade
func (ps *PlayerSystem) TryRoll() bool {
	if ps.player == nil || !ps.player.Player.IsAlive() {
		return false
	}
	
	staminaCost := 25.0
	if !ps.player.Player.UseStamina(staminaCost) {
		fmt.Println("Pas assez de stamina pour rouler!")
		return false
	}
	
	rollDirection := ps.player.Movement.Direction
	if rollDirection == components.DirectionNone {
		rollDirection = ps.player.Movement.FacingDir
	}
	
	rollSpeed := 400.0
	rollVector := rollDirection.ToVector2().Mul(rollSpeed)
	ps.player.Movement.Velocity = rollVector
	
	ps.player.Player.InvulnTime = time.Millisecond * 300
	
	fmt.Println("Roulade effectuée!")
	return true
}

// TryInteract tente une interaction
func (ps *PlayerSystem) TryInteract() bool {
	if ps.player == nil || !ps.player.Player.IsAlive() {
		return false
	}
	
	fmt.Println("Interaction (rien à proximité)")
	return false
}

// ===============================
// MÉTHODES UTILITAIRES
// ===============================

func (ps *PlayerSystem) IsPlayerAlive() bool {
	return ps.player != nil && ps.player.Player.IsAlive()
}

func (ps *PlayerSystem) GetPlayerHealth() (int, int) {
	if ps.player == nil {
		return 0, 0
	}
	return ps.player.Player.Health, ps.player.Player.MaxHealth
}

func (ps *PlayerSystem) GetPlayerStamina() (float64, float64) {
	if ps.player == nil {
		return 0, 0
	}
	return ps.player.Player.Stamina, ps.player.Player.MaxStamina
}

func minByte(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}