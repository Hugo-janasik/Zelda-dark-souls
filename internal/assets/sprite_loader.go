// internal/assets/sprite_loader.go - Système de chargement de sprites complet
package assets

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

	// Métadonnées
	SpriteWidth  int
	SpriteHeight int
	Loaded       bool
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

	fmt.Printf("✓ Image chargée: %s (%dx%d)\n", path, img.Bounds().Dx(), img.Bounds().Dy())
	return img, nil
}

// LoadPlayerSprites charge tous les sprites du joueur
func (sl *SpriteLoader) LoadPlayerSprites(assetsDir string) (*PlayerSpriteSet, error) {
	if sl.playerSprites != nil && sl.playerSprites.Loaded {
		fmt.Println("Sprites joueur déjà chargés, réutilisation du cache")
		return sl.playerSprites, nil
	}

	fmt.Println("=== CHARGEMENT DES SPRITES JOUEUR ===")

	playerSprites := &PlayerSpriteSet{}

	// 1. Charger le sprite principal
	mainSpritePath := filepath.Join(assetsDir, "textures", "player", "player.png")
	fmt.Printf("Tentative de chargement: %s\n", mainSpritePath)

	mainSprite, err := sl.LoadImage(mainSpritePath)
	if err != nil {
		fmt.Printf("⚠ Impossible de charger le sprite principal: %v\n", err)
		fmt.Println("Création d'un sprite de fallback...")
		mainSprite = sl.createFallbackSprite(32, 32)
	}
	playerSprites.MainSprite = mainSprite

	// Déterminer la taille du sprite
	bounds := mainSprite.Bounds()
	playerSprites.SpriteWidth = bounds.Dx()
	playerSprites.SpriteHeight = bounds.Dy()

	fmt.Printf("Sprite principal configuré: %dx%d\n", playerSprites.SpriteWidth, playerSprites.SpriteHeight)

	// 2. Créer les animations de base à partir du sprite principal
	fmt.Println("Création des animations de base...")

	baseAnimation := &SpriteAnimation{
		Frames: []image.Rectangle{
			image.Rect(0, 0, playerSprites.SpriteWidth, playerSprites.SpriteHeight),
		},
		FrameTime: 0.5,
		Loop:      true,
	}

	attackAnimation := &SpriteAnimation{
		Frames: []image.Rectangle{
			image.Rect(0, 0, playerSprites.SpriteWidth, playerSprites.SpriteHeight),
		},
		FrameTime: 0.2,
		Loop:      false,
	}

	// Assigner les animations (temporairement identiques)
	playerSprites.UpIdle = baseAnimation
	playerSprites.UpAttack = attackAnimation
	playerSprites.DownIdle = baseAnimation
	playerSprites.DownAttack = attackAnimation
	playerSprites.LeftIdle = baseAnimation
	playerSprites.LeftAttack = attackAnimation
	playerSprites.RightIdle = baseAnimation
	playerSprites.RightAttack = attackAnimation

	// 3. Essayer de charger les sprites spécifiques par direction (optionnel)
	sl.tryLoadDirectionalSprites(assetsDir, playerSprites)

	// 4. Marquer comme chargé
	playerSprites.Loaded = true
	sl.playerSprites = playerSprites

	fmt.Println("✓ Sprites du joueur chargés avec succès!")
	fmt.Printf("  - Sprite principal: %dx%d\n", playerSprites.SpriteWidth, playerSprites.SpriteHeight)
	fmt.Printf("  - Animations configurées: 8 (idle + attack pour 4 directions)\n")

	return playerSprites, nil
}

// tryLoadDirectionalSprites essaie de charger les sprites directionnels spécifiques
func (sl *SpriteLoader) tryLoadDirectionalSprites(assetsDir string, playerSprites *PlayerSpriteSet) {
	fmt.Println("Tentative de chargement des sprites directionnels...")

	// Liste des fichiers à essayer de charger
	spriteFiles := map[string]**SpriteAnimation{
		"up/idle_up.png":            &playerSprites.UpIdle,
		"up_idle/idle_up.png":       &playerSprites.UpIdle,
		"down/idle_down.png":        &playerSprites.DownIdle,
		"down_idle/idle_down.png":   &playerSprites.DownIdle,
		"left/idle_left.png":        &playerSprites.LeftIdle,
		"left_idle/idle_left.png":   &playerSprites.LeftIdle,
		"right/idle_right.png":      &playerSprites.RightIdle,
		"right_idle/idle_right.png": &playerSprites.RightIdle,
	}

	loadedCount := 0
	for relativePath, animationPtr := range spriteFiles {
		fullPath := filepath.Join(assetsDir, "textures", "player", relativePath)

		if sprite, err := sl.LoadImage(fullPath); err == nil {
			// Créer une animation avec ce sprite spécifique
			*animationPtr = &SpriteAnimation{
				Frames: []image.Rectangle{
					image.Rect(0, 0, sprite.Bounds().Dx(), sprite.Bounds().Dy()),
				},
				FrameTime: 0.5,
				Loop:      true,
			}
			loadedCount++
			fmt.Printf("  ✓ Sprite directionnel chargé: %s\n", relativePath)
		}
	}

	if loadedCount > 0 {
		fmt.Printf("✓ %d sprites directionnels chargés\n", loadedCount)
	} else {
		fmt.Println("  → Aucun sprite directionnel trouvé, utilisation du sprite principal")
	}
}

// createFallbackSprite crée un sprite de secours
func (sl *SpriteLoader) createFallbackSprite(width, height int) *ebiten.Image {
	fmt.Printf("Création d'un sprite de fallback %dx%d\n", width, height)

	img := ebiten.NewImage(width, height)

	// Dessiner un rectangle coloré avec des détails
	img.Fill(color.RGBA{100, 150, 255, 255}) // Bleu de base

	// Ajouter quelques détails pour le rendre reconnaissable
	// Dessiner une bordure
	borderColor := color.RGBA{80, 120, 200, 255}

	// Bordure simple (pixels du bord)
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)        // Haut
		img.Set(x, height-1, borderColor) // Bas
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)       // Gauche
		img.Set(width-1, y, borderColor) // Droite
	}

	// Ajouter un point central
	centerColor := color.RGBA{255, 255, 255, 255}
	centerX, centerY := width/2, height/2
	img.Set(centerX, centerY, centerColor)
	img.Set(centerX-1, centerY, centerColor)
	img.Set(centerX+1, centerY, centerColor)
	img.Set(centerX, centerY-1, centerColor)
	img.Set(centerX, centerY+1, centerColor)

	return img
}

// GetPlayerAnimation retourne l'animation appropriée selon l'état
func (pss *PlayerSpriteSet) GetPlayerAnimation(direction string, isAttacking bool) *SpriteAnimation {
	if !pss.Loaded {
		return nil
	}

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

// GetSpriteForAnimation retourne le sprite approprié pour une animation donnée
func (pss *PlayerSpriteSet) GetSpriteForAnimation(direction string, isMoving bool, isAttacking bool, frameIndex int) *ebiten.Image {
	if !pss.Loaded {
		return nil
	}

	// Pour l'instant, toujours retourner le sprite principal
	// TODO: Implémenter la sélection de frame dans les animations
	animation := pss.GetPlayerAnimation(direction, isAttacking)
	if animation != nil && len(animation.Frames) > 0 {
		// Pour l'instant, on retourne toujours le sprite principal
		// car toutes nos animations utilisent le même sprite de base
		return pss.MainSprite
	}

	return pss.MainSprite
}

// GetSpriteSize retourne la taille des sprites
func (pss *PlayerSpriteSet) GetSpriteSize() (int, int) {
	return pss.SpriteWidth, pss.SpriteHeight
}

// IsLoaded vérifie si les sprites sont chargés
func (pss *PlayerSpriteSet) IsLoaded() bool {
	return pss.Loaded && pss.MainSprite != nil
}

// GetMainSprite retourne le sprite principal
func (pss *PlayerSpriteSet) GetMainSprite() *ebiten.Image {
	return pss.MainSprite
}

// Cleanup libère les ressources
func (sl *SpriteLoader) Cleanup() {
	fmt.Println("Nettoyage SpriteLoader...")

	// Vider le cache d'images
	for path := range sl.loadedImages {
		delete(sl.loadedImages, path)
	}

	// Réinitialiser les sprites du joueur
	sl.playerSprites = nil

	fmt.Println("✓ SpriteLoader nettoyé")
}

// GetLoadedImageCount retourne le nombre d'images chargées
func (sl *SpriteLoader) GetLoadedImageCount() int {
	return len(sl.loadedImages)
}

// ReloadPlayerSprites force le rechargement des sprites du joueur
func (sl *SpriteLoader) ReloadPlayerSprites(assetsDir string) (*PlayerSpriteSet, error) {
	fmt.Println("Rechargement forcé des sprites du joueur...")
	sl.playerSprites = nil
	return sl.LoadPlayerSprites(assetsDir)
}

// CreateTestSprite crée un sprite de test pour le développement
func (sl *SpriteLoader) CreateTestSprite(width, height int, baseColor color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(baseColor)

	// Ajouter des détails de test
	borderColor := color.RGBA{
		R: uint8(int(baseColor.R) * 8 / 10),
		G: uint8(int(baseColor.G) * 8 / 10),
		B: uint8(int(baseColor.B) * 8 / 10),
		A: baseColor.A,
	}

	// Bordure
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)
		img.Set(x, height-1, borderColor)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)
		img.Set(width-1, y, borderColor)
	}

	return img
}
