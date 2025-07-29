// internal/core/game.go - Jeu complet avec states intégrés
package core

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// ===============================
// INTERFACES MINIMALES
// ===============================

// AssetManager interface
type AssetManager interface {
	LoadTexture(texturePath string) error
	LoadSound(soundPath string) error
	GetLoadedTextureCount() int
	GetLoadedSoundCount() int
	Cleanup()
}

// SaveManager interface
type SaveManager interface {
	SaveGame(slotID int, gameData interface{}) error
	LoadGame(slotID int) (interface{}, error)
	SlotExists(slotID int) bool
}

// Renderer interface
type Renderer interface {
	BeginFrame()
	EndFrame()
	DrawText(text string, pos Vector2, color Color)
	DrawRectangle(rect Rectangle, color Color, filled bool)
	GetMainImage() *ebiten.Image
	Cleanup()
}

// StateManager interface
type StateManager interface {
	Update(deltaTime time.Duration) error
	Render(renderer Renderer) error
	GetCurrentStateType() GameStateType
	ChangeState(stateType GameStateType)
}

// InputManager interface
type InputManager interface {
	Update()
	IsKeyJustPressed(key int) bool
	IsActionPressed(action int) bool
	IsWindowCloseRequested() bool
}

// ===============================
// GAME STRUCTURE
// ===============================

// Game structure principale
type Game struct {
	Config       *GameConfig
	AssetManager AssetManager
	SaveManager  SaveManager

	// Composants de base
	renderer     Renderer
	stateManager StateManager
	inputManager InputManager

	// État
	Running       bool
	Paused        bool
	LastFrameTime time.Time
	DeltaTime     time.Duration

	// Stats
	FrameCount    uint64
	FPS           float64
	LastFPSUpdate time.Time
}

// NewGame crée un jeu minimal
func NewGame(config *GameConfig, assetManager AssetManager, saveManager SaveManager) (*Game, error) {
	game := &Game{
		Config:        config,
		AssetManager:  assetManager,
		SaveManager:   saveManager,
		Running:       true,
		Paused:        false,
		LastFrameTime: time.Now(),
		LastFPSUpdate: time.Now(),
	}

	log.Println("Jeu minimal initialisé")
	return game, nil
}

// NewGameWithBuiltinStates crée un jeu avec StateManager intégré
func NewGameWithBuiltinStates(config *GameConfig, assetManager AssetManager, saveManager SaveManager) (*Game, error) {
	game := &Game{
		Config:        config,
		AssetManager:  assetManager,
		SaveManager:   saveManager,
		Running:       true,
		Paused:        false,
		LastFrameTime: time.Now(),
		LastFPSUpdate: time.Now(),
	}

	// Créer le StateManager avec les dimensions d'écran
	builtinStateManager := NewBuiltinStateManager(config.WindowWidth(), config.WindowHeight())

	// Configurer les callbacks du StateManager
	builtinStateManager.SetCallbacks(
		func() { // Nouvelle partie
			log.Println("Callback: Nouvelle partie démarrée")
			// TODO: Initialiser une nouvelle partie
		},
		func() { // Charger partie
			log.Println("Callback: Chargement de partie")
			// TODO: Charger une partie sauvegardée
		},
		func() { // Quitter
			log.Println("Callback: Fermeture du jeu")
			game.RequestExit()
		},
	)

	// Vérifier s'il y a des sauvegardes disponibles
	hasSaves := false
	if saveManager != nil {
		// Vérifier les slots de sauvegarde
		for i := 1; i <= 5; i++ {
			if saveManager.SlotExists(i) {
				hasSaves = true
				break
			}
		}
	}
	builtinStateManager.SetHasSaves(hasSaves)

	game.stateManager = builtinStateManager

	log.Printf("Jeu avec menu complet initialisé (sauvegardes: %t)", hasSaves)
	return game, nil
}

// Update met à jour le jeu (interface Ebiten)
func (g *Game) Update() error {
	// Calculer le delta time
	now := time.Now()
	g.DeltaTime = now.Sub(g.LastFrameTime)
	g.LastFrameTime = now

	// Mettre à jour les entrées si disponible
	if g.inputManager != nil {
		g.inputManager.Update()

		// Actions globales basiques
		if g.inputManager.IsKeyJustPressed(27) { // Escape
			log.Println("Pause/Menu (stub)")
		}
	}

	// Mettre à jour l'état si disponible
	if g.stateManager != nil && !g.Paused {
		// Essayer d'utiliser UpdateWithInput si disponible
		if bsm, ok := g.stateManager.(*BuiltinStateManager); ok {
			if err := bsm.UpdateWithInput(g.DeltaTime, g.inputManager); err != nil {
				return fmt.Errorf("erreur mise à jour état: %v", err)
			}
		} else {
			// Fallback vers Update normal
			if err := g.stateManager.Update(g.DeltaTime); err != nil {
				return fmt.Errorf("erreur mise à jour état: %v", err)
			}
		}
	}

	// Mettre à jour les stats
	g.updateStats()

	return nil
}

// Render effectue le rendu (interface Ebiten)
func (g *Game) Render(screen *ebiten.Image) {
	// Rendu basique si renderer disponible
	if g.renderer != nil {
		g.renderer.BeginFrame()

		// Rendre l'état actuel
		if g.stateManager != nil {
			if err := g.stateManager.Render(g.renderer); err != nil {
				log.Printf("Erreur rendu état: %v", err)
			}
		} else {
			// Rendu de fallback
			g.renderer.DrawText("Zelda Souls Game", Vector2{100, 100}, ColorWhite)
			g.renderer.DrawText("Moteur en cours d'initialisation...", Vector2{100, 130}, ColorGray)
		}

		// Debug info
		if g.Config.Debug.ShowFPS {
			g.renderDebugInfo()
		}

		g.renderer.EndFrame()

		// Copier vers l'écran Ebiten
		if mainImage := g.renderer.GetMainImage(); mainImage != nil {
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(mainImage, op)
		}
	} else {
		// Rendu direct sur l'écran Ebiten sans renderer
		g.renderFallback(screen)
	}
}

// renderFallback rendu de secours sans renderer
func (g *Game) renderFallback(screen *ebiten.Image) {
	// Dessiner un rectangle bleu pour indiquer que le jeu fonctionne
	screen.Fill(ColorBlue.ToEbitenColor())
}

// renderDebugInfo affiche les infos de debug
func (g *Game) renderDebugInfo() {
	if g.renderer != nil && g.Config.Debug.ShowFPS {
		fpsText := fmt.Sprintf("FPS: %.1f", g.FPS)
		g.renderer.DrawText(fpsText, Vector2{10, 10}, ColorWhite)
	}
}

// updateStats met à jour les statistiques
func (g *Game) updateStats() {
	g.FrameCount++

	if time.Since(g.LastFPSUpdate) >= time.Second {
		g.FPS = float64(g.FrameCount) / time.Since(g.LastFPSUpdate).Seconds()
		g.FrameCount = 0
		g.LastFPSUpdate = time.Now()
	}
}

// ===============================
// SETTERS POUR INJECTION DE DÉPENDANCES
// ===============================

// SetRenderer injecte le renderer
func (g *Game) SetRenderer(renderer Renderer) {
	g.renderer = renderer
	log.Println("Renderer injecté")
}

// SetStateManager injecte le gestionnaire d'états
func (g *Game) SetStateManager(stateManager StateManager) {
	g.stateManager = stateManager
	log.Println("StateManager injecté")
}

// SetInputManager injecte le gestionnaire d'entrées
func (g *Game) SetInputManager(inputManager InputManager) {
	g.inputManager = inputManager
	log.Println("InputManager injecté")
}

// ===============================
// GETTERS
// ===============================

func (g *Game) IsRunning() bool        { return g.Running }
func (g *Game) IsPaused() bool         { return g.Paused }
func (g *Game) GetConfig() *GameConfig { return g.Config }
func (g *Game) GetFPS() float64        { return g.FPS }

// RequestExit demande l'arrêt
func (g *Game) RequestExit() {
	g.Running = false
	log.Println("Arrêt demandé")
}

// Cleanup nettoie les ressources
func (g *Game) Cleanup() error {
	log.Println("Nettoyage des ressources...")

	if g.AssetManager != nil {
		g.AssetManager.Cleanup()
	}

	if g.renderer != nil {
		g.renderer.Cleanup()
	}

	log.Println("Nettoyage terminé")
	return nil
}

// internal/core/game.go - Ajout des méthodes manquantes pour le StateManager
// Ajoute ces méthodes à la fin de ton fichier game.go existant


// Ajoute cette méthode à la fin de internal/core/game.go

// GetBuiltinStateManager retourne le StateManager actuel
func (g *Game) GetBuiltinStateManager() interface{} {
	return g.stateManager
}

// GetStateManager retourne le StateManager actuel
func (g *Game) GetStateManager() StateManager {
	return g.stateManager
}

// SetEnhancedStateManager définit spécifiquement un EnhancedBuiltinStateManager
func (g *Game) SetEnhancedStateManager(esm *EnhancedBuiltinStateManager) {
	g.stateManager = esm
}