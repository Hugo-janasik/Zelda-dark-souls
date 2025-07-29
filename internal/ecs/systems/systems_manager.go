// internal/ecs/systems/systems_manager.go - Gestionnaire ECS
package systems

import "time"

type SystemsManager struct{}

func NewSystemsManager() *SystemsManager {
	return &SystemsManager{}
}

func (sm *SystemsManager) UpdateAll(entities []interface{}, deltaTime time.Duration) {
	// TODO: Implémenter la mise à jour des systèmes ECS
}

func (sm *SystemsManager) Cleanup() {
	// TODO: Nettoyer les systèmes
}
