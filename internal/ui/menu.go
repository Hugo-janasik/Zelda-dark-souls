// internal/ui/menu.go - Système de menu principal
package ui

import (
	"log"
)

// MenuManager gère le menu principal
type MenuManager struct {
	buttons      []*Button
	title        string
	screenWidth  int
	screenHeight int

	// Callbacks
	OnNewGame  func()
	OnLoadGame func()
	OnQuitGame func()

	// État
	hasSaves bool
}

// NewMenuManager crée un nouveau gestionnaire de menu
func NewMenuManager(screenWidth, screenHeight int) *MenuManager {
	menu := &MenuManager{
		title:        "ZELDA SOULS",
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		buttons:      make([]*Button, 0),
		hasSaves:     false, // À vérifier plus tard
	}

	menu.createButtons()
	return menu
}

// createButtons crée les boutons du menu
func (m *MenuManager) createButtons() {
	centerX := float64(m.screenWidth) / 2
	startY := float64(m.screenHeight) / 2
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
			if m.OnNewGame != nil {
				m.OnNewGame()
			}
		},
	)

	// Bouton "Charger Partie"
	loadGameBtn := NewButton(
		centerX-buttonWidth/2,
		startY,
		buttonWidth,
		buttonHeight,
		"Charger Partie",
		func() {
			log.Println("Charger Partie cliquée")
			if m.OnLoadGame != nil {
				m.OnLoadGame()
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
			if m.OnQuitGame != nil {
				m.OnQuitGame()
			}
		},
	)

	// Couleurs spéciales pour certains boutons
	newGameBtn.NormalColor = Color{50, 120, 50, 255} // Vert
	newGameBtn.HoverColor = Color{70, 150, 70, 255}

	quitBtn.NormalColor = Color{120, 50, 50, 255} // Rouge
	quitBtn.HoverColor = Color{150, 70, 70, 255}

	m.buttons = append(m.buttons, newGameBtn, loadGameBtn, quitBtn)
}

// SetHasSaves définit si des sauvegardes existent
func (m *MenuManager) SetHasSaves(hasSaves bool) {
	m.hasSaves = hasSaves

	// Activer/désactiver le bouton "Charger Partie"
	if len(m.buttons) >= 2 {
		m.buttons[1].SetEnabled(hasSaves)
	}
}

// Update met à jour le menu
func (m *MenuManager) Update(mousePos Vector2, mousePressed bool) {
	for _, button := range m.buttons {
		button.Update(mousePos, mousePressed)
	}
}

// Render dessine le menu
func (m *MenuManager) Render(renderer Renderer) {
	// Dessiner le titre
	titleX := float64(m.screenWidth)/2 - float64(len(m.title)*12)/2
	titleY := 100.0
	renderer.DrawText(m.title, Vector2{titleX, titleY}, Color{255, 255, 100, 255})

	// Dessiner un sous-titre
	subtitle := "Adventure Awaits"
	subtitleX := float64(m.screenWidth)/2 - float64(len(subtitle)*8)/2
	subtitleY := 140.0
	renderer.DrawText(subtitle, Vector2{subtitleX, subtitleY}, Color{200, 200, 200, 255})

	// Dessiner les boutons
	for _, button := range m.buttons {
		button.Render(renderer)
	}

	// Instructions en bas
	instructionY := float64(m.screenHeight) - 50
	instruction := "Utilisez la souris pour naviguer"
	instrX := float64(m.screenWidth)/2 - float64(len(instruction)*8)/2
	renderer.DrawText(instruction, Vector2{instrX, instructionY}, Color{150, 150, 150, 255})
}

// SetCallbacks définit les callbacks du menu
func (m *MenuManager) SetCallbacks(onNewGame, onLoadGame, onQuitGame func()) {
	m.OnNewGame = onNewGame
	m.OnLoadGame = onLoadGame
	m.OnQuitGame = onQuitGame
}
