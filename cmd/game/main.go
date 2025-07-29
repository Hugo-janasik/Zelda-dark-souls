// cmd/game/main.go - Menu complet avec boutons cliquables
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

// EbitenGame implémente l'interface ebiten.Game
type EbitenGame struct {
	coreGame     *core.Game
	config       *core.GameConfig
	renderer     *rendering.Renderer
	inputWrapper *input.EnhancedInputWrapper
	frameCount   int
}

// NewEbitenGame crée le jeu avec menu complet
func NewEbitenGame() (*EbitenGame, error) {
	// Charger la configuration
	config, err := core.LoadConfig("configs/game_config.yaml")
	if err != nil {
		log.Printf("Config non trouvée, utilisation des défauts: %v", err)
		config = core.GetDefaultConfig()
	}

	config.GameTitle = "Zelda Souls Game"
	config.GameVersion = "0.1.0"

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
	inputWrapper := input.NewEnhancedInputWrapper(inputManager)

	// Créer le jeu core avec menu intégré
	coreGame, err := core.NewGameWithBuiltinStates(config, assetManager, saveManager)
	if err != nil {
		return nil, fmt.Errorf("impossible de créer le jeu: %v", err)
	}

	// Injecter les dépendances
	coreGame.SetRenderer(renderer)
	coreGame.SetInputManager(inputWrapper)

	// Configurer l'input wrapper pour communiquer avec le core
	inputWrapper.SetCoreGame(coreGame)

	return &EbitenGame{
		coreGame:     coreGame,
		config:       config,
		renderer:     renderer,
		inputWrapper: inputWrapper,
		frameCount:   0,
	}, nil
}

// Update implémente ebiten.Game.Update
func (eg *EbitenGame) Update() error {
	eg.frameCount++
	return eg.coreGame.Update()
}

// Draw implémente ebiten.Game.Draw
func (eg *EbitenGame) Draw(screen *ebiten.Image) {
	eg.coreGame.Render(screen)
}

// Layout implémente ebiten.Game.Layout
func (eg *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return eg.config.WindowWidth(), eg.config.WindowHeight()
}

func main() {
	fmt.Println("Zelda Souls Game - Menu interactif avec boutons cliquables")

	// Initialiser les répertoires
	if err := initDirectories(); err != nil {
		log.Printf("Attention: %v", err)
	}

	// Créer le jeu
	game, err := NewEbitenGame()
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
	fmt.Println("=== MENU PRINCIPAL ===")
	fmt.Println("• Cliquez sur 'Nouvelle Partie' pour commencer")
	fmt.Println("• 'Charger Partie' sera grisé s'il n'y a pas de sauvegardes")
	fmt.Println("• 'Quitter' ferme le jeu")
	fmt.Println("• En jeu: ESC pour revenir au menu")
	fmt.Println("======================")
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
		"assets",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("impossible de créer %s: %v", dir, err)
		}
	}

	return nil
}
