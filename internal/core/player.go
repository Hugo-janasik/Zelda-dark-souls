// internal/core/player.go - Système de joueur
package core

import (
	"fmt"
	"time"
)

// Player représente le joueur
type Player struct {
	// Position et mouvement
	Position Vector2
	Velocity Vector2
	Speed    float64

	// Rendu
	Size  Vector2
	Color Color

	// État
	Moving    bool
	Direction Direction

	// Stats de base (pour plus tard)
	Health     int
	MaxHealth  int
	Stamina    float64
	MaxStamina float64
}

// NewPlayer crée un nouveau joueur
func NewPlayer(x, y float64) *Player {
	return &Player{
		Position:   Vector2{X: x, Y: y},
		Velocity:   Vector2{X: 0, Y: 0},
		Speed:      200.0,                     // pixels par seconde
		Size:       Vector2{X: 32, Y: 32},     // taille 32x32 pixels
		Color:      Color{100, 150, 255, 255}, // Bleu pour le joueur
		Moving:     false,
		Direction:  DirectionNone,
		Health:     100,
		MaxHealth:  100,
		Stamina:    100.0,
		MaxStamina: 100.0,
	}
}

// Update met à jour le joueur
func (p *Player) Update(deltaTime time.Duration, inputManager InputManager) {
	dt := deltaTime.Seconds()

	// Réinitialiser la vélocité
	p.Velocity = Vector2{X: 0, Y: 0}
	p.Moving = false
	p.Direction = DirectionNone

	// Gestion des entrées de mouvement
	if inputManager != nil {
		// Mouvement horizontal
		if p.isActionPressed(inputManager, int(ActionMoveLeft)) {
			p.Velocity.X = -p.Speed
			p.Direction = DirectionLeft
			p.Moving = true
		} else if p.isActionPressed(inputManager, int(ActionMoveRight)) {
			p.Velocity.X = p.Speed
			p.Direction = DirectionRight
			p.Moving = true
		}

		// Mouvement vertical
		if p.isActionPressed(inputManager, int(ActionMoveUp)) {
			p.Velocity.Y = -p.Speed
			if p.Direction == DirectionNone {
				p.Direction = DirectionUp
			}
			p.Moving = true
		} else if p.isActionPressed(inputManager, int(ActionMoveDown)) {
			p.Velocity.Y = p.Speed
			if p.Direction == DirectionNone {
				p.Direction = DirectionDown
			}
			p.Moving = true
		}

		// Normaliser la vélocité diagonale pour éviter que le joueur aille plus vite en diagonale
		if p.Velocity.X != 0 && p.Velocity.Y != 0 {
			length := p.Velocity.Length()
			if length > 0 {
				p.Velocity.X = (p.Velocity.X / length) * p.Speed
				p.Velocity.Y = (p.Velocity.Y / length) * p.Speed
			}
		}
	}

	// Appliquer le mouvement
	p.Position.X += p.Velocity.X * dt
	p.Position.Y += p.Velocity.Y * dt

	// Debug du mouvement
	if p.Moving {
		fmt.Printf("Joueur bouge: pos(%.1f,%.1f) dir=%s vel(%.1f,%.1f)\n",
			p.Position.X, p.Position.Y, p.directionString(), p.Velocity.X, p.Velocity.Y)
	}

	// Régénération de la stamina (pour plus tard)
	if p.Stamina < p.MaxStamina {
		p.Stamina += 25.0 * dt // Régénère 25 stamina par seconde
		if p.Stamina > p.MaxStamina {
			p.Stamina = p.MaxStamina
		}
	}
}

// isActionPressed vérifie si une action est pressée
func (p *Player) isActionPressed(inputManager InputManager, action int) bool {
	return inputManager.IsActionPressed(action)
}

// directionString retourne la direction en string pour le debug
func (p *Player) directionString() string {
	switch p.Direction {
	case DirectionUp:
		return "Haut"
	case DirectionDown:
		return "Bas"
	case DirectionLeft:
		return "Gauche"
	case DirectionRight:
		return "Droite"
	default:
		return "Aucune"
	}
}

// Render dessine le joueur
func (p *Player) Render(renderer Renderer) {
	// Dessiner le joueur comme un rectangle coloré pour l'instant
	playerRect := Rectangle{
		X:      p.Position.X - p.Size.X/2, // Centré sur la position
		Y:      p.Position.Y - p.Size.Y/2,
		Width:  p.Size.X,
		Height: p.Size.Y,
	}

	// Couleur différente selon l'état
	color := p.Color
	if p.Moving {
		// Légèrement plus clair quand il bouge
		color = Color{
			R: color.R + 30,
			G: color.G + 30,
			B: color.B + 30,
			A: color.A,
		}
	}

	// Dessiner le joueur
	renderer.DrawRectangle(playerRect, color, true)

	// Dessiner une bordure
	borderColor := Color{255, 255, 255, 255} // Blanc
	renderer.DrawRectangle(playerRect, borderColor, false)

	// Indicateur de direction (petite flèche)
	if p.Moving {
		p.drawDirectionIndicator(renderer)
	}
}

// drawDirectionIndicator dessine une flèche indiquant la direction
func (p *Player) drawDirectionIndicator(renderer Renderer) {
	centerX := p.Position.X
	centerY := p.Position.Y
	arrowSize := 10.0

	var arrowEnd Vector2

	switch p.Direction {
	case DirectionUp:
		arrowEnd = Vector2{centerX, centerY - arrowSize}
	case DirectionDown:
		arrowEnd = Vector2{centerX, centerY + arrowSize}
	case DirectionLeft:
		arrowEnd = Vector2{centerX - arrowSize, centerY}
	case DirectionRight:
		arrowEnd = Vector2{centerX + arrowSize, centerY}
	default:
		return
	}

	// Dessiner une ligne simple pour indiquer la direction
	// (On utilisera DrawRectangle pour faire une ligne épaisse)
	lineThickness := 2.0

	if p.Direction == DirectionUp || p.Direction == DirectionDown {
		// Ligne verticale
		lineRect := Rectangle{
			X:      centerX - lineThickness/2,
			Y:      min(centerY, arrowEnd.Y),
			Width:  lineThickness,
			Height: abs(arrowEnd.Y - centerY),
		}
		renderer.DrawRectangle(lineRect, Color{255, 255, 0, 255}, true) // Jaune
	} else {
		// Ligne horizontale
		lineRect := Rectangle{
			X:      min(centerX, arrowEnd.X),
			Y:      centerY - lineThickness/2,
			Width:  abs(arrowEnd.X - centerX),
			Height: lineThickness,
		}
		renderer.DrawRectangle(lineRect, Color{255, 255, 0, 255}, true) // Jaune
	}
}

// GetPosition retourne la position du joueur (pour la caméra)
func (p *Player) GetPosition() Vector2 {
	return p.Position
}

// GetVelocity retourne la vélocité du joueur (pour la caméra)
func (p *Player) GetVelocity() Vector2 {
	return p.Velocity
}

// SetPosition définit la position du joueur
func (p *Player) SetPosition(pos Vector2) {
	p.Position = pos
}

// GetBounds retourne les limites du joueur (pour les collisions futures)
func (p *Player) GetBounds() Rectangle {
	return Rectangle{
		X:      p.Position.X - p.Size.X/2,
		Y:      p.Position.Y - p.Size.Y/2,
		Width:  p.Size.X,
		Height: p.Size.Y,
	}
}

// Fonctions utilitaires
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
