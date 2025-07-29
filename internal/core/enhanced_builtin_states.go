// internal/core/enhanced_builtin_states.go - StateManager amélioré sans imports cycliques
package core

import (
	"fmt"
	"log"
	"time"
	"zelda-souls-game/internal/ecs/components"
	"zelda-souls-game/internal/ecs/systems"
)

// EnhancedBuiltinStateManager gestionnaire d'états avec joueur intégré
type EnhancedBuiltinStateManager struct {
	// État de base
	currentState     GameStateType
	frameCount       int
	showInstructions bool
	screenWidth      int
	screenHeight     int

	// Menu intégré (réutilisé du système précédent)
	buttons []*Button

	// Système de joueur
	playerSystem *systems.PlayerSystem

	// Callbacks
	onNewGame  func()
	onLoadGame func()
	onQuitGame func()

	// Entrées souris
	mousePos     Vector2
	mousePressed bool

	// Statistiques de jeu
	gameStartTime time.Time
}

// NewEnhancedBuiltinStateManager crée un gestionnaire d'états amélioré
func NewEnhancedBuiltinStateManager(screenWidth, screenHeight int) *EnhancedBuiltinStateManager {
	esm := &EnhancedBuiltinStateManager{
		currentState:     "menu",
		frameCount:       0,
		showInstructions: true,
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		playerSystem:     systems.NewPlayerSystem(),
		gameStartTime:    time.Now(),
	}

	esm.createButtons()
	return esm
}

// createButtons crée les boutons du menu (identique au système précédent)
func (esm *EnhancedBuiltinStateManager) createButtons() {
	centerX := float64(esm.screenWidth) / 2
	startY := float64(esm.screenHeight) / 2
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
			esm.startNewGame()
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
			if esm.onLoadGame != nil {
				esm.onLoadGame()
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
			if esm.onQuitGame != nil {
				esm.onQuitGame()
			}
		},
	)
	quitBtn.NormalColor = Color{120, 50, 50, 255} // Rouge
	quitBtn.HoverColor = Color{150, 70, 70, 255}

	esm.buttons = []*Button{newGameBtn, loadGameBtn, quitBtn}
}

// SetCallbacks définit les callbacks externes
func (esm *EnhancedBuiltinStateManager) SetCallbacks(onNewGame, onLoadGame, onQuitGame func()) {
	esm.onNewGame = onNewGame
	esm.onLoadGame = onLoadGame
	esm.onQuitGame = onQuitGame
}

// SetHasSaves définit si des sauvegardes existent
func (esm *EnhancedBuiltinStateManager) SetHasSaves(hasSaves bool) {
	if len(esm.buttons) >= 2 {
		esm.buttons[1].SetEnabled(hasSaves) // Bouton "Charger Partie"
	}
}

// SetInputManager injecte le gestionnaire d'entrées dans le système joueur
func (esm *EnhancedBuiltinStateManager) SetInputManager(inputManager interface{}) {
	// Adaptation de l'interface
	if im, ok := inputManager.(systems.InputManager); ok {
		esm.playerSystem.SetInputManager(im)
	}
}

// SetCamera injecte la caméra
func (esm *EnhancedBuiltinStateManager) SetCamera(camera interface{}) {
	esm.playerSystem.SetCamera(camera)
}

// startNewGame démarre une nouvelle partie
func (esm *EnhancedBuiltinStateManager) startNewGame() {
	fmt.Println("=== NOUVELLE PARTIE ===")
	
	// Créer le joueur au centre de l'écran
	playerX := float64(esm.screenWidth) / 2
	playerY := float64(esm.screenHeight) / 2
	esm.playerSystem.CreatePlayer(playerX, playerY)
	
	// Changer vers l'état de jeu
	esm.ChangeState("gameplay")
	esm.gameStartTime = time.Now()
	
	// Callback externe
	if esm.onNewGame != nil {
		esm.onNewGame()
	}
	
	fmt.Println("Nouvelle partie démarrée!")
}

// UpdateMouseInput met à jour les entrées souris - CORRIGÉ
func (esm *EnhancedBuiltinStateManager) UpdateMouseInput(mouseX, mouseY int, mousePressed bool) {
	esm.mousePos = Vector2{float64(mouseX), float64(mouseY)}
	esm.mousePressed = mousePressed
	
	// Debug pour vérifier que les entrées souris arrivent bien
	if mousePressed {
		fmt.Printf("StateManager reçoit clic souris à (%.0f, %.0f)\n", esm.mousePos.X, esm.mousePos.Y)
	}
}

// Update met à jour l'état
func (esm *EnhancedBuiltinStateManager) Update(deltaTime time.Duration) error {
	esm.frameCount++

	// Mettre à jour selon l'état actuel
	switch esm.currentState {
	case "menu":
		esm.updateMenuState(deltaTime)
	case "gameplay":
		esm.updateGameplayState(deltaTime)
	case "pause":
		esm.updatePauseState(deltaTime)
	}

	return nil
}

// updateMenuState met à jour l'état menu - CORRIGÉ AVEC DEBUG
func (esm *EnhancedBuiltinStateManager) updateMenuState(deltaTime time.Duration) {
	// Debug pour vérifier l'état de la souris
	if esm.frameCount%60 == 0 { // Toutes les secondes environ
		fmt.Printf("Menu - Souris: pos(%.0f,%.0f) pressed=%t\n", 
			esm.mousePos.X, esm.mousePos.Y, esm.mousePressed)
	}
	
	// Mettre à jour les boutons
	for i, button := range esm.buttons {
		button.Update(esm.mousePos, esm.mousePressed)
		
		// Debug pour voir si les boutons détectent la souris
		if button.Contains(esm.mousePos) && esm.frameCount%30 == 0 {
			fmt.Printf("Souris survole le bouton %d (%s)\n", i, button.Text)
		}
	}
}

// updateGameplayState met à jour l'état de jeu
func (esm *EnhancedBuiltinStateManager) updateGameplayState(deltaTime time.Duration) {
	// Mettre à jour le système de joueur
	esm.playerSystem.Update(deltaTime)
	
	// Vérifier si le joueur est mort
	if !esm.playerSystem.IsPlayerAlive() {
		fmt.Println("Joueur mort - retour au menu")
		esm.ChangeState("menu")
		return
	}
}

// updatePauseState met à jour l'état de pause
func (esm *EnhancedBuiltinStateManager) updatePauseState(deltaTime time.Duration) {
	// En pause, on ne met pas à jour le joueur
}

// UpdateWithInput met à jour avec InputManager (nouvelle méthode)
func (esm *EnhancedBuiltinStateManager) UpdateWithInput(deltaTime time.Duration, inputManager InputManager) error {
	// Injecter l'InputManager si pas encore fait
	esm.SetInputManager(inputManager)
	
	// Mise à jour normale
	return esm.Update(deltaTime)
}

// Render rend l'état actuel
func (esm *EnhancedBuiltinStateManager) Render(renderer Renderer) error {
	switch esm.currentState {
	case "menu":
		esm.renderMenuState(renderer)
	case "gameplay":
		esm.renderGameplayState(renderer)
	case "pause":
		esm.renderPauseState(renderer)
	default:
		esm.renderMenuState(renderer)
	}
	return nil
}

// renderMenuState rend l'état menu
func (esm *EnhancedBuiltinStateManager) renderMenuState(renderer Renderer) {
	// Titre
	titleX := float64(esm.screenWidth)/2 - float64(len("ZELDA SOULS")*12)/2
	renderer.DrawText("ZELDA SOULS", Vector2{titleX, 100}, ColorYellow)

	// Sous-titre
	subtitle := "Adventure Awaits"
	subtitleX := float64(esm.screenWidth)/2 - float64(len(subtitle)*8)/2
	renderer.DrawText(subtitle, Vector2{subtitleX, 140}, Color{200, 200, 200, 255})

	// Boutons
	for _, button := range esm.buttons {
		button.Render(renderer)
	}

	// Instructions
	instructionY := float64(esm.screenHeight) - 50
	instruction := "Utilisez la souris pour naviguer"
	instrX := float64(esm.screenWidth)/2 - float64(len(instruction)*8)/2
	renderer.DrawText(instruction, Vector2{instrX, instructionY}, Color{150, 150, 150, 255})
}

// renderGameplayState rend l'état gameplay
func (esm *EnhancedBuiltinStateManager) renderGameplayState(renderer Renderer) {
	// Interface de jeu
	renderer.DrawText("=== JEU EN COURS ===", Vector2{10, 10}, ColorWhite)
	renderer.DrawText("ESC - Retour menu", Vector2{10, 30}, ColorGreen)

	if esm.showInstructions {
		renderer.DrawText("ZQSD/WASD - Mouvement", Vector2{10, 60}, ColorWhite)
		renderer.DrawText("ESPACE - Attaque", Vector2{10, 80}, ColorWhite)
		renderer.DrawText("C - Roulade", Vector2{10, 100}, ColorWhite)
		renderer.DrawText("E - Interaction", Vector2{10, 120}, ColorWhite)
		renderer.DrawText("I - Toggle instructions", Vector2{10, 140}, ColorWhite)
	}

	// Informations du joueur
	esm.renderPlayerInfo(renderer)

	// Rendre le joueur avec une adaptation d'interface
	rendererAdapter := &RendererAdapter{coreRenderer: renderer}
	esm.playerSystem.Render(rendererAdapter)

	// Stats de jeu
	esm.renderGameStats(renderer)
}

// renderPauseState rend l'état de pause
func (esm *EnhancedBuiltinStateManager) renderPauseState(renderer Renderer) {
	// Assombrir l'arrière-plan
	overlay := Rectangle{X: 0, Y: 0, Width: float64(esm.screenWidth), Height: float64(esm.screenHeight)}
	renderer.DrawRectangle(overlay, Color{0, 0, 0, 128}, true)

	// Menu de pause
	centerX := float64(esm.screenWidth) / 2
	centerY := float64(esm.screenHeight) / 2

	renderer.DrawText("=== PAUSE ===", Vector2{centerX - 60, centerY - 50}, ColorYellow)
	renderer.DrawText("ESC - Reprendre", Vector2{centerX - 70, centerY - 20}, ColorWhite)
	renderer.DrawText("Q - Retour menu", Vector2{centerX - 70, centerY}, ColorWhite)
}

// renderPlayerInfo affiche les informations du joueur
func (esm *EnhancedBuiltinStateManager) renderPlayerInfo(renderer Renderer) {
	if !esm.playerSystem.IsPlayerAlive() {
		renderer.DrawText("JOUEUR MORT", Vector2{10, 180}, ColorRed)
		return
	}

	// Position du joueur
	playerPos := esm.playerSystem.GetPlayerPosition()
	posText := fmt.Sprintf("Position: (%.0f, %.0f)", playerPos.X, playerPos.Y)
	renderer.DrawText(posText, Vector2{10, 180}, ColorYellow)

	// Santé et stamina
	health, maxHealth := esm.playerSystem.GetPlayerHealth()
	stamina, maxStamina := esm.playerSystem.GetPlayerStamina()
	
	healthText := fmt.Sprintf("Vie: %d/%d", health, maxHealth)
	staminaText := fmt.Sprintf("Stamina: %.0f/%.0f", stamina, maxStamina)
	
	renderer.DrawText(healthText, Vector2{10, 200}, ColorGreen)
	renderer.DrawText(staminaText, Vector2{10, 220}, ColorCyan)

	// État du mouvement
	player := esm.playerSystem.GetPlayer()
	if player != nil && player.Movement.IsMoving {
		dirText := fmt.Sprintf("Direction: %s", player.Movement.Direction.String())
		renderer.DrawText(dirText, Vector2{10, 240}, ColorYellow)
		
		velocityLength := player.Movement.Velocity.Length()
		velocityText := fmt.Sprintf("Vitesse: %.1f", velocityLength)
		renderer.DrawText(velocityText, Vector2{10, 260}, ColorWhite)
	}
}

// renderGameStats affiche les statistiques de jeu
func (esm *EnhancedBuiltinStateManager) renderGameStats(renderer Renderer) {
	// Temps de jeu
	gameTime := time.Since(esm.gameStartTime)
	timeText := fmt.Sprintf("Temps: %s", formatDuration(gameTime))
	
	// Frames
	frameText := fmt.Sprintf("Frames: %d", esm.frameCount)
	
	// Affichage en bas à droite
	rightX := float64(esm.screenWidth) - 150
	bottomY := float64(esm.screenHeight) - 60
	
	renderer.DrawText(timeText, Vector2{rightX, bottomY}, ColorGray)
	renderer.DrawText(frameText, Vector2{rightX, bottomY + 20}, ColorGray)
}

// formatDuration formate une durée en string lisible
func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// GetCurrentStateType retourne le type d'état actuel
func (esm *EnhancedBuiltinStateManager) GetCurrentStateType() GameStateType {
	return GameStateType(esm.currentState)
}

// ChangeState change l'état
func (esm *EnhancedBuiltinStateManager) ChangeState(stateType GameStateType) {
	oldState := esm.currentState
	esm.currentState = GameStateType(stateType)
	fmt.Printf("Changement d'état: %s -> %s\n", oldState, esm.currentState)
}

// ToggleInstructions active/désactive les instructions
func (esm *EnhancedBuiltinStateManager) ToggleInstructions() {
	esm.showInstructions = !esm.showInstructions
	fmt.Printf("Instructions: %t\n", esm.showInstructions)
}

// GetPlayerSystem retourne le système de joueur
func (esm *EnhancedBuiltinStateManager) GetPlayerSystem() *systems.PlayerSystem {
	return esm.playerSystem
}

// IsInGame retourne si on est en jeu
func (esm *EnhancedBuiltinStateManager) IsInGame() bool {
	return esm.currentState == "gameplay"
}

// IsInMenu retourne si on est dans le menu
func (esm *EnhancedBuiltinStateManager) IsInMenu() bool {
	return esm.currentState == "menu"
}

// IsPaused retourne si le jeu est en pause
func (esm *EnhancedBuiltinStateManager) IsPaused() bool {
	return esm.currentState == "pause"
}

// ===============================
// ADAPTATEUR DE RENDERER
// ===============================

// RendererAdapter adapte le renderer core vers l'interface systems
type RendererAdapter struct {
	coreRenderer Renderer
}

// DrawRectangle adapte l'appel de rendu de rectangle vers components.Rectangle
func (r *RendererAdapter) DrawRectangle(rect components.Rectangle, color components.Color, filled bool) {
	coreRect := Rectangle{
		X:      rect.X,
		Y:      rect.Y,
		Width:  rect.Width,
		Height: rect.Height,
	}
	coreColor := Color{R: color.R, G: color.G, B: color.B, A: color.A}
	r.coreRenderer.DrawRectangle(coreRect, coreColor, filled)
}

// DrawText adapte l'appel de rendu de texte vers components.Vector2
func (r *RendererAdapter) DrawText(text string, pos components.Vector2, color components.Color) {
	corePos := Vector2{X: pos.X, Y: pos.Y}
	coreColor := Color{R: color.R, G: color.G, B: color.B, A: color.A}
	r.coreRenderer.DrawText(text, corePos, coreColor)
}