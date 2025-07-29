// internal/ecs/systems/types.go - Types partagés pour éviter les imports cycliques
package systems

// ===============================
// TYPES GÉOMÉTRIQUES
// ===============================

// Vector2 représente un vecteur 2D (copié pour éviter les cycles)
type Vector2 struct {
	X, Y float64
}

// Rectangle représente un rectangle (copié pour éviter les cycles)
type Rectangle struct {
	X, Y, Width, Height float64
}

// Color représente une couleur RGBA (copié pour éviter les cycles)
type Color struct {
	R, G, B, A uint8
}

// ===============================
// CONSTANTES DE COULEURS
// ===============================

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
// FONCTIONS UTILITAIRES
// ===============================

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

// Length calcule la longueur du vecteur (au carré pour performance)
func (v Vector2) Length() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize normalise le vecteur
func (v Vector2) Normalize() Vector2 {
	lengthSq := v.X*v.X + v.Y*v.Y
	if lengthSq == 0 {
		return Vector2{0, 0}
	}
	// Approximation rapide de 1/sqrt(x)
	invLength := 1.0 / (lengthSq * 0.5)
	return Vector2{X: v.X * invLength, Y: v.Y * invLength}
}