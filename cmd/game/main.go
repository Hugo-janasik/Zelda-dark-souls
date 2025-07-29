// cmd/game/main.go - Version avec support des sprites
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

// SpriteEbitenGame implémente l'interface ebiten.Game avec support des sprites
type SpriteEbitenGame struct {
	coreGame             *core.Game
	config               *core.GameConfig
	renderer             *rendering.Renderer
	inputWrapper         *input.FinalInputWrapper
	enhancedStateManager *core.EnhancedBuiltinStateManager
	spriteLoader         *assets.SpriteLoader // NOUVEAU
	frameCount           int
}

// NewSpriteEbitenGame crée le jeu avec support des sprites
func NewSpriteEbitenGame() (*SpriteEbitenGame, error) {
	// Charger la configuration
	config, err := core.LoadConfig("configs/game_config.yaml")
	if err != nil {
		log.Printf("Config non trouvée, utilisation des défauts: %v", err)
		config = core.GetDefaultConfig()
	}

	config.GameTitle = "Zelda Souls Game - Avec Sprites"
	config.GameVersion = "0.3.0"

	// Créer le renderer
	renderer, err := rendering.NewRenderer(config)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le renderer: %v", err)
	}

	// Créer l'asset manager et save manager
	assetManager := assets.NewAssetManager("assets")
	saveManager := save.NewSaveManager("saves")

	// Créer le sprite loader - NOUVEAU
	spriteLoader := assets.NewSpriteLoader()

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

	// Injecter les dépendances
	coreGame.SetRenderer(renderer)
	coreGame.SetStateManager(enhancedStateManager)
	coreGame.SetInputManager(inputWrapper)

	// Configurer l'input wrapper pour communiquer avec le core
	inputWrapper.SetCoreGame(coreGame)

	// Injecter la caméra et le sprite loader dans le système de joueur - NOUVEAU
	camera := renderer.GetCamera()
	enhancedStateManager.SetCamera(camera)
	enhancedStateManager.SetInputManager(inputWrapper)
	enhancedStateManager.SetSpriteLoader(spriteLoader) // NOUVEAU

	return &SpriteEbitenGame{
		coreGame:             coreGame,
		config:               config,
		renderer:             renderer,
		inputWrapper:         inputWrapper,
		enhancedStateManager: enhancedStateManager,
		spriteLoader:         spriteLoader, // NOUVEAU
		frameCount:           0,
	}, nil
}

// Update implémente ebiten.Game.Update
func (seg *SpriteEbitenGame) Update() error {
	seg.frameCount++
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

	// Initialiser les répertoires
	if err := initDirectories(); err != nil {
		log.Printf("Attention: %v", err)
	}

	// Créer le jeu avec sprites
	game, err := NewSpriteEbitenGame()
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
	fmt.Println("=== SPRITES SUPPORTÉS ===")
	fmt.Println("• assets/textures/player.png - Sprite principal")
	fmt.Println("• assets/textures/up_idle - Animation idle vers le haut")
	fmt.Println("• assets/textures/up_attack - Animation attaque vers le haut")
	fmt.Println("• assets/textures/down_idle - Animation idle vers le bas")
	fmt.Println("• assets/textures/down_attack - Animation attaque vers le bas")
	fmt.Println("• assets/textures/left_idle - Animation idle vers la gauche")
	fmt.Println("• assets/textures/left_attack - Animation attaque vers la gauche")
	fmt.Println("• assets/textures/right_idle - Animation idle vers la droite")
	fmt.Println("• assets/textures/right_attack - Animation attaque vers la droite")
	fmt.Println("")
	fmt.Println("=== CONTRÔLES ===")
	fmt.Println("Menu:")
	fmt.Println("• Souris - Navigation et clic")
	fmt.Println("")
	fmt.Println("Jeu:")
	fmt.Println("• ZQSD ou WASD - Mouvement (change l'animation)")
	fmt.Println("• ESPACE - Attaque (animation d'attaque)")
	fmt.Println("• C - Roulade (coûte 25 stamina)")
	fmt.Println("• E - Interaction")
	fmt.Println("• ESC - Retour au menu/Pause")
	fmt.Println("• I - Toggle instructions")
	fmt.Println("")
	fmt.Println("=== SPRITES ===")
	fmt.Println("• Le jeu chargera automatiquement les sprites depuis assets/textures/")
	fmt.Println("• Si un sprite n'est pas trouvé, un sprite de fallback sera utilisé")
	fmt.Println("• Les animations changent selon la direction et l'action")
	fmt.Println("• Redimensionnement automatique 2x pour une meilleure visibilité")
	fmt.Println("========================")
	fmt.Println("")

	// Lancer le jeu avec Ebiten
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Erreur:", err)
	}

	// Cleanup à la fin
	if err := game.coreGame.Cleanup(); err != nil {
		log.Printf("Erreur cleanup: %v", err)
	}
	
	// Cleanup des sprites
	if game.spriteLoader != nil {
		game.spriteLoader.Cleanup()
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