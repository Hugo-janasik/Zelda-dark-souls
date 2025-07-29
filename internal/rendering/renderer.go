// internal/rendering/renderer.go - Renderer basé sur Ebiten
package rendering

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"zelda-souls-game/internal/core"
)

// ===============================
// RENDERER STRUCTURE
// ===============================

// Renderer implémente le système de rendu avec Ebiten
type Renderer struct {
	// Configuration
	config *core.GameConfig
	width  int
	height int

	// Buffers de rendu
	mainImage  *ebiten.Image
	uiImage    *ebiten.Image
	debugImage *ebiten.Image

	// Gestion des textures
	textures     map[string]*ebiten.Image // Changé de core.TextureID à string
	textureCache map[string]*ebiten.Image

	// Batch rendering
	spriteBatch  *SpriteBatch
	drawCalls    int
	maxDrawCalls int

	// Caméra et viewport
	camera         *Camera
	viewportBounds core.Rectangle

	// Fonts et texte
	defaultFont font.Face
	fonts       map[string]font.Face

	// Debug info
	debugEnabled  bool
	showColliders bool
	showChunks    bool

	// Statistiques
	stats *RenderStats
}

// SpriteBatch optimise le rendu des sprites
type SpriteBatch struct {
	texture        *ebiten.Image
	vertices       []ebiten.Vertex
	indices        []uint16
	drawOptions    *ebiten.DrawTrianglesOptions
	currentTexture *ebiten.Image
	batchSize      int
	maxBatchSize   int
}

// RenderStats contient les statistiques de rendu
type RenderStats struct {
	DrawCalls      int
	TrianglesDrawn int
	SpritesDrawn   int
	TexturesUsed   int
	BatchesFlushed int
}

// ===============================
// RENDERER INITIALIZATION
// ===============================

// NewRenderer crée un nouveau renderer
func NewRenderer(config *core.GameConfig) (*Renderer, error) {
	renderer := &Renderer{
		config:        config,
		width:         config.WindowWidth(),
		height:        config.WindowHeight(),
		textures:      make(map[string]*ebiten.Image), // Changé
		textureCache:  make(map[string]*ebiten.Image),
		fonts:         make(map[string]font.Face),
		maxDrawCalls:  config.Rendering.MaxDrawCalls,
		debugEnabled:  config.Debug.EnableDebug,
		showColliders: config.Debug.ShowColliders,
		showChunks:    config.Debug.ShowChunkBorders,
		stats:         &RenderStats{},
	}

	// Initialiser les images de rendu
	renderer.mainImage = ebiten.NewImage(renderer.width, renderer.height)
	renderer.uiImage = ebiten.NewImage(renderer.width, renderer.height)
	renderer.debugImage = ebiten.NewImage(renderer.width, renderer.height)

	// Initialiser la caméra
	renderer.camera = NewCamera(
		core.Vector2{X: 0, Y: 0},
		float64(renderer.width),
		float64(renderer.height),
	)

	// Initialiser le batch de sprites
	renderer.spriteBatch = NewSpriteBatch(1000) // 1000 sprites max par batch

	// Charger la police par défaut
	renderer.defaultFont = basicfont.Face7x13
	renderer.fonts["default"] = renderer.defaultFont

	// Calculer le viewport
	renderer.updateViewport()

	return renderer, nil
}

// ===============================
// FRAME MANAGEMENT
// ===============================

// BeginFrame commence un nouveau frame de rendu
func (r *Renderer) BeginFrame() {
	// Réinitialiser les statistiques
	r.stats = &RenderStats{}
	r.drawCalls = 0

	// Vider les buffers
	r.mainImage.Clear()
	r.uiImage.Clear()
	if r.debugEnabled {
		r.debugImage.Clear()
	}

	// Commencer le batch
	r.spriteBatch.Begin()
}

// EndFrame termine le frame et affiche le résultat
func (r *Renderer) EndFrame() {
	// Terminer le batch
	r.spriteBatch.End()

	// Composer les couches finales
	r.composeFinalImage()
}

// Clear vide l'écran (méthode ajoutée pour compatibilité)
func (r *Renderer) Clear() {
	r.mainImage.Clear()
	r.uiImage.Clear()
	if r.debugEnabled {
		r.debugImage.Clear()
	}
}

// Present affiche le résultat final (méthode ajoutée pour compatibilité)
func (r *Renderer) Present() {
	// Dans notre cas, c'est géré par Ebiten automatiquement
	// Cette méthode existe pour la compatibilité avec l'interface
}

// composeFinalImage compose toutes les couches en une image finale
func (r *Renderer) composeFinalImage() {
	// L'image principale est déjà dans mainImage

	// Ajouter l'UI par dessus
	op := &ebiten.DrawImageOptions{}
	r.mainImage.DrawImage(r.uiImage, op)

	// Ajouter le debug en dernier
	if r.debugEnabled {
		r.mainImage.DrawImage(r.debugImage, op)
	}
}

// ===============================
// DRAWING METHODS
// ===============================

// DrawSprite dessine un sprite avec transformation
func (r *Renderer) DrawSprite(textureID string, position core.Vector2, options *DrawSpriteOptions) {
	texture := r.getTexture(textureID)
	if texture == nil {
		return
	}

	if options == nil {
		options = &DrawSpriteOptions{}
	}

	// Vérification de culling
	if r.config.Rendering.EnableCulling {
		bounds := core.Rectangle{
			X: position.X, Y: position.Y,
			Width:  float64(texture.Bounds().Dx()) * options.ScaleX,
			Height: float64(texture.Bounds().Dy()) * options.ScaleY,
		}

		if !r.isInViewport(bounds) {
			return
		}
	}

	// Utiliser le batch si possible
	if r.config.Rendering.EnableBatching {
		r.spriteBatch.DrawSprite(texture, position, options)
	} else {
		r.drawSpriteDirect(texture, position, options)
	}

	r.stats.SpritesDrawn++
}

// DrawTile dessine une tile de la carte
func (r *Renderer) DrawTile(textureID string, srcRect core.Rectangle, destRect core.Rectangle) {
	texture := r.getTexture(textureID)
	if texture == nil {
		return
	}

	// Culling de la tile
	if r.config.Rendering.EnableCulling && !r.isInViewport(destRect) {
		return
	}

	// Convertir en coordonnées écran
	screenPos := r.camera.WorldToScreen(core.Vector2{X: destRect.X, Y: destRect.Y})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenPos.X, screenPos.Y)

	// Créer une sous-image pour la partie source
	srcImage := texture.SubImage(image.Rect(
		int(srcRect.X), int(srcRect.Y),
		int(srcRect.X+srcRect.Width), int(srcRect.Y+srcRect.Height),
	)).(*ebiten.Image)

	r.mainImage.DrawImage(srcImage, op)
	r.drawCalls++
	r.stats.SpritesDrawn++
}

// DrawText dessine du texte à l'écran
func (r *Renderer) DrawText(textStr string, position core.Vector2, color core.Color) {
	clr := r.coreColorToEbiten(color)
	text.Draw(r.uiImage, textStr, r.defaultFont, int(position.X), int(position.Y), clr)
}

// DrawTextWithFont dessine du texte avec une police spécifique
func (r *Renderer) DrawTextWithFont(textStr string, position core.Vector2, fontName string, color core.Color) {
	font := r.fonts[fontName]
	if font == nil {
		font = r.defaultFont
	}

	clr := r.coreColorToEbiten(color)
	text.Draw(r.uiImage, textStr, font, int(position.X), int(position.Y), clr)
}

// DrawRectangle dessine un rectangle (pour debug principalement)
func (r *Renderer) DrawRectangle(rect core.Rectangle, color core.Color, filled bool) {
	clr := r.coreColorToEbiten(color)

	if filled {
		// Rectangle plein
		vector.DrawFilledRect(
			r.debugImage,
			float32(rect.X), float32(rect.Y),
			float32(rect.Width), float32(rect.Height),
			clr,
			false,
		)
	} else {
		// Contour seulement
		thickness := float32(1.0)
		// Haut
		vector.DrawFilledRect(r.debugImage,
			float32(rect.X), float32(rect.Y),
			float32(rect.Width), thickness, clr, false)
		// Bas
		vector.DrawFilledRect(r.debugImage,
			float32(rect.X), float32(rect.Y+rect.Height-float64(thickness)),
			float32(rect.Width), thickness, clr, false)
		// Gauche
		vector.DrawFilledRect(r.debugImage,
			float32(rect.X), float32(rect.Y),
			thickness, float32(rect.Height), clr, false)
		// Droite
		vector.DrawFilledRect(r.debugImage,
			float32(rect.X+rect.Width-float64(thickness)), float32(rect.Y),
			thickness, float32(rect.Height), clr, false)
	}
}

// DrawLine dessine une ligne
func (r *Renderer) DrawLine(start, end core.Vector2, color core.Color, thickness float32) {
	clr := r.coreColorToEbiten(color)
	vector.StrokeLine(
		r.debugImage,
		float32(start.X), float32(start.Y),
		float32(end.X), float32(end.Y),
		thickness, clr, false,
	)
}

// DrawCircle dessine un cercle
func (r *Renderer) DrawCircle(center core.Vector2, radius float32, color core.Color, filled bool) {
	clr := r.coreColorToEbiten(color)

	if filled {
		vector.DrawFilledCircle(
			r.debugImage,
			float32(center.X), float32(center.Y),
			radius, clr, false,
		)
	} else {
		vector.StrokeCircle(
			r.debugImage,
			float32(center.X), float32(center.Y),
			radius, 1.0, clr, false,
		)
	}
}

// ===============================
// TEXTURE MANAGEMENT
// ===============================

// LoadTexture charge une texture depuis un fichier
func (r *Renderer) LoadTexture(id string, filepath string) error {
	// Vérifier si déjà chargée
	if _, exists := r.textures[id]; exists {
		return nil
	}

	// Charger l'image depuis le cache si possible
	if cached, exists := r.textureCache[filepath]; exists {
		r.textures[id] = cached
		return nil
	}

	// Charger depuis le disque
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		return fmt.Errorf("impossible de charger la texture %s: %v", filepath, err)
	}

	// Stocker dans le cache et la map
	r.textureCache[filepath] = img
	r.textures[id] = img

	return nil
}

// UnloadTexture décharge une texture
func (r *Renderer) UnloadTexture(id string) {
	delete(r.textures, id)
}

// GetTexture retourne une texture chargée
func (r *Renderer) GetTexture(id string) *ebiten.Image {
	return r.textures[id]
}

// getTexture retourne une texture (interne)
func (r *Renderer) getTexture(id string) *ebiten.Image {
	return r.textures[id]
}

// ===============================
// CAMERA METHODS
// ===============================

// GetCamera retourne la caméra
func (r *Renderer) GetCamera() *Camera {
	return r.camera
}

// SetCameraPosition définit la position de la caméra
func (r *Renderer) SetCameraPosition(position core.Vector2) {
	r.camera.SetPosition(position)
	r.updateViewport()
}

// SetCameraZoom définit le zoom de la caméra
func (r *Renderer) SetCameraZoom(zoom float64) {
	r.camera.SetZoom(zoom)
	r.updateViewport()
}

// updateViewport met à jour les limites du viewport
func (r *Renderer) updateViewport() {
	r.viewportBounds = core.Rectangle{
		X:      r.camera.Position.X - r.camera.Width/2/r.camera.Zoom,
		Y:      r.camera.Position.Y - r.camera.Height/2/r.camera.Zoom,
		Width:  r.camera.Width / r.camera.Zoom,
		Height: r.camera.Height / r.camera.Zoom,
	}
}

// isInViewport vérifie si un rectangle est visible
func (r *Renderer) isInViewport(bounds core.Rectangle) bool {
	margin := r.config.Rendering.CullingMargin
	expandedViewport := core.Rectangle{
		X:      r.viewportBounds.X - margin,
		Y:      r.viewportBounds.Y - margin,
		Width:  r.viewportBounds.Width + 2*margin,
		Height: r.viewportBounds.Height + 2*margin,
	}

	return bounds.Intersects(expandedViewport)
}

// ===============================
// SPRITE BATCH IMPLEMENTATION
// ===============================

// NewSpriteBatch crée un nouveau batch de sprites
func NewSpriteBatch(maxSprites int) *SpriteBatch {
	return &SpriteBatch{
		vertices:     make([]ebiten.Vertex, 0, maxSprites*4),
		indices:      make([]uint16, 0, maxSprites*6),
		maxBatchSize: maxSprites,
		drawOptions:  &ebiten.DrawTrianglesOptions{},
	}
}

// Begin commence un nouveau batch
func (sb *SpriteBatch) Begin() {
	sb.vertices = sb.vertices[:0]
	sb.indices = sb.indices[:0]
	sb.currentTexture = nil
	sb.batchSize = 0
}

// DrawSprite ajoute un sprite au batch
func (sb *SpriteBatch) DrawSprite(texture *ebiten.Image, position core.Vector2, options *DrawSpriteOptions) {
	// Changer de batch si différente texture
	if sb.currentTexture != nil && sb.currentTexture != texture {
		sb.Flush()
	}

	sb.currentTexture = texture

	// Calculer les sommets du quad
	w := float64(texture.Bounds().Dx()) * options.ScaleX
	h := float64(texture.Bounds().Dy()) * options.ScaleY

	// Rotation et position
	cos := 1.0
	sin := 0.0
	if options.Rotation != 0 {
		cos = core.Cos(options.Rotation)
		sin = core.Sin(options.Rotation)
	}

	// Les 4 coins du sprite
	corners := []core.Vector2{
		{X: -w / 2, Y: -h / 2}, // Top-left
		{X: w / 2, Y: -h / 2},  // Top-right
		{X: w / 2, Y: h / 2},   // Bottom-right
		{X: -w / 2, Y: h / 2},  // Bottom-left
	}

	// Appliquer rotation et translation
	baseIdx := uint16(len(sb.vertices))
	for i, corner := range corners {
		// Rotation
		x := corner.X*cos - corner.Y*sin
		y := corner.X*sin + corner.Y*cos

		// Translation
		x += position.X
		y += position.Y

		// Coordonnées UV
		u := float32(0)
		v := float32(0)
		if i == 1 || i == 2 { // Right side
			u = 1
		}
		if i == 2 || i == 3 { // Bottom side
			v = 1
		}

		vertex := ebiten.Vertex{
			DstX:   float32(x),
			DstY:   float32(y),
			SrcX:   u * float32(texture.Bounds().Dx()),
			SrcY:   v * float32(texture.Bounds().Dy()),
			ColorR: float32(options.ColorR) / 255.0,
			ColorG: float32(options.ColorG) / 255.0,
			ColorB: float32(options.ColorB) / 255.0,
			ColorA: float32(options.ColorA) / 255.0,
		}

		sb.vertices = append(sb.vertices, vertex)
	}

	// Indices pour 2 triangles
	indices := []uint16{
		baseIdx, baseIdx + 1, baseIdx + 2,
		baseIdx, baseIdx + 2, baseIdx + 3,
	}
	sb.indices = append(sb.indices, indices...)

	sb.batchSize++

	// Flush si batch plein
	if sb.batchSize >= sb.maxBatchSize {
		sb.Flush()
	}
}

// Flush dessine tous les sprites du batch
func (sb *SpriteBatch) Flush() {
	if len(sb.vertices) == 0 || sb.currentTexture == nil {
		return
	}

	// Dessiner les triangles
	// Note: Cette partie nécessite une image cible, elle sera appelée par le renderer

	// Réinitialiser le batch
	sb.vertices = sb.vertices[:0]
	sb.indices = sb.indices[:0]
	sb.batchSize = 0
}

// End termine le batch
func (sb *SpriteBatch) End() {
	sb.Flush()
}

// ===============================
// DRAW OPTIONS
// ===============================

// DrawSpriteOptions options pour dessiner un sprite
type DrawSpriteOptions struct {
	ScaleX   float64
	ScaleY   float64
	Rotation float64
	ColorR   uint8
	ColorG   uint8
	ColorB   uint8
	ColorA   uint8
	FlipX    bool
	FlipY    bool
}

// NewDrawSpriteOptions crée des options par défaut
func NewDrawSpriteOptions() *DrawSpriteOptions {
	return &DrawSpriteOptions{
		ScaleX: 1.0,
		ScaleY: 1.0,
		ColorR: 255,
		ColorG: 255,
		ColorB: 255,
		ColorA: 255,
	}
}

// ===============================
// UTILITY METHODS
// ===============================

// drawSpriteDirect dessine un sprite directement (sans batch)
func (r *Renderer) drawSpriteDirect(texture *ebiten.Image, position core.Vector2, options *DrawSpriteOptions) {
	op := &ebiten.DrawImageOptions{}

	// Scale
	op.GeoM.Scale(options.ScaleX, options.ScaleY)

	// Rotation
	if options.Rotation != 0 {
		w := float64(texture.Bounds().Dx()) * options.ScaleX
		h := float64(texture.Bounds().Dy()) * options.ScaleY
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(options.Rotation)
		op.GeoM.Translate(w/2, h/2)
	}

	// Position (convertie par la caméra)
	screenPos := r.camera.WorldToScreen(position)
	op.GeoM.Translate(screenPos.X, screenPos.Y)

	// Couleur
	op.ColorM.Scale(
		float64(options.ColorR)/255.0,
		float64(options.ColorG)/255.0,
		float64(options.ColorB)/255.0,
		float64(options.ColorA)/255.0,
	)

	r.mainImage.DrawImage(texture, op)
	r.drawCalls++
}

// coreColorToEbiten convertit une couleur core en couleur Ebiten
func (r *Renderer) coreColorToEbiten(c core.Color) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

// SaveScreenshot sauvegarde une capture d'écran
func (r *Renderer) SaveScreenshot(filename string) error {
	// Créer le répertoire si nécessaire
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Créer le fichier
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encoder l'image
	return png.Encode(file, r.mainImage)
}

// GetStats retourne les statistiques de rendu
func (r *Renderer) GetStats() *RenderStats {
	return r.stats
}

// GetMainImage retourne l'image principale pour Ebiten
func (r *Renderer) GetMainImage() *ebiten.Image {
	return r.mainImage
}

// Cleanup nettoie les ressources
func (r *Renderer) Cleanup() {
	r.textures = nil
	r.textureCache = nil
	r.mainImage = nil
	r.uiImage = nil
	r.debugImage = nil
}
