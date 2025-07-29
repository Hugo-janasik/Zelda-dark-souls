// cmd/game/main.go - Version finale sans imports cycliques
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"zelda-souls-game/internal/assets"
	"zelda-souls-game/internal/core"
	"zelda-souls-game/internal/input"
	"zelda-souls-game/internal/rendering"
	"zelda-souls-game/internal/save"
)

// FixedEbitenGame implémente l'interface ebiten.Game sans imports cycliques
type FixedEbitenGame struct {
	coreGame          *core.Game
	config            *core.GameConfig
	renderer          *rendering.Renderer
	inputWrapper      *input.FinalInputWrapper
	enhancedStateManager *core.EnhancedBuiltinStateManager
	frameCount        int
}

// NewFixedEbitenGame crée le jeu sans imports cycliques
func NewFixedEbitenGame() (*FixedEbitenGame, error) {
	// Charger la configuration
	config, err := core.LoadConfig("configs/game_config.yaml")
	if err != nil {
		log.Printf("Config non trouvée, utilisation des défauts: %v", err)
		config = core.GetDefaultConfig()
	}

	config.GameTitle = "Zelda Souls Game - Système Joueur"
	config.GameVersion = "0.2.0"

	// Créer le renderer
	renderer, err := rendering.NewRenderer(config)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le renderer: %v", err)
	}

	// Créer l'asset manager et save manager
	assetManager := assets.NewAssetManager("assets")
	saveManager := save.NewSaveManager("saves")

	// Créer l'input manager et son wrapper
	inputManager := input.NewInputManager(config)
	inputWrapper := input.NewFinalInputWrapper(inputManager)

	// Créer le jeu core avec le système de base
	coreGame, err := core.NewGame(config, assetManager, saveManager)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le jeu: %v", err)
	}

	// Créer le gestionnaire d'états amélioré
	enhancedStateManager := core.NewEnhancedBuiltinStateManager(
		config.WindowWidth(), 
		config.WindowHeight(),
	)

	// Configurer les callbacks du StateManager
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

	// Injecter les dépendances - CORRIGÉ SANS CAST D'INTERFACE
	coreGame.SetRenderer(renderer)
	
	// Ajouter la méthode SetEnhancedStateManager directement au core.Game
	// Pour l'instant, utilisons SetStateManager normal
	coreGame.SetStateManager(enhancedStateManager)
	
	coreGame.SetInputManager(inputWrapper)

	// Configurer l'input wrapper pour communiquer avec le core
	inputWrapper.SetCoreGame(coreGame)

	// Injecter la caméra dans le système de joueur
	camera := renderer.GetCamera()
	enhancedStateManager.SetCamera(camera)
	enhancedStateManager.SetInputManager(inputWrapper)

	return &FixedEbitenGame{
		coreGame:             coreGame,
		config:               config,
		renderer:             renderer,
		inputWrapper:         inputWrapper,
		enhancedStateManager: enhancedStateManager,
		frameCount:           0,
	}, nil
}

// Update implémente ebiten.Game.Update
func (feg *FixedEbitenGame) Update() error {
	feg.frameCount++
	return feg.coreGame.Update()
}

// Draw implémente ebiten.Game.Draw
func (feg *FixedEbitenGame) Draw(screen *ebiten.Image) {
	feg.coreGame.Render(screen)
}

// Layout implémente ebiten.Game.Layout
func (feg *FixedEbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return feg.config.WindowWidth(), feg.config.WindowHeight()
}

// GetBuiltinStateManager retourne le gestionnaire d'états pour l'input wrapper
func (feg *FixedEbitenGame) GetBuiltinStateManager() interface{} {
	return feg.enhancedStateManager // CORRIGÉ : retourne directement le bon StateManager
}

func main() {
	fmt.Println("Zelda Souls Game - Système de Joueur Complet")
	fmt.Println("===========================================")

	// Initialiser les répertoires
	if err := initDirectories(); err != nil {
		log.Printf("Attention: %v", err)
	}

	// Créer le jeu
	game, err := NewFixedEbitenGame()
	if err != nil {
		log.Fatal("Erreur création jeu:", err)
	}

	// Configuration Ebiten
	config := game.config
	ebiten.SetWindowSize(config.WindowWidth(), config.WindowHeight())
	ebiten.SetWindowTitle(config.GameTitle)
	ebiten.SetVsyncEnabled(config.Window.VSync)

	fmt.Printf("Lancement: %s v%s\n", config.GameTitle, config.GameVersion)
	fmt.Printf("Résolution: %dx%d\n", config.WindowWidth(), config.WindowHeight())
	fmt.Println("")
	fmt.Println("=== CONTRÔLES ===")
	fmt.Println("Menu:")
	fmt.Println("• Souris - Navigation et clic")
	fmt.Println("")
	fmt.Println("Jeu:")
	fmt.Println("• ZQSD ou WASD - Mouvement")
	fmt.Println("• ESPACE - Attaque (coûte 15 stamina)")
	fmt.Println("• C - Roulade (coûte 25 stamina)")
	fmt.Println("• E - Interaction")
	fmt.Println("• ESC - Retour au menu/Pause")
	fmt.Println("• I - Toggle instructions")
	fmt.Println("")
	fmt.Println("=== CARACTÉRISTIQUES ===")
	fmt.Println("• Joueur avec sprite et mouvement fluide")
	fmt.Println("• Système de stamina avec régénération")
	fmt.Println("• Caméra qui suit le joueur")
	fmt.Println("• Barres de vie et stamina")
	fmt.Println("• Indicateur de direction")
	fmt.Println("• Animations de base (idle/marche)")
	fmt.Println("• Système d'invulnérabilité")
	fmt.Println("• Limites d'écran temporaires")
	fmt.Println("• AUCUN IMPORT CYCLIQUE!")
	fmt.Println("=================")
	fmt.Println("")

	// Lancer le jeu avec Ebiten
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Erreur:", err)
	}

	// Cleanup à la fin
	if err := game.coreGame.Cleanup(); err != nil {
		log.Printf("Erreur cleanup: %v", err)
	}
}

func initDirectories() error {
	dirs := []string{
		"saves",
		"logs",
		"screenshots",
		"configs",
		"assets/textures/player",
		"assets/textures/enemies",
		"assets/textures/environment",
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