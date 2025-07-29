// internal/input/adapter.go - Adaptateur pour l'interface core
package input

// InputManagerWrapper adapte InputManager pour l'interface core
type InputManagerWrapper struct {
	inputManager *InputManager
}

// NewInputManagerWrapper crée un wrapper
func NewInputManagerWrapper(im *InputManager) *InputManagerWrapper {
	return &InputManagerWrapper{
		inputManager: im,
	}
}

// Update met à jour l'input manager
func (w *InputManagerWrapper) Update() {
	w.inputManager.Update()
}

// IsKeyJustPressed avec un int (interface core)
func (w *InputManagerWrapper) IsKeyJustPressed(key int) bool {
	return w.inputManager.IsKeyCorePressed(key)
}

// IsActionPressed avec un int (interface core)
func (w *InputManagerWrapper) IsActionPressed(action int) bool {
	return w.inputManager.IsActionCorePressed(action)
}

// IsWindowCloseRequested (interface core)
func (w *InputManagerWrapper) IsWindowCloseRequested() bool {
	return w.inputManager.IsWindowCloseRequested()
}

// Constantes pour la compatibilité avec core
const (
	KeyEscape = 27 // Code de la touche Escape
	KeyW      = 87 // Code de la touche W
	KeyS      = 83 // Code de la touche S
	KeyA      = 65 // Code de la touche A
	KeyD      = 68 // Code de la touche D
)
