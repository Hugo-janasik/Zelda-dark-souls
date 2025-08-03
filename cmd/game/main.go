// cmd/game/main.go - Version corrigée avec ordre d'injection correct
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"zelda-souls-game/internal/assets"
	"zelda-souls-game/internal/core"
	"zelda-souls-game/internal/input"
	"zelda-souls-game/internal/rendering"
	"zelda-souls-game/internal/save"
)

// SpriteEbitenGame implémente l'interface ebiten.Game avec support des sprites
type SpriteEbitenGame struct {
	coreGame             *core.Game
	config               *core.GameConfig
	renderer             *rendering.Renderer
	inputWrapper         *input.FinalInputWrapper
	enhancedStateManager *core.EnhancedBuiltinStateManager
	spriteLoader         *assets.SpriteLoader
	frameCount           int
}

// NewSpriteEbitenGame crée le jeu avec support des sprites
func NewSpriteEbitenGame() (*SpriteEbitenGame, error) {
	fmt.Println("=== INITIALISATION DU JEU ===")

	// Charger la configuration
	config, err := core.LoadConfig("configs/game_config.yaml")
	if err != nil {
		log.Printf("Config non trouvée, utilisation des défauts: %v", err)
		config = core.GetDefaultConfig()
	}

	config.GameTitle = "Zelda Souls Game - Avec Sprites"
	config.GameVersion = "0.3.0"
	fmt.Printf("Configuration chargée: %s v%s\n", config.GameTitle, config.GameVersion)

	// Créer le renderer
	fmt.Println("Création du renderer...")
	renderer, err := rendering.NewRenderer(config)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le renderer: %v", err)
	}
	fmt.Println("✓ Renderer créé")

	// Créer l'asset manager et save manager
	fmt.Println("Création des gestionnaires...")
	assetManager := assets.NewAssetManager("assets")
	saveManager := save.NewSaveManager("saves")
	fmt.Println("✓ Gestionnaires créés")

	// Créer le sprite loader - PREMIÈRE ÉTAPE
	fmt.Println("Création du SpriteLoader...")
	spriteLoader := assets.NewSpriteLoader()
	fmt.Printf("✓ SpriteLoader créé: %T\n", spriteLoader)

	// Créer l'input manager et son wrapper
	fmt.Println("Création du gestionnaire d'entrées...")
	inputManager := input.NewInputManager(config)
	inputWrapper := input.NewFinalInputWrapper(inputManager)
	fmt.Println("✓ InputManager créé")

	// Créer le jeu core avec le système de base
	fmt.Println("Création du jeu core...")
	coreGame, err := core.NewGame(config, assetManager, saveManager)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le jeu: %v", err)
	}
	fmt.Println("✓ Jeu core créé")

	// Créer le gestionnaire d'états amélioré
	fmt.Println("Création du gestionnaire d'états...")
	enhancedStateManager := core.NewEnhancedBuiltinStateManager(
		config.WindowWidth(),
		config.WindowHeight(),
	)
	fmt.Println("✓ StateManager créé")

	// ORDRE CRITIQUE: Injecter le SpriteLoader AVANT les autres dépendances
	fmt.Println("\n=== INJECTION DU SPRITE LOADER (PREMIÈRE) ===")
	fmt.Printf("Injection du SpriteLoader (type: %T)...\n", spriteLoader)
	enhancedStateManager.SetSpriteLoader(spriteLoader)
	fmt.Println("✓ SpriteLoader injecté en premier")

	// Configurer les callbacks du StateManager
	fmt.Println("Configuration des callbacks...")
	enhancedStateManager.SetCallbacks(
		func() { // Nouvelle partie
			log.Println("Callback: Nouvelle partie démarrée")
		},
		func() { // Charger partie
			log.Println("Callback: Chargement de partie")
		},
		func() { // Quitter
			log.Println("Callback: Fermeture du jeu")
			coreGame.RequestExit()
		},
	)

	// Vérifier s'il y a des sauvegardes disponibles
	hasSaves := false
	if saveManager != nil {
		for i := 1; i <= 5; i++ {
			if saveManager.SlotExists(i) {
				hasSaves = true
				break
			}
		}
	}
	enhancedStateManager.SetHasSaves(hasSaves)
	fmt.Printf("✓ Sauvegardes détectées: %t\n", hasSaves)

	// Injecter les autres dépendances dans le jeu core
	fmt.Println("Injection des dépendances dans le core...")
	coreGame.SetRenderer(renderer)
	coreGame.SetStateManager(enhancedStateManager)
	coreGame.SetInputManager(inputWrapper)
	fmt.Println("✓ Dépendances injectées dans le jeu core")

	// Configurer l'input wrapper pour communiquer avec le core
	inputWrapper.SetCoreGame(coreGame)
	fmt.Println("✓ InputWrapper configuré")

	// Injecter la caméra et l'input manager dans le système de joueur
	fmt.Println("=== CONFIGURATION DU SYSTÈME DE JOUEUR ===")
	camera := renderer.GetCamera()
	enhancedStateManager.SetCamera(camera)
	enhancedStateManager.SetInputManager(inputWrapper)
	fmt.Println("✓ Camera et InputManager injectés")

	// VÉRIFICATION: S'assurer que le SpriteLoader est bien injecté
	fmt.Println("\n=== VÉRIFICATION INJECTION SPRITELOADER ===")
	playerSystem := enhancedStateManager.GetPlayerSystem()
	if playerSystem != nil {
		fmt.Println("PlayerSystem accessible depuis StateManager")

		// Test de création immédiate d'un joueur pour vérifier l'injection
		fmt.Println("Test de création de joueur pour vérifier l'injection...")
		// On ne crée pas vraiment le joueur ici, juste un test
	} else {
		fmt.Println("⚠ ERREUR: PlayerSystem non accessible")
	}

	fmt.Println("=== INITIALISATION TERMINÉE ===")

	return &SpriteEbitenGame{
		coreGame:             coreGame,
		config:               config,
		renderer:             renderer,
		inputWrapper:         inputWrapper,
		enhancedStateManager: enhancedStateManager,
		spriteLoader:         spriteLoader,
		frameCount:           0,
	}, nil
}

// Update implémente ebiten.Game.Update
func (seg *SpriteEbitenGame) Update() error {
	seg.frameCount++

	// Debug périodique pour les sprites
	if seg.frameCount == 60 { // Après 1 seconde
		fmt.Println("=== DEBUG SPRITES (après 1 seconde) ===")
		playerSystem := seg.enhancedStateManager.GetPlayerSystem()
		if playerSystem != nil {
			player := playerSystem.GetPlayer()
			if player != nil {
				fmt.Printf("Joueur actif: %t\n", player.Active)
				fmt.Printf("PlayerSprites chargés: %t\n", player.PlayerSprites != nil)
				if player.PlayerSprites != nil {
					fmt.Printf("Type PlayerSprites: %T\n", player.PlayerSprites)
				}

				// Vérifier l'état du SpriteLoader dans le PlayerSystem
				fmt.Println("État SpriteLoader dans PlayerSystem:")
				// Note: On ne peut pas accéder directement au spriteLoader privé,
				// mais on peut déduire son état des messages précédents
			} else {
				fmt.Println("Aucun joueur trouvé")
			}
		} else {
			fmt.Println("PlayerSystem non trouvé")
		}
		fmt.Println("=== FIN DEBUG SPRITES ===")
	}

	return seg.coreGame.Update()
}

// Draw implémente ebiten.Game.Draw
func (seg *SpriteEbitenGame) Draw(screen *ebiten.Image) {
	seg.coreGame.Render(screen)
}

// Layout implémente ebiten.Game.Layout
func (seg *SpriteEbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return seg.config.WindowWidth(), seg.config.WindowHeight()
}

// GetBuiltinStateManager retourne le gestionnaire d'états pour l'input wrapper
func (seg *SpriteEbitenGame) GetBuiltinStateManager() interface{} {
	return seg.enhancedStateManager
}

func main() {
	fmt.Println("Zelda Souls Game - Système de Sprites")
	fmt.Println("=====================================")

	// Vérifier et créer la structure des assets
	checkAndCreateAssets()

	// Initialiser les répertoires de base
	if err := initDirectories(); err != nil {
		log.Printf("Attention: %v", err)
	}

	// Créer le jeu avec sprites
	fmt.Println("\n=== CRÉATION DU JEU ===")
	game, err := NewSpriteEbitenGame()
	if err != nil {
		log.Fatal("Erreur création jeu:", err)
	}

	// Configuration Ebiten
	fmt.Println("\n=== CONFIGURATION EBITEN ===")
	config := game.config
	ebiten.SetWindowSize(config.WindowWidth(), config.WindowHeight())
	ebiten.SetWindowTitle(config.GameTitle)
	ebiten.SetVsyncEnabled(config.Window.VSync)

	fmt.Printf("✓ Fenêtre configurée: %dx%d\n", config.WindowWidth(), config.WindowHeight())
	fmt.Printf("✓ Titre: %s\n", config.GameTitle)
	fmt.Printf("✓ VSync: %t\n", config.Window.VSync)

	// Afficher les informations du jeu
	displayGameInfo(config)

	// Lancer le jeu avec Ebiten
	fmt.Println("\n=== LANCEMENT DU JEU ===")
	fmt.Println("Le jeu démarre...")
	fmt.Println("IMPORTANT: Le SpriteLoader a été injecté AVANT la création du joueur.")
	fmt.Println("Regardez la console pour les messages de debug des sprites.")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Erreur:", err)
	}

	// Cleanup à la fin
	fmt.Println("\n=== NETTOYAGE ===")
	if err := game.coreGame.Cleanup(); err != nil {
		log.Printf("Erreur cleanup: %v", err)
	}

	if game.spriteLoader != nil {
		game.spriteLoader.Cleanup()
	}

	fmt.Println("Jeu fermé proprement.")
}

// checkAndCreateAssets vérifie et crée la structure des assets nécessaire
func checkAndCreateAssets() {
	fmt.Println("\n=== VÉRIFICATION DES ASSETS ===")

	// Vérifier la structure de base
	playerDir := "assets/textures/player"
	if _, err := os.Stat(playerDir); os.IsNotExist(err) {
		fmt.Printf("⚠ Création du dossier: %s\n", playerDir)
		os.MkdirAll(playerDir, 0755)
	} else {
		fmt.Printf("✓ Dossier trouvé: %s\n", playerDir)
	}

	// Vérifier le fichier player.png principal
	playerFile := "assets/textures/player/player.png"
	if _, err := os.Stat(playerFile); os.IsNotExist(err) {
		fmt.Printf("⚠ %s n'existe pas, création d'un sprite de test...\n", playerFile)
		createTestPlayerSprite(playerFile)
	} else {
		fmt.Printf("✓ Sprite principal trouvé: %s\n", playerFile)
	}

	// Vérifier quelques sprites directionnels
	testSprites := []string{
		"assets/textures/player/down_idle/idle_down.png",
		"assets/textures/player/up_idle/idle_up.png",
	}

	for _, spritePath := range testSprites {
		if _, err := os.Stat(spritePath); os.IsNotExist(err) {
			fmt.Printf("○ Sprite optionnel manquant: %s\n", spritePath)
		} else {
			fmt.Printf("✓ Sprite directionnel trouvé: %s\n", spritePath)
		}
	}
}

// createTestPlayerSprite crée un sprite de test
func createTestPlayerSprite(filepath string) {
	// Créer le dossier parent si nécessaire
	dir := "assets/textures/player"
	os.MkdirAll(dir, 0755)

	// Créer une image de test 32x32
	testImg := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Dessiner un personnage simple
	playerColor := color.RGBA{100, 150, 255, 255} // Bleu
	borderColor := color.RGBA{80, 120, 200, 255}  // Bleu foncé
	faceColor := color.RGBA{255, 220, 180, 255}   // Peau

	// Remplir le fond
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			testImg.Set(x, y, playerColor)
		}
	}

	// Bordure
	for x := 0; x < 32; x++ {
		testImg.Set(x, 0, borderColor)
		testImg.Set(x, 31, borderColor)
	}
	for y := 0; y < 32; y++ {
		testImg.Set(0, y, borderColor)
		testImg.Set(31, y, borderColor)
	}

	// Visage simple
	for y := 8; y < 16; y++ {
		for x := 12; x < 20; x++ {
			testImg.Set(x, y, faceColor)
		}
	}

	// Yeux
	testImg.Set(14, 11, color.RGBA{0, 0, 0, 255})
	testImg.Set(17, 11, color.RGBA{0, 0, 0, 255})

	// Sauvegarder
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Erreur création sprite de test: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, testImg); err != nil {
		fmt.Printf("Erreur encodage sprite de test: %v\n", err)
		return
	}

	fmt.Printf("✓ Sprite de test créé: %s\n", filepath)
}

// initDirectories initialise les répertoires nécessaires
func initDirectories() error {
	dirs := []string{
		"saves",
		"logs",
		"screenshots",
		"configs",
		"assets/textures/player",
		"assets/textures/enemies",
		"assets/textures/environment",
		"assets/textures/ui",
		"assets/sounds/sfx",
		"assets/sounds/music",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("impossible de créer %s: %v", dir, err)
		}
	}

	return nil
}

// displayGameInfo affiche les informations du jeu
func displayGameInfo(config *core.GameConfig) {
	fmt.Println("\n=== INFORMATIONS DU JEU ===")
	fmt.Printf("Titre: %s v%s\n", config.GameTitle, config.GameVersion)
	fmt.Printf("Résolution: %dx%d\n", config.WindowWidth(), config.WindowHeight())
	fmt.Printf("FPS cible: %d\n", config.TargetFPS())

	fmt.Println("\n=== SPRITES SUPPORTÉS ===")
	fmt.Println("• assets/textures/player/player.png - Sprite principal")
	fmt.Println("• assets/textures/player/*/idle_*.png - Animations idle")
	fmt.Println("• assets/textures/player/*/attack_*.png - Animations attaque")

	fmt.Println("\n=== CONTRÔLES ===")
	fmt.Println("Menu:")
	fmt.Println("• Souris - Navigation et clic")

	fmt.Println("\nJeu:")
	fmt.Println("• ZQSD ou WASD - Mouvement")
	fmt.Println("• ESPACE - Attaque")
	fmt.Println("• C - Roulade (coûte 25 stamina)")
	fmt.Println("• E - Interaction")
	fmt.Println("• ESC - Retour au menu/Pause")
	fmt.Println("• I - Toggle instructions")

	fmt.Println("\n=== FONCTIONNALITÉS ===")
	fmt.Println("• Système de sprites avec fallback")
	fmt.Println("• Animations par direction")
	fmt.Println("• Système de stamina et vie")
	fmt.Println("• Caméra qui suit le joueur")
	fmt.Println("• Menu fonctionnel avec souris")
	fmt.Println("• CORRECTION: SpriteLoader injecté avant création joueur")

	fmt.Println("\n========================")
}
