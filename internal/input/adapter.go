// internal/input/adapter.go - Adaptateur corrigé pour l'interface core
package input

// InputManagerWrapperFixed adapte InputManagerImpl pour l'interface core
type InputManagerWrapperFixed struct {
	inputManager *InputManagerImpl
}

// NewInputManagerWrapperFixed crée un wrapper
func NewInputManagerWrapperFixed(im *InputManagerImpl) *InputManagerWrapperFixed {
	return &InputManagerWrapperFixed{
		inputManager: im,
	}
}

// Update met à jour l'input manager
func (w *InputManagerWrapperFixed) Update() {
	w.inputManager.Update()
}

// IsKeyJustPressed avec un int (interface core)
func (w *InputManagerWrapperFixed) IsKeyJustPressed(key int) bool {
	return w.inputManager.IsKeyCorePressed(key)
}

// IsActionPressed avec un int (interface core)
func (w *InputManagerWrapperFixed) IsActionPressed(action int) bool {
	return w.inputManager.IsActionCorePressed(action)
}

// IsWindowCloseRequested (interface core)
func (w *InputManagerWrapperFixed) IsWindowCloseRequested() bool {
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