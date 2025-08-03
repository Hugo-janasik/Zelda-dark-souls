// internal/assets/asset_manager.go - Gestionnaire d'assets
package assets

import (
	"path/filepath"
)

// TextureID représente l'identifiant d'une texture
type TextureID string

// SoundID représente l'identifiant d'un son
type SoundID string

// MapID représente l'identifiant d'une carte
type MapID string

// AssetManager gère le chargement et la mise en cache des ressources
type AssetManager struct {
	basePath       string
	loadedTextures map[TextureID]string
	loadedSounds   map[SoundID]string
	loadedMaps     map[MapID]string
	textureCount   int
	soundCount     int
}

// NewAssetManager crée un nouveau gestionnaire d'assets
func NewAssetManager(basePath string) *AssetManager {
	return &AssetManager{
		basePath:       basePath,
		loadedTextures: make(map[TextureID]string),
		loadedSounds:   make(map[SoundID]string),
		loadedMaps:     make(map[MapID]string),
	}
}

// LoadTexture charge une texture
func (am *AssetManager) LoadTexture(texturePath string) error {
	fullPath := filepath.Join(am.basePath, texturePath)
	textureID := TextureID(texturePath)
	am.loadedTextures[textureID] = fullPath
	am.textureCount++
	return nil
}

// LoadSound charge un son
func (am *AssetManager) LoadSound(soundPath string) error {
	fullPath := filepath.Join(am.basePath, soundPath)
	soundID := SoundID(soundPath)
	am.loadedSounds[soundID] = fullPath
	am.soundCount++
	return nil
}

// GetLoadedTextureCount retourne le nombre de textures chargées
func (am *AssetManager) GetLoadedTextureCount() int {
	return am.textureCount
}

// GetLoadedSoundCount retourne le nombre de sons chargés
func (am *AssetManager) GetLoadedSoundCount() int {
	return am.soundCount
}

// Cleanup nettoie les ressources
func (am *AssetManager) Cleanup() {
	am.loadedTextures = make(map[TextureID]string)
	am.loadedSounds = make(map[SoundID]string)
	am.textureCount = 0
	am.soundCount = 0
}
