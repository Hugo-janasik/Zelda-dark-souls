// internal/core/enhanced_builtin_states.go - StateManager amélioré avec sprites complet
package core

import (
	"fmt"
	"image"
	"log"
	"time"
	"zelda-souls-game/internal/assets"
	"zelda-souls-game/internal/ecs/components"
	"zelda-souls-game/internal/ecs/systems"

	"github.com/hajimehoshi/ebiten/v2"
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

	// Debug
	debugSprites bool
}

// NewEnhancedBuiltinStateManager crée un gestionnaire d'états amélioré
func NewEnhancedBuiltinStateManager(screenWidth, screenHeight int) *EnhancedBuiltinStateManager {
	fmt.Printf("Création EnhancedBuiltinStateManager (%dx%d)\n", screenWidth, screenHeight)

	esm := &EnhancedBuiltinStateManager{
		currentState:     "menu",
		frameCount:       0,
		showInstructions: true,
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		playerSystem:     systems.NewPlayerSystem(),
		gameStartTime:    time.Now(),
		debugSprites:     true,
	}

	esm.createButtons()
	fmt.Println("✓ EnhancedBuiltinStateManager créé")
	return esm
}

// createButtons crée les boutons du menu
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
	fmt.Printf("✓ %d boutons de menu créés\n", len(esm.buttons))
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
	if esm.debugSprites {
		fmt.Printf("SetInputManager appelé avec: %T\n", inputManager)
	}

	// Adaptation de l'interface
	if im, ok := inputManager.(systems.InputManager); ok {
		esm.playerSystem.SetInputManager(im)
		if esm.debugSprites {
			fmt.Println("✓ InputManager injecté dans PlayerSystem")
		}
	} else {
		fmt.Printf("⚠ InputManager type incompatible: %T\n", inputManager)
	}
}

// SetCamera injecte la caméra
func (esm *EnhancedBuiltinStateManager) SetCamera(camera interface{}) {
	if esm.debugSprites {
		fmt.Printf("SetCamera appelé avec: %T\n", camera)
	}

	esm.playerSystem.SetCamera(camera)

	if esm.debugSprites {
		fmt.Println("✓ Camera injectée dans PlayerSystem")
	}
}

// SetSpriteLoader injecte le chargeur de sprites dans le système de joueur
func (esm *EnhancedBuiltinStateManager) SetSpriteLoader(loader interface{}) {
	fmt.Printf("\n=== SetSpriteLoader appelé ===\n")
	fmt.Printf("Type du loader: %T\n", loader)

	// Vérifier que c'est bien notre SpriteLoader des assets
	if spriteLoader, ok := loader.(*assets.SpriteLoader); ok {
		fmt.Println("✓ Type SpriteLoader correct détecté")

		// Injecter dans le PlayerSystem
		if esm.playerSystem != nil {
			fmt.Println("PlayerSystem disponible, injection du SpriteLoader...")
			esm.playerSystem.SetSpriteLoader(spriteLoader)
			fmt.Println("✓ SpriteLoader injecté dans le PlayerSystem")

			// Vérifier si on a déjà un joueur créé
			if esm.playerSystem.GetPlayer() != nil {
				fmt.Println("Joueur déjà existant, sprites seront chargés automatiquement...")
			}
		} else {
			fmt.Println("⚠ ERREUR: PlayerSystem est nil!")
		}
	} else {
		fmt.Printf("⚠ Type inattendu pour SpriteLoader: %T\n", loader)
	}

	fmt.Println("=== Fin SetSpriteLoader ===\n")
}

// startNewGame démarre une nouvelle partie
func (esm *EnhancedBuiltinStateManager) startNewGame() {
	fmt.Println("\n=== DÉMARRAGE NOUVELLE PARTIE ===")

	// Créer le joueur au centre de l'écran
	playerX := float64(esm.screenWidth) / 2
	playerY := float64(esm.screenHeight) / 2

	fmt.Printf("Création du joueur à la position (%.1f, %.1f)\n", playerX, playerY)
	esm.playerSystem.CreatePlayer(playerX, playerY)

	// Vérifier que le joueur est bien créé
	if esm.playerSystem.GetPlayer() != nil {
		fmt.Println("✓ Joueur créé avec succès")

		player := esm.playerSystem.GetPlayer()
		fmt.Printf("  - Position: (%.1f, %.1f)\n", player.Position.Position.X, player.Position.Position.Y)
		fmt.Printf("  - Actif: %t\n", player.Active)
		fmt.Printf("  - Sprites: %t\n", player.PlayerSprites != nil)
		if player.PlayerSprites != nil {
			fmt.Printf("  - Type sprites: %T\n", player.PlayerSprites)
		}
	} else {
		fmt.Println("⚠ ERREUR: Échec de création du joueur")
	}

	// Changer vers l'état de jeu
	esm.ChangeState("gameplay")
	esm.gameStartTime = time.Now()

	// Callback externe
	if esm.onNewGame != nil {
		esm.onNewGame()
	}

	fmt.Println("✓ Nouvelle partie démarrée!")
	fmt.Println("=== FIN DÉMARRAGE NOUVELLE PARTIE ===\n")
}

// UpdateMouseInput met à jour les entrées souris
func (esm *EnhancedBuiltinStateManager) UpdateMouseInput(mouseX, mouseY int, mousePressed bool) {
	esm.mousePos = Vector2{float64(mouseX), float64(mouseY)}
	esm.mousePressed = mousePressed

	// Debug pour vérifier que les entrées souris arrivent bien
	if mousePressed && esm.frameCount%30 == 0 {
		fmt.Printf("Souris: pos(%.0f, %.0f) pressed=%t\n", esm.mousePos.X, esm.mousePos.Y, mousePressed)
	}
}

// Update met à jour l'état
func (esm *EnhancedBuiltinStateManager) Update(deltaTime time.Duration) error {
	esm.frameCount++

	// Debug sprites périodique
	if esm.debugSprites && esm.frameCount%600 == 0 { // Toutes les 10 secondes
		esm.debugSpriteState()
	}

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

// debugSpriteState affiche l'état des sprites pour debug
func (esm *EnhancedBuiltinStateManager) debugSpriteState() {
	fmt.Println("\n=== DEBUG ÉTAT SPRITES ===")
	fmt.Printf("Frame: %d, État: %s\n", esm.frameCount, esm.currentState)

	if esm.playerSystem != nil {
		player := esm.playerSystem.GetPlayer()
		if player != nil {
			fmt.Printf("Joueur actif: %t\n", player.Active)
			fmt.Printf("PlayerSprites: %t\n", player.PlayerSprites != nil)
			if player.PlayerSprites != nil {
				fmt.Printf("Type PlayerSprites: %T\n", player.PlayerSprites)
			}
			fmt.Printf("Position: (%.1f, %.1f)\n", player.Position.Position.X, player.Position.Position.Y)
			fmt.Printf("SpriteRenderer: %t\n", player.SpriteRenderer != nil)
			if player.SpriteRenderer != nil {
				fmt.Printf("  - Visible: %t\n", player.SpriteRenderer.Visible)
				fmt.Printf("  - Position: (%.1f, %.1f)\n", player.SpriteRenderer.Position.X, player.SpriteRenderer.Position.Y)
			}
		} else {
			fmt.Println("Aucun joueur")
		}
	} else {
		fmt.Println("PlayerSystem null")
	}
	fmt.Println("=== FIN DEBUG SPRITES ===\n")
}

// updateMenuState met à jour l'état menu
func (esm *EnhancedBuiltinStateManager) updateMenuState(deltaTime time.Duration) {
	// Debug souris périodique
	if esm.frameCount%180 == 0 { // Toutes les 3 secondes
		fmt.Printf("Menu - Souris: pos(%.0f,%.0f) pressed=%t\n",
			esm.mousePos.X, esm.mousePos.Y, esm.mousePressed)
	}

	// Mettre à jour les boutons
	for i, button := range esm.buttons {
		button.Update(esm.mousePos, esm.mousePressed)

		// Debug pour voir si les boutons détectent la souris
		if button.Contains(esm.mousePos) && esm.frameCount%60 == 0 {
			fmt.Printf("Souris survole le bouton %d (%s)\n", i, button.Text)
		}
	}
}

// updateGameplayState met à jour l'état de jeu
func (esm *EnhancedBuiltinStateManager) updateGameplayState(deltaTime time.Duration) {
	// Debug périodique du gameplay
	if esm.frameCount%300 == 0 { // Toutes les 5 secondes environ
		player := esm.playerSystem.GetPlayer()
		if player != nil {
			fmt.Printf("Gameplay - Joueur: pos(%.1f,%.1f), actif=%t, sprites=%t\n",
				player.Position.Position.X, player.Position.Position.Y,
				player.Active, player.PlayerSprites != nil)
		}
	}

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

	// Debug info sprites (si activé)
	if esm.debugSprites {
		debugText := fmt.Sprintf("Debug: Frame %d", esm.frameCount)
		renderer.DrawText(debugText, Vector2{10, float64(esm.screenHeight) - 30}, Color{100, 100, 100, 255})
	}
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

	// Debug sprites info
	if esm.debugSprites {
		esm.renderSpriteDebugInfo(renderer)
	}
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

// renderSpriteDebugInfo affiche les informations de debug des sprites
func (esm *EnhancedBuiltinStateManager) renderSpriteDebugInfo(renderer Renderer) {
	startY := 300.0

	player := esm.playerSystem.GetPlayer()
	if player == nil {
		renderer.DrawText("DEBUG: Aucun joueur", Vector2{10, startY}, ColorRed)
		return
	}

	// Informations sur les sprites
	debugTexts := []string{
		fmt.Sprintf("DEBUG SPRITES:"),
		fmt.Sprintf("PlayerSprites: %t", player.PlayerSprites != nil),
		fmt.Sprintf("SpriteRenderer: %t", player.SpriteRenderer != nil),
	}

	if player.PlayerSprites != nil {
		debugTexts = append(debugTexts, fmt.Sprintf("Type: %T", player.PlayerSprites))

		// Essayer de caster pour avoir plus d'infos
		if sprites, ok := player.PlayerSprites.(*systems.PlayerSpriteSet); ok {
			debugTexts = append(debugTexts,
				fmt.Sprintf("Loaded: %t", sprites.Loaded),
				fmt.Sprintf("MainSprite: %t", sprites.MainSprite != nil),
			)
		}
	}

	if player.SpriteRenderer != nil {
		debugTexts = append(debugTexts,
			fmt.Sprintf("Visible: %t", player.SpriteRenderer.Visible),
			fmt.Sprintf("Direction: %s", player.SpriteRenderer.LastDirection),
			fmt.Sprintf("Attacking: %t", player.SpriteRenderer.IsAttacking),
		)
	}

	for i, text := range debugTexts {
		color := ColorWhite
		if i == 0 {
			color = ColorYellow // Titre en jaune
		}
		renderer.DrawText(text, Vector2{10, startY + float64(i*15)}, color)
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

// ToggleDebugSprites active/désactive le debug des sprites
func (esm *EnhancedBuiltinStateManager) ToggleDebugSprites() {
	esm.debugSprites = !esm.debugSprites
	fmt.Printf("Debug sprites: %t\n", esm.debugSprites)
}

// ===============================
// ADAPTATEUR DE RENDERER - AVEC VRAIS SPRITES
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

// DrawSprite adapte l'appel de rendu de sprite avec support des vrais sprites
func (r *RendererAdapter) DrawSprite(sprite interface{}, position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, rotation float64, tint components.Color) {
	// Essayer de rendre avec un vrai sprite Ebiten
	if spriteImage, ok := sprite.(*ebiten.Image); ok && spriteImage != nil {
		// Utiliser le vrai système de sprites
		r.drawEbitenSprite(spriteImage, position, sourceRect, scale, rotation, tint)
		return
	}

	// Fallback vers un rectangle coloré si pas de sprite
	width := sourceRect.Width * scale.X
	height := sourceRect.Height * scale.Y

	rect := components.Rectangle{
		X:      position.X - width/2,
		Y:      position.Y - height/2,
		Width:  width,
		Height: height,
	}

	r.DrawRectangle(rect, tint, true)

	// Dessiner une bordure pour indiquer que c'est un fallback
	borderColor := components.Color{255, 255, 255, 100} // Blanc semi-transparent
	r.DrawRectangle(rect, borderColor, false)
}

// drawEbitenSprite dessine un sprite Ebiten réel
func (r *RendererAdapter) drawEbitenSprite(spriteImage *ebiten.Image, position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, rotation float64, tint components.Color) {
	// Vérifier si le renderer core supporte les sprites Ebiten
	if ebitenRenderer, ok := r.coreRenderer.(interface {
		GetMainImage() *ebiten.Image
	}); ok {
		mainImage := ebitenRenderer.GetMainImage()
		if mainImage != nil {
			// Dessiner directement sur l'image principale
			r.drawSpriteToEbitenImage(mainImage, spriteImage, position, sourceRect, scale, rotation, tint)
			return
		}
	}

	// Fallback si le renderer ne supporte pas Ebiten
	r.fallbackRectangleRender(position, sourceRect, scale, tint)
}

// drawSpriteToEbitenImage dessine le sprite sur l'image Ebiten
func (r *RendererAdapter) drawSpriteToEbitenImage(targetImage *ebiten.Image, spriteImage *ebiten.Image, position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, rotation float64, tint components.Color) {
	op := &ebiten.DrawImageOptions{}

	// Source rect (partie du sprite à dessiner)
	subImage := spriteImage
	if sourceRect.Width > 0 && sourceRect.Height > 0 {
		srcBounds := image.Rect(
			int(sourceRect.X),
			int(sourceRect.Y),
			int(sourceRect.X+sourceRect.Width),
			int(sourceRect.Y+sourceRect.Height),
		)
		subImage = spriteImage.SubImage(srcBounds).(*ebiten.Image)
	}

	// Scale
	op.GeoM.Scale(scale.X, scale.Y)

	// Rotation autour du centre
	if rotation != 0 {
		w := sourceRect.Width * scale.X
		h := sourceRect.Height * scale.Y
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(w/2, h/2)
	}

	// Position finale (centrer le sprite sur la position)
	finalWidth := sourceRect.Width * scale.X
	finalHeight := sourceRect.Height * scale.Y
	op.GeoM.Translate(position.X-finalWidth/2, position.Y-finalHeight/2)

	// Appliquer la teinte
	op.ColorM.Scale(
		float64(tint.R)/255.0,
		float64(tint.G)/255.0,
		float64(tint.B)/255.0,
		float64(tint.A)/255.0,
	)

	// Dessiner le sprite
	targetImage.DrawImage(subImage, op)
}

// fallbackRectangleRender rendu de fallback rectangulaire
func (r *RendererAdapter) fallbackRectangleRender(position components.Vector2, sourceRect components.Rectangle, scale components.Vector2, tint components.Color) {
	width := sourceRect.Width * scale.X
	height := sourceRect.Height * scale.Y

	rect := components.Rectangle{
		X:      position.X - width/2,
		Y:      position.Y - height/2,
		Width:  width,
		Height: height,
	}

	r.DrawRectangle(rect, tint, true)

	// Bordure pour indiquer le fallback
	borderColor := components.Color{255, 255, 255, 150}
	r.DrawRectangle(rect, borderColor, false)
}
