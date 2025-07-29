// internal/core/builtin_states.go - StateManager avec menu intégré sans dépendances
package core

import (
	"fmt"
	"log"
	"time"
)

// Button structure intégrée dans core
type Button struct {
	Bounds  Rectangle
	Text    string
	Enabled bool
	Visible bool
	OnClick func()

	// État
	State      int // 0=normal, 1=hover, 2=pressed, 3=disabled
	wasPressed bool

	// Couleurs
	NormalColor   Color
	HoverColor    Color
	PressedColor  Color
	DisabledColor Color
	TextColor     Color
}

// NewButton crée un bouton intégré
func NewButton(x, y, width, height float64, text string, onClick func()) *Button {
	return &Button{
		Bounds:  Rectangle{X: x, Y: y, Width: width, Height: height},
		Text:    text,
		Enabled: true,
		Visible: true,
		OnClick: onClick,

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

// Update met à jour le bouton
func (b *Button) Update(mousePos Vector2, mousePressed bool) {
	if !b.Visible || !b.Enabled {
		b.State = 3 // disabled
		return
	}

	isHovering := b.Contains(mousePos)

	if isHovering {
		if mousePressed && !b.wasPressed {
			b.State = 2 // pressed
			if b.OnClick != nil {
				b.OnClick()
			}
		} else if mousePressed {
			b.State = 2 // pressed
		} else {
			b.State = 1 // hover
		}
	} else {
		b.State = 0 // normal
	}

	b.wasPressed = mousePressed
}

// Render dessine le bouton
func (b *Button) Render(renderer Renderer) {
	if !b.Visible {
		return
	}

	// Choisir la couleur
	var bgColor Color
	switch b.State {
	case 1:
		bgColor = b.HoverColor
	case 2:
		bgColor = b.PressedColor
	case 3:
		bgColor = b.DisabledColor
	default:
		bgColor = b.NormalColor
	}

	// Dessiner le fond
	renderer.DrawRectangle(b.Bounds, bgColor, true)

	// Dessiner la bordure
	borderColor := Color{200, 200, 200, 255}
	if b.State == 3 {
		borderColor = Color{100, 100, 100, 255}
	}
	renderer.DrawRectangle(b.Bounds, borderColor, false)

	// Dessiner le texte centré
	textColor := b.TextColor
	if b.State == 3 {
		textColor = Color{150, 150, 150, 255}
	}

	textX := b.Bounds.X + b.Bounds.Width/2 - float64(len(b.Text)*8)/2
	textY := b.Bounds.Y + b.Bounds.Height/2 - 8

	renderer.DrawText(b.Text, Vector2{textX, textY}, textColor)
}

// SetEnabled active/désactive le bouton
func (b *Button) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

// BuiltinStateManager gestionnaire d'états avec menu intégré
type BuiltinStateManager struct {
	currentState     GameStateType
	frameCount       int
	showInstructions bool

	// Menu intégré
	buttons      []*Button
	screenWidth  int
	screenHeight int

	// Joueur
	player *Player

	// Callbacks
	onNewGame  func()
	onLoadGame func()
	onQuitGame func()

	// Souris
	mousePos     Vector2
	mousePressed bool
}

// NewBuiltinStateManager crée un gestionnaire d'états avec menu
func NewBuiltinStateManager(screenWidth, screenHeight int) *BuiltinStateManager {
	bsm := &BuiltinStateManager{
		currentState:     "menu",
		frameCount:       0,
		showInstructions: true,
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
	}

	// Créer le joueur au centre de l'écran
	playerX := float64(screenWidth) / 2
	playerY := float64(screenHeight) / 2
	bsm.player = NewPlayer(playerX, playerY)

	bsm.createButtons()
	return bsm
}

// createButtons crée les boutons du menu
func (bsm *BuiltinStateManager) createButtons() {
	centerX := float64(bsm.screenWidth) / 2
	startY := float64(bsm.screenHeight) / 2
	buttonWidth := 200.0
	buttonHeight := 50.0
	buttonSpacing := 70.0

	// Bouton "Nouvelle Partie"
	newGameBtn := NewButton(
		centerX-buttonWidth/2,
		startY-buttonSpacing,
		buttonWidth,
		buttonHeight,
		"Nouvelle Partie",
		func() {
			log.Println("Nouvelle Partie cliquée")
			bsm.ChangeState("gameplay")
			if bsm.onNewGame != nil {
				bsm.onNewGame()
			}
		},
	)
	newGameBtn.NormalColor = Color{50, 120, 50, 255} // Vert
	newGameBtn.HoverColor = Color{70, 150, 70, 255}

	// Bouton "Charger Partie"
	loadGameBtn := NewButton(
		centerX-buttonWidth/2,
		startY,
		buttonWidth,
		buttonHeight,
		"Charger Partie",
		func() {
			log.Println("Charger Partie cliquée")
			if bsm.onLoadGame != nil {
				bsm.onLoadGame()
			}
		},
	)

	// Bouton "Quitter"
	quitBtn := NewButton(
		centerX-buttonWidth/2,
		startY+buttonSpacing,
		buttonWidth,
		buttonHeight,
		"Quitter",
		func() {
			log.Println("Quitter cliqué")
			if bsm.onQuitGame != nil {
				bsm.onQuitGame()
			}
		},
	)
	quitBtn.NormalColor = Color{120, 50, 50, 255} // Rouge
	quitBtn.HoverColor = Color{150, 70, 70, 255}

	bsm.buttons = []*Button{newGameBtn, loadGameBtn, quitBtn}
}

// SetCallbacks définit les callbacks externes
func (bsm *BuiltinStateManager) SetCallbacks(onNewGame, onLoadGame, onQuitGame func()) {
	bsm.onNewGame = onNewGame
	bsm.onLoadGame = onLoadGame
	bsm.onQuitGame = onQuitGame
}

// SetHasSaves définit si des sauvegardes existent
func (bsm *BuiltinStateManager) SetHasSaves(hasSaves bool) {
	if len(bsm.buttons) >= 2 {
		bsm.buttons[1].SetEnabled(hasSaves) // Bouton "Charger Partie"
	}
}

// UpdateMouseInput met à jour les entrées souris
func (bsm *BuiltinStateManager) UpdateMouseInput(mouseX, mouseY int, mousePressed bool) {
	bsm.mousePos = Vector2{float64(mouseX), float64(mouseY)}
	bsm.mousePressed = mousePressed
}

// Update met à jour l'état
func (bsm *BuiltinStateManager) Update(deltaTime time.Duration) error {
	bsm.frameCount++

	// Mettre à jour selon l'état actuel
	switch bsm.currentState {
	case "menu":
		for _, button := range bsm.buttons {
			button.Update(bsm.mousePos, bsm.mousePressed)
		}
	case "gameplay":
		// Mettre à jour le joueur
		// Note: on passera l'InputManager plus tard
		// bsm.player.Update(deltaTime, inputManager)
	}

	return nil
}

// UpdateWithInput met à jour avec InputManager (nouvelle méthode)
func (bsm *BuiltinStateManager) UpdateWithInput(deltaTime time.Duration, inputManager InputManager) error {
	bsm.frameCount++

	// Mettre à jour selon l'état actuel
	switch bsm.currentState {
	case "menu":
		for _, button := range bsm.buttons {
			button.Update(bsm.mousePos, bsm.mousePressed)
		}
	case "gameplay":
		// Mettre à jour le joueur avec les entrées
		if bsm.player != nil && inputManager != nil {
			bsm.player.Update(deltaTime, inputManager)
		}
	}

	return nil
}

// Render rend l'état actuel
func (bsm *BuiltinStateManager) Render(renderer Renderer) error {
	switch bsm.currentState {
	case "menu":
		bsm.renderMenuState(renderer)
	case "gameplay":
		bsm.renderGameplayState(renderer)
	default:
		bsm.renderMenuState(renderer)
	}
	return nil
}

// renderMenuState rend l'état menu
func (bsm *BuiltinStateManager) renderMenuState(renderer Renderer) {
	// Titre
	titleX := float64(bsm.screenWidth)/2 - float64(len("ZELDA SOULS")*12)/2
	renderer.DrawText("ZELDA SOULS", Vector2{titleX, 100}, ColorYellow)

	// Sous-titre
	subtitle := "Adventure Awaits"
	subtitleX := float64(bsm.screenWidth)/2 - float64(len(subtitle)*8)/2
	renderer.DrawText(subtitle, Vector2{subtitleX, 140}, Color{200, 200, 200, 255})

	// Boutons
	for _, button := range bsm.buttons {
		button.Render(renderer)
	}

	// Instructions
	instructionY := float64(bsm.screenHeight) - 50
	instruction := "Utilisez la souris pour naviguer"
	instrX := float64(bsm.screenWidth)/2 - float64(len(instruction)*8)/2
	renderer.DrawText(instruction, Vector2{instrX, instructionY}, Color{150, 150, 150, 255})
}

// renderGameplayState rend l'état gameplay
func (bsm *BuiltinStateManager) renderGameplayState(renderer Renderer) {
	// Dessiner le joueur
	if bsm.player != nil {
		bsm.player.Render(renderer)
	}

	// Interface de jeu
	renderer.DrawText("=== JEU EN COURS ===", Vector2{10, 10}, ColorWhite)
	renderer.DrawText("ESC - Retour menu", Vector2{10, 30}, ColorGreen)

	if bsm.showInstructions {
		renderer.DrawText("ZQSD - Mouvement", Vector2{10, 60}, ColorWhite)
		renderer.DrawText("I - Toggle instructions", Vector2{10, 80}, ColorWhite)
	}

	// Infos du joueur
	if bsm.player != nil {
		playerInfo := fmt.Sprintf("Joueur: (%.0f,%.0f)", bsm.player.Position.X, bsm.player.Position.Y)
		renderer.DrawText(playerInfo, Vector2{10, 120}, ColorYellow)

		if bsm.player.Moving {
			renderer.DrawText(fmt.Sprintf("Direction: %s", bsm.player.directionString()), Vector2{10, 140}, ColorYellow)
		}
	}

	frameText := fmt.Sprintf("Frames: %d", bsm.frameCount)
	renderer.DrawText(frameText, Vector2{10, 180}, ColorWhite)
}

// GetCurrentStateType retourne le type d'état actuel
func (bsm *BuiltinStateManager) GetCurrentStateType() GameStateType {
	return GameStateType(bsm.currentState)
}

// ChangeState change l'état
func (bsm *BuiltinStateManager) ChangeState(stateType GameStateType) {
	oldState := bsm.currentState
	bsm.currentState = GameStateType(stateType)
	fmt.Printf("Changement d'état: %s -> %s\n", oldState, bsm.currentState)
}

// ToggleInstructions active/désactive les instructions
func (bsm *BuiltinStateManager) ToggleInstructions() {
	bsm.showInstructions = !bsm.showInstructions
	fmt.Printf("Instructions: %t\n", bsm.showInstructions)
}
