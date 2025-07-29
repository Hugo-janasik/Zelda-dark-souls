// internal/audio/audio_manager.go - Gestionnaire audio
package audio

// AudioConfig configuration audio (copié de core pour éviter le cycle)
type AudioConfig struct {
	MasterVolume float64
	MusicVolume  float64
	SFXVolume    float64
	VoiceVolume  float64
	EnableAudio  bool
	SampleRate   int
	BufferSize   int
	MaxSounds    int
}

// GameConfig interface minimale pour éviter le cycle d'import
type GameConfig interface {
	GetAudio() AudioConfig
}

type AudioManager struct {
	config *AudioConfig
}

func NewAudioManager(config GameConfig) (*AudioManager, error) {
	audioConfig := config.GetAudio()
	return &AudioManager{config: &audioConfig}, nil
}

func (am *AudioManager) UpdateConfig(config *AudioConfig) {
	am.config = config
}

func (am *AudioManager) PauseAll() {
	// TODO: Implémenter pause audio
}

func (am *AudioManager) ResumeAll() {
	// TODO: Implémenter resume audio
}

func (am *AudioManager) Cleanup() {
	// TODO: Nettoyer les ressources audio
}
