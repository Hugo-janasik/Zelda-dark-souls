
package systems

import (
	"fmt"
	"math"
	"time"
	"zelda-souls-game/internal/ecs/components"
)

// ===============================
// INTERFACES MINIMALES
// ===============================

// InputManager interface minimale pour éviter les cycles
type InputManager interface {
	IsActionPressed(action int) bool
	IsKeyJustPressed(key int) bool
}

// Renderer interface minimale pour le rendu
type Renderer interface {
	DrawRectangle(rect components.Rectangle, color components.Color, filled bool)
	DrawText(text string, pos components.Vector2, color components.Color)
}

// Camera interface pour la caméra
type Camera interface {
	SetTarget(position components.Vector2)
	FollowTarget(target interface{}, speed float64, offset components.Vector2)
}

// ===============================
// ENTITÉ JOUEUR
// ===============================

// PlayerEntity représente l'entité joueur complète
type PlayerEntity struct {
	// Composants
	Position   *components.PositionComponent
	Movement   *components.MovementComponent
	Sprite     *components.SpriteComponent
	Animation  *components.AnimationComponent
	Collider   *components.ColliderComponent
	Player     *components.PlayerComponent
	Input      *components.InputComponent
	
	// État interne
	EntityID   uint32
	Active     bool
}

// NewPlayerEntity crée une nouvelle entité joueur
func NewPlayerEntity(x, y float64) *PlayerEntity {
	entity := &PlayerEntity{
		Position:  components.NewPositionComponent(x, y),
		Movement:  components.NewMovementComponent(200.0, 250.0), // Vitesse base et max
		Sprite:    components.NewSpriteComponent("player", 32, 32),
		Animation: components.NewAnimationComponent(),
		Collider:  components.NewColliderComponent(24, 24, components.LayerPlayer),
		Player:    components.NewPlayerComponent(),
		Input:     components.NewInputComponent(),
		EntityID:  1, // ID fixe pour le joueur
		Active:    true,
	}
	
	// Configuration du sprite
	entity.Sprite.Layer = 10 // Couche élevée pour être visible
	entity.Sprite.Color = components.Color{100, 150, 255, 255} // Bleu pour le test
	
	// Configuration du collider
	entity.Collider.Offset = components.Vector2{X: 0, Y: 4} // Légèrement vers le bas
	
	// Initialiser les animations
	entity.setupAnimations()
	
	return entity
}

// setupAnimations configure les animations du joueur
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
	
	// Animation de marche (pour l'instant on simule avec 4 frames)
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
		PlayRate: 1.5, // Un peu plus rapide
	}
	pe.Animation.AddAnimation("walk", walkAnim)
	
	// Démarrer avec l'animation idle
	pe.Animation.Play("idle")
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
// SYSTÈME JOUEUR
// ===============================

// PlayerSystem gère la logique du joueur
type PlayerSystem struct {
	player      *PlayerEntity
	inputManager InputManager
	camera      Camera
}

// NewPlayerSystem crée un nouveau système joueur
func NewPlayerSystem() *PlayerSystem {
	return &PlayerSystem{
		player: nil,
	}
}

// SetInputManager injecte le gestionnaire d'entrées
func (ps *PlayerSystem) SetInputManager(inputManager InputManager) {
	ps.inputManager = inputManager
}

// SetCamera injecte la caméra
func (ps *PlayerSystem) SetCamera(camera interface{}) {
	if cam, ok := camera.(Camera); ok {
		ps.camera = cam
	}
}

// CreatePlayer crée l'entité joueur
func (ps *PlayerSystem) CreatePlayer(x, y float64) {
	ps.player = NewPlayerEntity(x, y)
	fmt.Printf("Joueur créé à la position (%.1f, %.1f)\n", x, y)
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

// Update met à jour le système joueur
func (ps *PlayerSystem) Update(deltaTime time.Duration) {
	if ps.player == nil || !ps.player.Active {
		return
	}
	
	// Mise à jour dans l'ordre logique
	ps.updateInput(deltaTime)
	ps.updateMovement(deltaTime)
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
	input.MoveUp = ps.inputManager.IsActionPressed(0)    // ActionMoveUp
	input.MoveDown = ps.inputManager.IsActionPressed(1)  // ActionMoveDown
	input.MoveLeft = ps.inputManager.IsActionPressed(2)  // ActionMoveLeft
	input.MoveRight = ps.inputManager.IsActionPressed(3) // ActionMoveRight
	
	// Actions "just pressed"
	input.AttackJustPressed = ps.inputManager.IsKeyJustPressed(32)  // Espace
	input.RollJustPressed = ps.inputManager.IsKeyJustPressed(99)    // C
	input.InteractJustPressed = ps.inputManager.IsKeyJustPressed(101) // E
	
	// Actions maintenues
	input.Block = ps.inputManager.IsActionPressed(5) // ActionBlock
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
		// Le joueur veut bouger
		movement.IsMoving = true
		
		// Calculer la vélocité cible
		targetVelocity := inputVector.Mul(movement.Speed)
		
		// Appliquer l'accélération vers la vélocité cible
		velocityDiff := targetVelocity.Sub(movement.Velocity)
		acceleration := velocityDiff.Mul(movement.Acceleration * dt)
		
		movement.Velocity = movement.Velocity.Add(acceleration)
		
		// Limiter à la vitesse maximale
		if movement.Velocity.Length() > movement.MaxSpeed*movement.MaxSpeed {
			movement.Velocity = movement.Velocity.Normalize().Mul(movement.MaxSpeed)
		}
		
		// Déterminer la direction de mouvement
		movement.Direction = ps.vectorToDirection(inputVector)
		movement.FacingDir = movement.Direction
		
	} else {
		// Le joueur ne bouge pas, appliquer la friction
		movement.IsMoving = false
		
		frictionForce := movement.Velocity.Normalize().Mul(-movement.Friction * dt)
		movement.Velocity = movement.Velocity.Add(frictionForce)
		
		// Arrêter si la vélocité devient très faible
		if movement.Velocity.Length() < 25.0 { // 5^2
			movement.Velocity = components.Vector2{X: 0, Y: 0}
		}
	}
	
	// Sauvegarder la position précédente
	position.LastPosition = position.Position
	
	// Appliquer le mouvement
	position.Position = position.Position.Add(movement.Velocity.Mul(dt))
	
	// Limiter aux bords de l'écran (temporaire)
	ps.applyScreenBounds()
}

// vectorToDirection convertit un vecteur en direction
func (ps *PlayerSystem) vectorToDirection(vector components.Vector2) components.Direction {
	if vector.X == 0 && vector.Y == 0 {
		return components.DirectionNone
	}
	
	// Directions cardinales pures
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
	
	// Directions diagonales
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

// applyScreenBounds limite le joueur aux bords de l'écran (temporaire)
func (ps *PlayerSystem) applyScreenBounds() {
	position := ps.player.Position
	size := ps.player.Sprite.Size
	
	// Limites temporaires (sera remplacé par les collisions de monde)
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
	
	// Mettre à jour l'animation
	if animation.Playing {
		animation.ElapsedTime += deltaTime
		
		currentAnimation := animation.Animations[animation.CurrentAnim]
		if currentAnimation != nil && len(currentAnimation.Frames) > 0 {
			// Calculer la frame actuelle
			currentFrame := &currentAnimation.Frames[animation.CurrentFrame]
			frameDuration := time.Duration(float64(currentFrame.Duration) / animation.PlayRate)
			
			if animation.ElapsedTime >= frameDuration {
				animation.ElapsedTime = 0
				animation.CurrentFrame++
				
				// Vérifier si on a atteint la fin
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
			
			// Mettre à jour le sprite avec la frame actuelle
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

// updatePlayer met à jour les stats du joueur
func (ps *PlayerSystem) updatePlayer(deltaTime time.Duration) {
	ps.player.Player.Update(deltaTime)
	
	// Traiter les actions spéciales du joueur
	ps.handlePlayerActions()
}

// handlePlayerActions traite les actions spéciales du joueur
func (ps *PlayerSystem) handlePlayerActions() {
	if !ps.player.Player.IsAlive() {
		return
	}
	
	input := ps.player.Input
	
	// Traiter les actions "just pressed"
	if input.AttackJustPressed {
		ps.TryAttack()
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
	
	// Faire suivre la caméra au joueur avec un léger décalage et lissage
	offset := components.Vector2{X: 0, Y: -20} // Légèrement au-dessus du joueur
	ps.camera.FollowTarget(ps.player, 3.0, offset) // Vitesse de suivi modérée
}

// Render rend le joueur
func (ps *PlayerSystem) Render(renderer Renderer) {
	if ps.player == nil || !ps.player.Active || !ps.player.Sprite.Visible {
		return
	}
	
	position := ps.player.Position.Position
	sprite := ps.player.Sprite
	
	// Pour l'instant, on dessine un rectangle coloré
	// Plus tard, on utilisera DrawSprite avec une vraie texture
	playerRect := components.Rectangle{
		X:      position.X - sprite.Size.X/2 + sprite.Offset.X,
		Y:      position.Y - sprite.Size.Y/2 + sprite.Offset.Y,
		Width:  sprite.Size.X,
		Height: sprite.Size.Y,
	}
	
	// Couleur différente selon l'état
	color := sprite.Color
	if ps.player.Movement.IsMoving {
		// Plus brillant quand il bouge
		color = components.Color{
			R: minByte(sprite.Color.R + 30, 255),
			G: minByte(sprite.Color.G + 30, 255),
			B: minByte(sprite.Color.B + 30, 255),
			A: sprite.Color.A,
		}
	}
	
	// Clignotement si invulnérable
	if ps.player.Player.InvulnTime > 0 {
		// Clignoter toutes les 100ms
		if (ps.player.Player.InvulnTime.Milliseconds()/100)%2 == 0 {
			color.A = 128 // Semi-transparent
		}
	}
	
	// Dessiner le joueur
	renderer.DrawRectangle(playerRect, color, true)
	
	// Bordure
	borderColor := components.ColorWhite
	if ps.player.Player.InvulnTime > 0 {
		borderColor = components.ColorYellow // Bordure jaune si invulnérable
	}
	renderer.DrawRectangle(playerRect, borderColor, false)
	
	// Indicateur de direction
	if ps.player.Movement.IsMoving {
		ps.renderDirectionIndicator(renderer, position)
	}
	
	// Barre de vie (debug)
	ps.renderHealthBar(renderer, position)
	
	// Barre de stamina (debug)
	ps.renderStaminaBar(renderer, position)
}

// renderDirectionIndicator dessine un indicateur de direction
func (ps *PlayerSystem) renderDirectionIndicator(renderer Renderer, position components.Vector2) {
	direction := ps.player.Movement.Direction
	if direction == components.DirectionNone {
		return
	}
	
	// Calculer la position de la flèche
	arrowLength := 15.0
	dirVector := direction.ToVector2()
	arrowEnd := position.Add(dirVector.Mul(arrowLength))
	
	// Dessiner une ligne simple pour la direction
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
	
	// Fond de la barre
	bgRect := components.Rectangle{X: barX, Y: barY, Width: barWidth, Height: barHeight}
	renderer.DrawRectangle(bgRect, components.ColorBlack, true)
	
	// Barre de vie
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
	
	// Bordure
	renderer.DrawRectangle(bgRect, components.ColorWhite, false)
}

// renderStaminaBar dessine une barre de stamina
func (ps *PlayerSystem) renderStaminaBar(renderer Renderer, position components.Vector2) {
	player := ps.player.Player
	
	barWidth := 30.0
	barHeight := 3.0
	barY := position.Y - ps.player.Sprite.Size.Y/2 - 14
	barX := position.X - barWidth/2
	
	// Fond de la barre
	bgRect := components.Rectangle{X: barX, Y: barY, Width: barWidth, Height: barHeight}
	renderer.DrawRectangle(bgRect, components.ColorBlack, true)
	
	// Barre de stamina
	staminaPercent := player.Stamina / player.MaxStamina
	staminaWidth := barWidth * staminaPercent
	
	if staminaWidth > 0 {
		staminaRect := components.Rectangle{X: barX, Y: barY, Width: staminaWidth, Height: barHeight}
		renderer.DrawRectangle(staminaRect, components.ColorCyan, true)
	}
	
	// Bordure
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
	
	// Coût en stamina pour attaquer
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
	
	// Coût en stamina pour rouler
	staminaCost := 25.0
	if !ps.player.Player.UseStamina(staminaCost) {
		fmt.Println("Pas assez de stamina pour rouler!")
		return false
	}
	
	// Boost de vitesse temporaire
	rollDirection := ps.player.Movement.Direction
	if rollDirection == components.DirectionNone {
		rollDirection = ps.player.Movement.FacingDir
	}
	
	rollSpeed := 400.0
	rollVector := rollDirection.ToVector2().Mul(rollSpeed)
	ps.player.Movement.Velocity = rollVector
	
	// Invulnérabilité temporaire
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

// Méthodes utilitaires pour le système
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

// Fonctions utilitaires
func minByte(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}