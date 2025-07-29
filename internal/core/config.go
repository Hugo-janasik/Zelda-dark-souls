// internal/core/config.go - Configuration du jeu
package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// ===============================
// CONFIGURATION STRUCTURES
// ===============================

// GameConfig contient toute la configuration du jeu
type GameConfig struct {
	// Informations du jeu
	GameTitle   string `yaml:"game_title"`
	GameVersion string `yaml:"game_version"`

	// Configuration de la fenêtre
	Window WindowConfig `yaml:"window"`

	// Configuration du rendu
	Rendering RenderingConfig `yaml:"rendering"`

	// Configuration audio
	Audio AudioConfig `yaml:"audio"`

	// Configuration des contrôles
	Input InputConfig `yaml:"input"`

	// Configuration du gameplay
	Gameplay GameplayConfig `yaml:"gameplay"`

	// Configuration de débogage
	Debug DebugConfig `yaml:"debug"`

	// Chemins des ressources
	Paths PathsConfig `yaml:"paths"`
}

// WindowConfig configuration de la fenêtre
type WindowConfig struct {
	Width      int    `yaml:"width"`
	Height     int    `yaml:"height"`
	Title      string `yaml:"title"`
	Fullscreen bool   `yaml:"fullscreen"`
	Resizable  bool   `yaml:"resizable"`
	VSync      bool   `yaml:"vsync"`
	Maximized  bool   `yaml:"maximized"`
}

// RenderingConfig configuration du rendu
type RenderingConfig struct {
	TargetFPS      int     `yaml:"target_fps"`
	TileSize       int     `yaml:"tile_size"`
	ChunkSize      int     `yaml:"chunk_size"`
	MaxDrawCalls   int     `yaml:"max_draw_calls"`
	EnableBatching bool    `yaml:"enable_batching"`
	EnableCulling  bool    `yaml:"enable_culling"`
	CullingMargin  float64 `yaml:"culling_margin"`
	MaxTextures    int     `yaml:"max_textures"`

	// Effets visuels
	EnableParticles      bool `yaml:"enable_particles"`
	EnableLighting       bool `yaml:"enable_lighting"`
	EnableShadows        bool `yaml:"enable_shadows"`
	EnablePostProcessing bool `yaml:"enable_post_processing"`

	// Qualité
	TextureQuality  string `yaml:"texture_quality"` // "low", "medium", "high"
	ParticleQuality string `yaml:"particle_quality"`
}

// AudioConfig configuration audio
type AudioConfig struct {
	MasterVolume float64 `yaml:"master_volume"`
	MusicVolume  float64 `yaml:"music_volume"`
	SFXVolume    float64 `yaml:"sfx_volume"`
	VoiceVolume  float64 `yaml:"voice_volume"`
	EnableAudio  bool    `yaml:"enable_audio"`
	SampleRate   int     `yaml:"sample_rate"`
	BufferSize   int     `yaml:"buffer_size"`
	MaxSounds    int     `yaml:"max_sounds"`
}

// InputConfig configuration des contrôles
type InputConfig struct {
	KeyboardEnabled bool `yaml:"keyboard_enabled"`
	MouseEnabled    bool `yaml:"mouse_enabled"`
	GamepadEnabled  bool `yaml:"gamepad_enabled"`

	// Sensibilité
	MouseSensitivity   float64 `yaml:"mouse_sensitivity"`
	GamepadSensitivity float64 `yaml:"gamepad_sensitivity"`

	// Mapping des touches
	KeyMapping     map[string]string `yaml:"key_mapping"`
	GamepadMapping map[string]string `yaml:"gamepad_mapping"`

	// Zones mortes
	GamepadDeadzone float64 `yaml:"gamepad_deadzone"`
}

// GameplayConfig configuration du gameplay
type GameplayConfig struct {
	// Difficulté
	Difficulty            string  `yaml:"difficulty"` // "easy", "normal", "hard"
	DamageMultiplier      float64 `yaml:"damage_multiplier"`
	EnemyHealthMultiplier float64 `yaml:"enemy_health_multiplier"`

	// Progression
	ExperienceMultiplier float64 `yaml:"experience_multiplier"`
	SoulGainMultiplier   float64 `yaml:"soul_gain_multiplier"`

	// Système de stamina
	StaminaRegenRate float64 `yaml:"stamina_regen_rate"`
	StaminaPenalty   float64 `yaml:"stamina_penalty"`

	// Combat
	InvulnerabilityTime float64 `yaml:"invulnerability_time"`
	PerfectBlockWindow  float64 `yaml:"perfect_block_window"`

	// Monde
	EnemyRespawnTime float64 `yaml:"enemy_respawn_time"`
	ItemDespawnTime  float64 `yaml:"item_despawn_time"`

	// Sauvegarde
	AutoSaveEnabled  bool    `yaml:"auto_save_enabled"`
	AutoSaveInterval float64 `yaml:"auto_save_interval"` // en minutes
}

// DebugConfig configuration de débogage
type DebugConfig struct {
	EnableDebug      bool   `yaml:"enable_debug"`
	ShowFPS          bool   `yaml:"show_fps"`
	ShowColliders    bool   `yaml:"show_colliders"`
	ShowEntityInfo   bool   `yaml:"show_entity_info"`
	ShowChunkBorders bool   `yaml:"show_chunk_borders"`
	ShowPathfinding  bool   `yaml:"show_pathfinding"`
	EnableGodMode    bool   `yaml:"enable_god_mode"`
	EnableNoclip     bool   `yaml:"enable_noclip"`
	LogLevel         string `yaml:"log_level"` // "debug", "info", "warn", "error"
	ConsoleEnabled   bool   `yaml:"console_enabled"`
}

// PathsConfig chemins des ressources
type PathsConfig struct {
	AssetsDir      string `yaml:"assets_dir"`
	TexturesDir    string `yaml:"textures_dir"`
	SoundsDir      string `yaml:"sounds_dir"`
	MusicDir       string `yaml:"music_dir"`
	MapsDir        string `yaml:"maps_dir"`
	DataDir        string `yaml:"data_dir"`
	SavesDir       string `yaml:"saves_dir"`
	ConfigsDir     string `yaml:"configs_dir"`
	LogsDir        string `yaml:"logs_dir"`
	ScreenshotsDir string `yaml:"screenshots_dir"`
}

// ===============================
// CONFIGURATION METHODS
// ===============================

// LoadConfig charge la configuration depuis un fichier YAML
func LoadConfig(configPath string) (*GameConfig, error) {
	// Vérifier si le fichier existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("fichier de configuration introuvable: %s", configPath)
	}

	// Lire le fichier
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier de configuration: %v", err)
	}

	// Parser le YAML
	var config GameConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("erreur de parsing YAML: %v", err)
	}

	// Valider la configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration invalide: %v", err)
	}

	return &config, nil
}

// SaveConfig sauvegarde la configuration dans un fichier YAML
func (c *GameConfig) SaveConfig(configPath string) error {
	// Sérialiser en YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation YAML: %v", err)
	}

	// Écrire le fichier
	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("impossible d'écrire le fichier de configuration: %v", err)
	}

	return nil
}

// Validate valide la configuration
func (c *GameConfig) Validate() error {
	// Validation de la fenêtre
	if c.Window.Width <= 0 || c.Window.Height <= 0 {
		return fmt.Errorf("dimensions de fenêtre invalides: %dx%d", c.Window.Width, c.Window.Height)
	}

	// Validation du rendu
	if c.Rendering.TargetFPS <= 0 {
		return fmt.Errorf("FPS cible invalide: %d", c.Rendering.TargetFPS)
	}

	if c.Rendering.TileSize <= 0 {
		return fmt.Errorf("taille de tile invalide: %d", c.Rendering.TileSize)
	}

	// Validation audio
	if c.Audio.MasterVolume < 0.0 || c.Audio.MasterVolume > 1.0 {
		return fmt.Errorf("volume master invalide: %f", c.Audio.MasterVolume)
	}

	// Validation des chemins
	if c.Paths.AssetsDir == "" {
		return fmt.Errorf("répertoire d'assets non spécifié")
	}

	return nil
}

// GetDefaultConfig retourne une configuration par défaut
func GetDefaultConfig() *GameConfig {
	return &GameConfig{
		GameTitle:   "Zelda Souls Game",
		GameVersion: "0.1.0",

		Window: WindowConfig{
			Width:      1280,
			Height:     720,
			Title:      "Zelda Souls Game",
			Fullscreen: false,
			Resizable:  true,
			VSync:      true,
			Maximized:  false,
		},

		Rendering: RenderingConfig{
			TargetFPS:            60,
			TileSize:             32,
			ChunkSize:            16,
			MaxDrawCalls:         1000,
			EnableBatching:       true,
			EnableCulling:        true,
			CullingMargin:        100.0,
			MaxTextures:          256,
			EnableParticles:      true,
			EnableLighting:       false,
			EnableShadows:        false,
			EnablePostProcessing: false,
			TextureQuality:       "high",
			ParticleQuality:      "medium",
		},

		Audio: AudioConfig{
			MasterVolume: 1.0,
			MusicVolume:  0.7,
			SFXVolume:    0.8,
			VoiceVolume:  1.0,
			EnableAudio:  true,
			SampleRate:   44100,
			BufferSize:   1024,
			MaxSounds:    32,
		},

		Input: InputConfig{
			KeyboardEnabled:    true,
			MouseEnabled:       true,
			GamepadEnabled:     true,
			MouseSensitivity:   1.0,
			GamepadSensitivity: 1.0,
			GamepadDeadzone:    0.15,
			KeyMapping:         getDefaultKeyMapping(),
			GamepadMapping:     getDefaultGamepadMapping(),
		},

		Gameplay: GameplayConfig{
			Difficulty:            "normal",
			DamageMultiplier:      1.0,
			EnemyHealthMultiplier: 1.0,
			ExperienceMultiplier:  1.0,
			SoulGainMultiplier:    1.0,
			StaminaRegenRate:      25.0,
			StaminaPenalty:        0.5,
			InvulnerabilityTime:   1.0,
			PerfectBlockWindow:    0.2,
			EnemyRespawnTime:      30.0,
			ItemDespawnTime:       300.0,
			AutoSaveEnabled:       true,
			AutoSaveInterval:      5.0,
		},

		Debug: DebugConfig{
			EnableDebug:      false,
			ShowFPS:          false,
			ShowColliders:    false,
			ShowEntityInfo:   false,
			ShowChunkBorders: false,
			ShowPathfinding:  false,
			EnableGodMode:    false,
			EnableNoclip:     false,
			LogLevel:         "info",
			ConsoleEnabled:   false,
		},

		Paths: PathsConfig{
			AssetsDir:      "assets",
			TexturesDir:    "assets/textures",
			SoundsDir:      "assets/sounds/sfx",
			MusicDir:       "assets/sounds/music",
			MapsDir:        "assets/maps",
			DataDir:        "assets/data",
			SavesDir:       "saves",
			ConfigsDir:     "configs",
			LogsDir:        "logs",
			ScreenshotsDir: "screenshots",
		},
	}
}

// getDefaultKeyMapping retourne le mapping par défaut des touches
func getDefaultKeyMapping() map[string]string {
	return map[string]string{
		"move_up":       "W",
		"move_down":     "S",
		"move_left":     "A",
		"move_right":    "D",
		"attack":        "Space",
		"block":         "Shift",
		"roll":          "LeftControl",
		"interact":      "E",
		"inventory":     "I",
		"map":           "M",
		"pause":         "Escape",
		"quick_slot_1":  "1",
		"quick_slot_2":  "2",
		"quick_slot_3":  "3",
		"quick_slot_4":  "4",
		"cast_spell":    "F",
		"camera_reset":  "R",
		"screenshot":    "F12",
		"debug_console": "BackQuote",
	}
}

// getDefaultGamepadMapping retourne le mapping par défaut de la manette
func getDefaultGamepadMapping() map[string]string {
	return map[string]string{
		"move":         "LeftStick",
		"camera":       "RightStick",
		"attack":       "X",
		"block":        "RightTrigger",
		"roll":         "B",
		"interact":     "A",
		"inventory":    "Y",
		"map":          "Back",
		"pause":        "Start",
		"quick_slot_1": "LeftBumper",
		"quick_slot_2": "RightBumper",
		"target_lock":  "RightStickClick",
	}
}

// ===============================
// DIFFICULTY SETTINGS
// ===============================

// DifficultySettings contient les paramètres par difficulté
type DifficultySettings struct {
	DamageMultiplier      float64
	EnemyHealthMultiplier float64
	EnemySpeedMultiplier  float64
	StaminaRegenRate      float64
	SoulGainMultiplier    float64
	DeathPenalty          float64
}

// GetDifficultySettings retourne les paramètres pour une difficulté donnée
func GetDifficultySettings(difficulty string) DifficultySettings {
	switch difficulty {
	case "easy":
		return DifficultySettings{
			DamageMultiplier:      0.7,
			EnemyHealthMultiplier: 0.8,
			EnemySpeedMultiplier:  0.9,
			StaminaRegenRate:      35.0,
			SoulGainMultiplier:    1.2,
			DeathPenalty:          0.5,
		}
	case "hard":
		return DifficultySettings{
			DamageMultiplier:      1.5,
			EnemyHealthMultiplier: 1.3,
			EnemySpeedMultiplier:  1.1,
			StaminaRegenRate:      15.0,
			SoulGainMultiplier:    0.8,
			DeathPenalty:          1.0,
		}
	case "nightmare":
		return DifficultySettings{
			DamageMultiplier:      2.0,
			EnemyHealthMultiplier: 1.5,
			EnemySpeedMultiplier:  1.2,
			StaminaRegenRate:      10.0,
			SoulGainMultiplier:    0.6,
			DeathPenalty:          1.0,
		}
	default: // normal
		return DifficultySettings{
			DamageMultiplier:      1.0,
			EnemyHealthMultiplier: 1.0,
			EnemySpeedMultiplier:  1.0,
			StaminaRegenRate:      25.0,
			SoulGainMultiplier:    1.0,
			DeathPenalty:          0.8,
		}
	}
}

// ===============================
// CONFIG UTILITIES
// ===============================

// CreateDefaultConfigFile crée un fichier de configuration par défaut
func CreateDefaultConfigFile(configPath string) error {
	config := GetDefaultConfig()
	return config.SaveConfig(configPath)
}

// WindowWidth retourne la largeur configurée
func (c *GameConfig) WindowWidth() int {
	return c.Window.Width
}

// WindowHeight retourne la hauteur configurée
func (c *GameConfig) WindowHeight() int {
	return c.Window.Height
}

// TargetFPS retourne les FPS cibles
func (c *GameConfig) TargetFPS() int {
	return c.Rendering.TargetFPS
}

// TileSize retourne la taille des tiles
func (c *GameConfig) TileSize() int {
	return c.Rendering.TileSize
}

// IsDebugEnabled retourne si le mode debug est activé
func (c *GameConfig) IsDebugEnabled() bool {
	return c.Debug.EnableDebug
}

// GetAudio retourne la configuration audio (pour éviter les cycles d'import)
func (c *GameConfig) GetAudio() AudioConfig {
	return c.Audio
}
