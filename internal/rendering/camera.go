// internal/rendering/camera.go - Système de caméra
package rendering

import (
	"math"
	"time"

	"zelda-souls-game/internal/core"
)

// ===============================
// CAMERA STRUCTURE
// ===============================

// Camera gère la vue et la transformation du monde vers l'écran
type Camera struct {
	// Position et dimensions
	Position core.Vector2 // Position de la caméra dans le monde
	Width    float64      // Largeur de la vue
	Height   float64      // Hauteur de la vue
	Zoom     float64      // Niveau de zoom (1.0 = normal)

	// Limites de la caméra
	Bounds *core.Rectangle // Limites dans lesquelles la caméra peut bouger

	// Suivi d'entité
	Target      interface{}  // Entité à suivre (Player, etc.)
	FollowSpeed float64      // Vitesse de suivi (0-1, 1=instantané)
	Offset      core.Vector2 // Décalage par rapport à la cible

	// Effets de caméra
	Shake *CameraShake

	// Interpolation
	targetPosition core.Vector2
	velocity       core.Vector2
	smoothing      float64

	// Limites de zoom
	MinZoom float64
	MaxZoom float64

	// État interne
	viewMatrix [6]float64 // Matrice de transformation
	needUpdate bool
}

// CameraShake gère les effets de tremblement
type CameraShake struct {
	Intensity float64
	Duration  time.Duration
	Frequency float64

	// État interne
	startTime   time.Time
	currentTime float64
	offset      core.Vector2
	active      bool
}

// ===============================
// CAMERA INITIALIZATION
// ===============================

// NewCamera crée une nouvelle caméra
func NewCamera(position core.Vector2, width, height float64) *Camera {
	camera := &Camera{
		Position:    position,
		Width:       width,
		Height:      height,
		Zoom:        1.0,
		FollowSpeed: 5.0,
		MinZoom:     0.1,
		MaxZoom:     5.0,
		smoothing:   0.1,
		needUpdate:  true,
	}

	camera.targetPosition = position
	camera.updateViewMatrix()

	return camera
}

// ===============================
// CAMERA CONTROL
// ===============================

// SetPosition définit immédiatement la position de la caméra
func (c *Camera) SetPosition(position core.Vector2) {
	c.Position = position
	c.targetPosition = position
	c.needUpdate = true
}

// SetTarget définit la position cible (avec interpolation)
func (c *Camera) SetTarget(position core.Vector2) {
	c.targetPosition = position
}

// SetZoom définit le niveau de zoom
func (c *Camera) SetZoom(zoom float64) {
	// Limiter le zoom
	if zoom < c.MinZoom {
		zoom = c.MinZoom
	}
	if zoom > c.MaxZoom {
		zoom = c.MaxZoom
	}

	c.Zoom = zoom
	c.needUpdate = true
}

// SetBounds définit les limites de mouvement de la caméra
func (c *Camera) SetBounds(bounds core.Rectangle) {
	c.Bounds = &bounds
	c.applyBounds()
}

// RemoveBounds supprime les limites de la caméra
func (c *Camera) RemoveBounds() {
	c.Bounds = nil
}

// FollowTarget fait suivre une entité à la caméra
func (c *Camera) FollowTarget(target interface{}, speed float64, offset core.Vector2) {
	c.Target = target
	c.FollowSpeed = speed
	c.Offset = offset
}

// StopFollowing arrête le suivi de cible
func (c *Camera) StopFollowing() {
	c.Target = nil
}

// ===============================
// CAMERA UPDATE
// ===============================

// Update met à jour la caméra
func (c *Camera) Update(deltaTime time.Duration) {
	dt := deltaTime.Seconds()

	// Mise à jour du suivi de cible
	c.updateTargetFollowing(dt)

	// Interpolation vers la position cible
	c.updateMovementSmoothing(dt)

	// Mise à jour des effets de tremblement
	c.updateShake(dt)

	// Appliquer les limites
	c.applyBounds()

	// Mettre à jour la matrice de vue si nécessaire
	if c.needUpdate {
		c.updateViewMatrix()
		c.needUpdate = false
	}
}

// updateTargetFollowing met à jour le suivi de la cible
func (c *Camera) updateTargetFollowing(deltaTime float64) {
	if c.Target == nil {
		return
	}

	// Obtenir la position de la cible
	var targetPos core.Vector2

	// Interface pour les objets avec position
	type Positionable interface {
		GetPosition() core.Vector2
	}

	if positionable, ok := c.Target.(Positionable); ok {
		targetPos = positionable.GetPosition()
	} else {
		return // Cible non compatible
	}

	// Ajouter le décalage
	targetPos = targetPos.Add(c.Offset)

	// Interpolation vers la cible
	if c.FollowSpeed >= 1.0 {
		// Suivi instantané
		c.targetPosition = targetPos
	} else {
		// Suivi avec lissage
		diff := targetPos.Sub(c.targetPosition)
		c.targetPosition = c.targetPosition.Add(diff.Mul(c.FollowSpeed * deltaTime))
	}
}

// updateMovementSmoothing applique le lissage de mouvement
func (c *Camera) updateMovementSmoothing(deltaTime float64) {
	if c.smoothing <= 0 {
		c.Position = c.targetPosition
		return
	}

	// Calculer la force vers la cible
	diff := c.targetPosition.Sub(c.Position)
	force := diff.Mul(1.0 / c.smoothing)

	// Appliquer la force à la vélocité
	c.velocity = c.velocity.Add(force.Mul(deltaTime))

	// Appliquer un amortissement
	damping := 1.0 - (5.0 * deltaTime)
	if damping < 0 {
		damping = 0
	}
	c.velocity = c.velocity.Mul(damping)

	// Mettre à jour la position
	oldPos := c.Position
	c.Position = c.Position.Add(c.velocity.Mul(deltaTime))

	// Marquer pour mise à jour si la position a changé
	if oldPos.X != c.Position.X || oldPos.Y != c.Position.Y {
		c.needUpdate = true
	}
}

// updateShake met à jour les effets de tremblement
func (c *Camera) updateShake(deltaTime float64) {
	if c.Shake == nil || !c.Shake.active {
		return
	}

	elapsed := time.Since(c.Shake.startTime)

	// Vérifier si le shake est terminé
	if elapsed >= c.Shake.Duration {
		c.Shake.active = false
		c.Shake.offset = core.Vector2{X: 0, Y: 0}
		c.needUpdate = true
		return
	}

	// Calculer l'intensité décroissante
	progress := float64(elapsed) / float64(c.Shake.Duration)
	intensity := c.Shake.Intensity * (1.0 - progress)

	// Générer un offset aléatoire basé sur le temps et la fréquence
	c.Shake.currentTime += deltaTime
	time := c.Shake.currentTime * c.Shake.Frequency

	// Utilisation de fonctions sinusoïdales pour un shake plus naturel
	c.Shake.offset.X = intensity * math.Sin(time*2.3) * math.Cos(time*1.7)
	c.Shake.offset.Y = intensity * math.Cos(time*2.1) * math.Sin(time*1.9)

	c.needUpdate = true
}

// applyBounds applique les limites de mouvement
func (c *Camera) applyBounds() {
	if c.Bounds == nil {
		return
	}

	// Calculer les limites effectives en tenant compte du zoom
	halfWidth := (c.Width / c.Zoom) / 2
	halfHeight := (c.Height / c.Zoom) / 2

	// Limites min/max
	minX := c.Bounds.X + halfWidth
	maxX := c.Bounds.X + c.Bounds.Width - halfWidth
	minY := c.Bounds.Y + halfHeight
	maxY := c.Bounds.Y + c.Bounds.Height - halfHeight

	// Appliquer les contraintes
	oldPos := c.Position

	if c.Position.X < minX {
		c.Position.X = minX
	}
	if c.Position.X > maxX {
		c.Position.X = maxX
	}
	if c.Position.Y < minY {
		c.Position.Y = minY
	}
	if c.Position.Y > maxY {
		c.Position.Y = maxY
	}

	// Mettre à jour la cible aussi pour éviter le conflit
	if oldPos.X != c.Position.X || oldPos.Y != c.Position.Y {
		c.targetPosition = c.Position
		c.needUpdate = true
	}
}

// updateViewMatrix met à jour la matrice de transformation
func (c *Camera) updateViewMatrix() {
	// Position finale avec shake
	finalPos := c.Position
	if c.Shake != nil && c.Shake.active {
		finalPos = finalPos.Add(c.Shake.offset)
	}

	// Matrice de transformation 2D
	// Translation pour centrer la caméra
	tx := -finalPos.X + c.Width/(2*c.Zoom)
	ty := -finalPos.Y + c.Height/(2*c.Zoom)

	// Matrice: [zoom, 0, 0, zoom, tx*zoom, ty*zoom]
	c.viewMatrix[0] = c.Zoom      // scaleX
	c.viewMatrix[1] = 0           // skewY
	c.viewMatrix[2] = 0           // skewX
	c.viewMatrix[3] = c.Zoom      // scaleY
	c.viewMatrix[4] = tx * c.Zoom // translateX
	c.viewMatrix[5] = ty * c.Zoom // translateY
}

// ===============================
// TRANSFORMATION METHODS
// ===============================

// WorldToScreen convertit des coordonnées monde en coordonnées écran
func (c *Camera) WorldToScreen(worldPos core.Vector2) core.Vector2 {
	// Position finale avec shake
	finalCamPos := c.Position
	if c.Shake != nil && c.Shake.active {
		finalCamPos = finalCamPos.Add(c.Shake.offset)
	}

	// Transformation
	screenX := (worldPos.X - finalCamPos.X + c.Width/(2*c.Zoom)) * c.Zoom
	screenY := (worldPos.Y - finalCamPos.Y + c.Height/(2*c.Zoom)) * c.Zoom

	return core.Vector2{X: screenX, Y: screenY}
}

// ScreenToWorld convertit des coordonnées écran en coordonnées monde
func (c *Camera) ScreenToWorld(screenPos core.Vector2) core.Vector2 {
	// Position finale avec shake
	finalCamPos := c.Position
	if c.Shake != nil && c.Shake.active {
		finalCamPos = finalCamPos.Add(c.Shake.offset)
	}

	// Transformation inverse
	worldX := (screenPos.X / c.Zoom) + finalCamPos.X - c.Width/(2*c.Zoom)
	worldY := (screenPos.Y / c.Zoom) + finalCamPos.Y - c.Height/(2*c.Zoom)

	return core.Vector2{X: worldX, Y: worldY}
}

// GetViewBounds retourne les limites de ce qui est visible
func (c *Camera) GetViewBounds() core.Rectangle {
	halfWidth := (c.Width / c.Zoom) / 2
	halfHeight := (c.Height / c.Zoom) / 2

	return core.Rectangle{
		X:      c.Position.X - halfWidth,
		Y:      c.Position.Y - halfHeight,
		Width:  c.Width / c.Zoom,
		Height: c.Height / c.Zoom,
	}
}

// IsVisible vérifie si un rectangle est visible par la caméra
func (c *Camera) IsVisible(bounds core.Rectangle) bool {
	viewBounds := c.GetViewBounds()
	return bounds.Intersects(viewBounds)
}

// IsPointVisible vérifie si un point est visible par la caméra
func (c *Camera) IsPointVisible(point core.Vector2) bool {
	viewBounds := c.GetViewBounds()
	return viewBounds.Contains(point)
}

// ===============================
// CAMERA EFFECTS
// ===============================

// StartShake démarre un effet de tremblement
func (c *Camera) StartShake(intensity float64, duration time.Duration) {
	c.StartShakeWithFrequency(intensity, duration, 30.0) // 30 Hz par défaut
}

// StartShakeWithFrequency démarre un tremblement avec fréquence personnalisée
func (c *Camera) StartShakeWithFrequency(intensity float64, duration time.Duration, frequency float64) {
	if c.Shake == nil {
		c.Shake = &CameraShake{}
	}

	c.Shake.Intensity = intensity
	c.Shake.Duration = duration
	c.Shake.Frequency = frequency
	c.Shake.startTime = time.Now()
	c.Shake.currentTime = 0
	c.Shake.active = true
	c.Shake.offset = core.Vector2{X: 0, Y: 0}
}

// StopShake arrête immédiatement le tremblement
func (c *Camera) StopShake() {
	if c.Shake != nil {
		c.Shake.active = false
		c.Shake.offset = core.Vector2{X: 0, Y: 0}
		c.needUpdate = true
	}
}

// IsShaking retourne true si la caméra tremble actuellement
func (c *Camera) IsShaking() bool {
	return c.Shake != nil && c.Shake.active
}

// ===============================
// CAMERA ANIMATION
// ===============================

// MoveTo anime la caméra vers une position
func (c *Camera) MoveTo(targetPos core.Vector2, duration time.Duration) {
	// TODO: Implémenter une animation fluide vers la position
	// Pour l'instant, utilisation du système de target
	c.SetTarget(targetPos)
}

// ZoomTo anime le zoom vers une valeur
func (c *Camera) ZoomTo(targetZoom float64, duration time.Duration) {
	// TODO: Implémenter une animation de zoom fluide
	// Pour l'instant, changement direct
	c.SetZoom(targetZoom)
}

// ===============================
// UTILITY METHODS
// ===============================

// Reset remet la caméra à sa position et zoom par défaut
func (c *Camera) Reset() {
	c.Position = core.Vector2{X: 0, Y: 0}
	c.targetPosition = core.Vector2{X: 0, Y: 0}
	c.velocity = core.Vector2{X: 0, Y: 0}
	c.Zoom = 1.0
	c.StopShake()
	c.StopFollowing()
	c.needUpdate = true
}

// SetSize change la taille de la vue de la caméra
func (c *Camera) SetSize(width, height float64) {
	c.Width = width
	c.Height = height
	c.needUpdate = true
}

// GetCenter retourne le centre de la vue de la caméra
func (c *Camera) GetCenter() core.Vector2 {
	return c.Position
}

// GetSize retourne la taille de la vue
func (c *Camera) GetSize() (float64, float64) {
	return c.Width, c.Height
}

// GetZoomedSize retourne la taille effective avec le zoom
func (c *Camera) GetZoomedSize() (float64, float64) {
	return c.Width / c.Zoom, c.Height / c.Zoom
}

// GetViewMatrix retourne la matrice de transformation
func (c *Camera) GetViewMatrix() [6]float64 {
	return c.viewMatrix
}

// SetSmoothing définit le niveau de lissage du mouvement
func (c *Camera) SetSmoothing(smoothing float64) {
	if smoothing < 0 {
		smoothing = 0
	}
	c.smoothing = smoothing
}

// ===============================
// ADVANCED CAMERA FEATURES
// ===============================

// PanTo fait un panoramique vers une position avec vitesse contrôlée
func (c *Camera) PanTo(targetPos core.Vector2, speed float64) {
	direction := targetPos.Sub(c.Position).Normalize()
	c.velocity = direction.Mul(speed)
}

// LookAt fait regarder la caméra vers un point avec un décalage temporel
func (c *Camera) LookAt(targetPos core.Vector2, leadTime float64) {
	// Prédire où sera la cible dans leadTime secondes
	if c.Target != nil {
		// Si on suit une cible, essayer de prédire son mouvement
		type Moveable interface {
			GetVelocity() core.Vector2
		}

		if moveable, ok := c.Target.(Moveable); ok {
			velocity := moveable.GetVelocity()
			predictedPos := targetPos.Add(velocity.Mul(leadTime))
			c.SetTarget(predictedPos)
			return
		}
	}

	c.SetTarget(targetPos)
}

// ConstrainToTarget garde la caméra dans certaines limites par rapport à sa cible
func (c *Camera) ConstrainToTarget(maxDistance float64) {
	if c.Target == nil {
		return
	}

	type Positionable interface {
		GetPosition() core.Vector2
	}

	if positionable, ok := c.Target.(Positionable); ok {
		targetPos := positionable.GetPosition()
		distance := c.Position.Distance(targetPos)

		if distance > maxDistance {
			direction := c.Position.Sub(targetPos).Normalize()
			c.Position = targetPos.Add(direction.Mul(maxDistance))
			c.targetPosition = c.Position
			c.needUpdate = true
		}
	}
}
