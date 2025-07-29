// internal/ecs/components/sprite_renderer.go - Composant de rendu de sprites amélioré
package components

import (
	"image"
	"time"
)

// SpriteRendererComponent gère le rendu des sprites avec animations
type SpriteRendererComponent struct {
	// Sprite actuel
	CurrentSprite    *image.Image // Référence vers l'image Ebiten (interface{})
	SourceRect       Rectangle    // Rectangle source dans la texture
	
	// Animation
	CurrentAnimation *SpriteAnimationData
	AnimationTime    float64
	CurrentFrame     int
	IsPlaying        bool
	
	// Rendu
	Position         Vector2
	Scale            Vector2
	Rotation         float64
	FlipX            bool
	FlipY            bool
	Tint             Color
	Visible          bool
	Layer            int
	
	// État du joueur
	LastDirection    string
	IsAttacking      bool
	AttackTime       float64
}

// SpriteAnimationData représente une animation de sprite
type SpriteAnimationData struct {
	Frames       []Rectangle
	FrameDuration float64 // Durée de chaque frame en secondes
	Loop         bool
	Name         string
}

// NewSpriteRendererComponent crée un nouveau composant de rendu de sprites
func NewSpriteRendererComponent() *SpriteRendererComponent {
	return &SpriteRendererComponent{
		Scale:         Vector2{X: 1.0, Y: 1.0},
		Tint:          ColorWhite,
		Visible:       true,
		Layer:         0,
		IsPlaying:     true,
		LastDirection: "down",
		IsAttacking:   false,
	}
}

// SetAnimation définit l'animation actuelle
func (src *SpriteRendererComponent) SetAnimation(animation *SpriteAnimationData) {
	if src.CurrentAnimation != animation {
		src.CurrentAnimation = animation
		src.CurrentFrame = 0
		src.AnimationTime = 0
		src.IsPlaying = true
	}
}

// Update met à jour l'animation
func (src *SpriteRendererComponent) Update(deltaTime time.Duration) {
	if !src.IsPlaying || src.CurrentAnimation == nil || len(src.CurrentAnimation.Frames) == 0 {
		return
	}
	
	dt := deltaTime.Seconds()
	
	// Mettre à jour le temps d'attaque
	if src.IsAttacking {
		src.AttackTime += dt
		// L'attaque dure 0.3 secondes
		if src.AttackTime >= 0.3 {
			src.IsAttacking = false
			src.AttackTime = 0
		}
	}
	
	// Mettre à jour l'animation
	src.AnimationTime += dt
	
	if src.AnimationTime >= src.CurrentAnimation.FrameDuration {
		src.AnimationTime = 0
		src.CurrentFrame++
		
		// Vérifier si on a atteint la fin de l'animation
		if src.CurrentFrame >= len(src.CurrentAnimation.Frames) {
			if src.CurrentAnimation.Loop {
				src.CurrentFrame = 0
			} else {
				src.CurrentFrame = len(src.CurrentAnimation.Frames) - 1
				src.IsPlaying = false
				
				// Si c'était une attaque, revenir à idle
				if src.IsAttacking {
					src.IsAttacking = false
				}
			}
		}
	}
	
	// Mettre à jour le rectangle source
	if src.CurrentAnimation != nil && src.CurrentFrame < len(src.CurrentAnimation.Frames) {
		src.SourceRect = src.CurrentAnimation.Frames[src.CurrentFrame]
	}
}

// StartAttack démarre une animation d'attaque
func (src *SpriteRendererComponent) StartAttack() {
	src.IsAttacking = true
	src.AttackTime = 0
}

// GetCurrentAnimationName retourne le nom de l'animation actuelle basée sur l'état
func (src *SpriteRendererComponent) GetCurrentAnimationName() string {
	if src.IsAttacking {
		return src.LastDirection + "_attack"
	}
	return src.LastDirection + "_idle"
}

// SetDirection définit la direction et change l'animation si nécessaire
func (src *SpriteRendererComponent) SetDirection(direction string, isMoving bool) {
	src.LastDirection = direction
	
	// L'animation sera mise à jour par le système de rendu
	// en fonction de l'état actuel (attaque ou idle/mouvement)
}