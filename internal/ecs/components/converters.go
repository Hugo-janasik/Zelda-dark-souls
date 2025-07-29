// internal/ecs/components/converters.go - Fonctions de conversion pour éviter les imports cycliques
package components

// ===============================
// CONVERSIONS VERS SYSTEMS
// ===============================

// ToSystemsVector2 convertit vers le type systems.Vector2
func (v Vector2) ToSystemsVector2() interface{} {
	return struct {
		X, Y float64
	}{X: v.X, Y: v.Y}
}

// ToSystemsRectangle convertit vers le type systems.Rectangle
func (r Rectangle) ToSystemsRectangle() interface{} {
	return struct {
		X, Y, Width, Height float64
	}{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height}
}

// ToSystemsColor convertit vers le type systems.Color
func (c Color) ToSystemsColor() interface{} {
	return struct {
		R, G, B, A uint8
	}{R: c.R, G: c.G, B: c.B, A: c.A}
}

// ===============================
// CONVERSIONS DEPUIS SYSTEMS
// ===============================

// Vector2FromSystems convertit depuis le type systems
func Vector2FromSystems(v interface{}) Vector2 {
	switch val := v.(type) {
	case struct{ X, Y float64 }:
		return Vector2{X: val.X, Y: val.Y}
	default:
		return Vector2{X: 0, Y: 0}
	}
}

// RectangleFromSystems convertit depuis le type systems
func RectangleFromSystems(r interface{}) Rectangle {
	switch val := r.(type) {
	case struct{ X, Y, Width, Height float64 }:
		return Rectangle{X: val.X, Y: val.Y, Width: val.Width, Height: val.Height}
	default:
		return Rectangle{X: 0, Y: 0, Width: 0, Height: 0}
	}
}

// ColorFromSystems convertit depuis le type systems
func ColorFromSystems(c interface{}) Color {
	switch val := c.(type) {
	case struct{ R, G, B, A uint8 }:
		return Color{R: val.R, G: val.G, B: val.B, A: val.A}
	default:
		return ColorWhite
	}
}

// ===============================
// MÉTHODES UTILITAIRES SUPPLÉMENTAIRES
// ===============================

// Distance calcule la distance entre deux points
func (v Vector2) Distance(other Vector2) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy // Distance au carré pour performance
}

// Abs retourne la valeur absolue
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
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

// Min retourne le minimum de deux valeurs
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max retourne le maximum de deux valeurs
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}