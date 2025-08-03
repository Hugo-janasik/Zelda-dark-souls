// internal/assets/sprite_loader.go - Système de chargement de sprites
package assets

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SpriteSheet représente une feuille de sprites
type SpriteSheet struct {
	Image      *ebiten.Image
	TileWidth  int
	TileHeight int
	Columns    int
	Rows       int
}

// SpriteAnimation représente une animation de sprites
type SpriteAnimation struct {
	Frames    []image.Rectangle
	FrameTime float64 // Durée de chaque frame en secondes
	Loop      bool
}

// PlayerSpriteSet contient tous les sprites du joueur
type PlayerSpriteSet struct {
	// Sprites par direction et état
	UpIdle      *SpriteAnimation
	UpAttack    *SpriteAnimation
	DownIdle    *SpriteAnimation
	DownAttack  *SpriteAnimation
	LeftIdle    *SpriteAnimation
	LeftAttack  *SpriteAnimation
	RightIdle   *SpriteAnimation
	RightAttack *SpriteAnimation

	// Sprite principal
	MainSprite *ebiten.Image
}

// SpriteLoader gère le chargement des sprites
type SpriteLoader struct {
	loadedImages  map[string]*ebiten.Image
	playerSprites *PlayerSpriteSet
}

// NewSpriteLoader crée un nouveau chargeur de sprites
func NewSpriteLoader() *SpriteLoader {
	return &SpriteLoader{
		loadedImages: make(map[string]*ebiten.Image),
	}
}

// LoadImage charge une image depuis un fichier
func (sl *SpriteLoader) LoadImage(path string) (*ebiten.Image, error) {
	// Vérifier si déjà chargée
	if img, exists := sl.loadedImages[path]; exists {
		return img, nil
	}

	// Charger l'image
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("impossible de charger l'image %s: %v", path, err)
	}

	// Mettre en cache
	sl.loadedImages[path] = img

	fmt.Printf("Image chargée: %s (%dx%d)\n", path, img.Bounds().Dx(), img.Bounds().Dy())
	return img, nil
}

// LoadPlayerSprites charge tous les sprites du joueur
func (sl *SpriteLoader) LoadPlayerSprites(assetsDir string) (*PlayerSpriteSet, error) {
	if sl.playerSprites != nil {
		return sl.playerSprites, nil // Déjà chargé
	}

	fmt.Println("Chargement des sprites du joueur...")

	playerSprites := &PlayerSpriteSet{}

	// Charger le sprite principal
	mainSpritePath := filepath.Join(assetsDir, "textures/player", "player.png")
	mainSprite, err := sl.LoadImage(mainSpritePath)
	if err != nil {
		fmt.Printf("Attention: impossible de charger le sprite principal: %v\n", err)
		// Créer un sprite de fallback
		mainSprite = sl.createFallbackSprite(32, 32)
	}
	playerSprites.MainSprite = mainSprite

	// Charger les animations par direction et état
	directions := []struct {
		name      string
		animation **SpriteAnimation
	}{
		{"up_idle", &playerSprites.UpIdle},
		{"up_attack", &playerSprites.UpAttack},
		{"down_idle", &playerSprites.DownIdle},
		{"down_attack", &playerSprites.DownAttack},
		{"left_idle", &playerSprites.LeftIdle},
		{"left_attack", &playerSprites.LeftAttack},
		{"right_idle", &playerSprites.RightIdle},
		{"right_attack", &playerSprites.RightAttack},
	}

	for _, dir := range directions {
		animation, err := sl.loadPlayerAnimation(assetsDir, dir.name)
		if err != nil {
			fmt.Printf("Attention: %v - utilisation d'animation par défaut\n", err)
			*dir.animation = sl.createDefaultAnimation()
		} else {
			*dir.animation = animation
		}
	}

	sl.playerSprites = playerSprites
	fmt.Println("Sprites du joueur chargés avec succès!")
	return playerSprites, nil
}

// loadPlayerAnimation charge une animation spécifique du joueur
func (sl *SpriteLoader) loadPlayerAnimation(assetsDir, animationName string) (*SpriteAnimation, error) {
	// Construire le chemin du fichier
	spritePath := filepath.Join(assetsDir, "textures", animationName)

	// Charger l'image
	img, err := sl.LoadImage(spritePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de charger l'animation %s: %v", animationName, err)
	}

	// Déterminer le nombre de frames selon le nom
	var frameCount int
	var frameTime float64

	if filepath.Ext(animationName) == "" {
		// C'est un dossier, chercher des fichiers individuels
		return sl.loadAnimationFromFiles(spritePath)
	}

	// C'est un fichier unique, supposer que c'est une spritesheet
	switch {
	case contains(animationName, "attack"):
		frameCount = 3  // 3 frames pour l'attaque
		frameTime = 0.1 // 100ms par frame
	case contains(animationName, "idle"):
		frameCount = 1  // 1 frame pour idle
		frameTime = 1.0 // 1 seconde
	default:
		frameCount = 1
		frameTime = 0.2
	}

	// Créer les rectangles des frames
	frameWidth := img.Bounds().Dx() / frameCount
	frameHeight := img.Bounds().Dy()

	frames := make([]image.Rectangle, frameCount)
	for i := 0; i < frameCount; i++ {
		frames[i] = image.Rect(
			i*frameWidth, 0,
			(i+1)*frameWidth, frameHeight,
		)
	}

	return &SpriteAnimation{
		Frames:    frames,
		FrameTime: frameTime,
		Loop:      contains(animationName, "idle"), // Loop seulement pour idle
	}, nil
}

// loadAnimationFromFiles charge une animation depuis plusieurs fichiers
func (sl *SpriteLoader) loadAnimationFromFiles(dirPath string) (*SpriteAnimation, error) {
	// Pour l'instant, retourner une animation par défaut
	// TODO: Implémenter le chargement depuis plusieurs fichiers
	return sl.createDefaultAnimation(), nil
}

// createDefaultAnimation crée une animation par défaut
func (sl *SpriteLoader) createDefaultAnimation() *SpriteAnimation {
	return &SpriteAnimation{
		Frames: []image.Rectangle{
			image.Rect(0, 0, 32, 32), // Frame unique
		},
		FrameTime: 0.2,
		Loop:      true,
	}
}

// createFallbackSprite crée un sprite de secours
func (sl *SpriteLoader) createFallbackSprite(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Dessiner un rectangle coloré simple
	img.Fill(color.RGBA{100, 150, 255, 255}) // Bleu

	return img
}

// GetPlayerAnimation retourne l'animation appropriée selon l'état
func (pss *PlayerSpriteSet) GetPlayerAnimation(direction string, isAttacking bool) *SpriteAnimation {
	switch direction {
	case "up":
		if isAttacking {
			return pss.UpAttack
		}
		return pss.UpIdle

	case "down":
		if isAttacking {
			return pss.DownAttack
		}
		return pss.DownIdle

	case "left":
		if isAttacking {
			return pss.LeftAttack
		}
		return pss.LeftIdle

	case "right":
		if isAttacking {
			return pss.RightAttack
		}
		return pss.RightIdle

	default:
		// Direction par défaut (down)
		if isAttacking {
			return pss.DownAttack
		}
		return pss.DownIdle
	}
}

// Fonction utilitaire pour vérifier si une string contient un substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && len(s) > 0 &&
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())
}

// Cleanup libère les ressources
func (sl *SpriteLoader) Cleanup() {
	sl.loadedImages = nil
	sl.playerSprites = nil
}
