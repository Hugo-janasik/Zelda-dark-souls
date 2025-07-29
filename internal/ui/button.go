// internal/ui/button.go - Système de boutons
package ui

// Vector2 pour éviter les imports
type Vector2 struct {
	X, Y float64
}

// Rectangle pour éviter les imports
type Rectangle struct {
	X, Y, Width, Height float64
}

// Color pour éviter les imports
type Color struct {
	R, G, B, A uint8
}

// ButtonState état d'un bouton
type ButtonState int

const (
	ButtonNormal ButtonState = iota
	ButtonHover
	ButtonPressed
	ButtonDisabled
)

// Button représente un bouton cliquable
type Button struct {
	// Position et taille
	Bounds Rectangle

	// Texte
	Text     string
	TextSize int

	// État
	State   ButtonState
	Enabled bool
	Visible bool

	// Callback
	OnClick func()

	// Style
	NormalColor   Color
	HoverColor    Color
	PressedColor  Color
	DisabledColor Color
	TextColor     Color

	// État interne
	wasPressed bool
}

// NewButton crée un nouveau bouton
func NewButton(x, y, width, height float64, text string, onClick func()) *Button {
	return &Button{
		Bounds: Rectangle{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Text:     text,
		TextSize: 16,
		State:    ButtonNormal,
		Enabled:  true,
		Visible:  true,
		OnClick:  onClick,

		// Couleurs par défaut
		NormalColor:   Color{70, 70, 70, 255},
		HoverColor:    Color{100, 100, 100, 255},
		PressedColor:  Color{50, 50, 50, 255},
		DisabledColor: Color{40, 40, 40, 255},
		TextColor:     Color{255, 255, 255, 255},
	}
}

// Contains vérifie si un point est dans le bouton
func (b *Button) Contains(point Vector2) bool {
	return point.X >= b.Bounds.X &&
		point.X <= b.Bounds.X+b.Bounds.Width &&
		point.Y >= b.Bounds.Y &&
		point.Y <= b.Bounds.Y+b.Bounds.Height
}

// Update met à jour l'état du bouton
func (b *Button) Update(mousePos Vector2, mousePressed bool) {
	if !b.Visible || !b.Enabled {
		b.State = ButtonDisabled
		return
	}

	isHovering := b.Contains(mousePos)

	if isHovering {
		if mousePressed && !b.wasPressed {
			b.State = ButtonPressed
			if b.OnClick != nil {
				b.OnClick()
			}
		} else if mousePressed {
			b.State = ButtonPressed
		} else {
			b.State = ButtonHover
		}
	} else {
		b.State = ButtonNormal
	}

	b.wasPressed = mousePressed
}

// Render dessine le bouton
func (b *Button) Render(renderer Renderer) {
	if !b.Visible {
		return
	}

	// Choisir la couleur selon l'état
	var bgColor Color
	switch b.State {
	case ButtonHover:
		bgColor = b.HoverColor
	case ButtonPressed:
		bgColor = b.PressedColor
	case ButtonDisabled:
		bgColor = b.DisabledColor
	default:
		bgColor = b.NormalColor
	}

	// Dessiner le fond du bouton
	renderer.DrawRectangle(b.Bounds, bgColor, true)

	// Dessiner la bordure
	borderColor := Color{200, 200, 200, 255}
	if b.State == ButtonDisabled {
		borderColor = Color{100, 100, 100, 255}
	}
	renderer.DrawRectangle(b.Bounds, borderColor, false)

	// Dessiner le texte centré
	textColor := b.TextColor
	if b.State == ButtonDisabled {
		textColor = Color{150, 150, 150, 255}
	}

	// Position du texte (approximativement centré)
	textX := b.Bounds.X + b.Bounds.Width/2 - float64(len(b.Text)*8)/2
	textY := b.Bounds.Y + b.Bounds.Height/2 - 8

	renderer.DrawText(b.Text, Vector2{textX, textY}, textColor)
}

// SetEnabled active/désactive le bouton
func (b *Button) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

// IsEnabled retourne si le bouton est activé
func (b *Button) IsEnabled() bool {
	return b.Enabled
}

// Renderer interface minimale
type Renderer interface {
	DrawText(text string, pos Vector2, color Color)
	DrawRectangle(rect Rectangle, color Color, filled bool)
}
