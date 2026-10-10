package Application

import (
	"encoding/json"
	"engine/Global"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var app_config = Global.GlobalConfig.ApplicationConfigPath

type Application_config map[uuid.UUID]*Application

func CreateApplicationConfig() *Application_config {
	config := make(Application_config)
	return &config
}

type Application struct {
	Application_id uuid.UUID            `json:"application_id"`
	Parent_id      uuid.UUID            `json:"parent_id"`
	Name           string               `json:"name"`
	Actions        map[string]uuid.UUID `json:"actions"`
	Sequences      map[string]uuid.UUID `json:"sequences"`
	Children       map[string]uuid.UUID `json:"children"`
}

type Application_State struct {
}

func CreateApplication(name string, parentID uuid.UUID) *Application {
	return &Application{
		Application_id: uuid.New(),
		Parent_id:      parentID,
		Name:           name,
		Actions:        make(map[string]uuid.UUID),
		Sequences:      make(map[string]uuid.UUID),
		Children:       make(map[string]uuid.UUID),
	}
}

func (config *Application_config) RestoreAppConfig() error {
	data, err := os.ReadFile(app_config)
	if err != nil {
		return fmt.Errorf("Engine: failed to read application config: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("Engine: failed to parse application config: %w", err)
	}

	if *config == nil {
		*config = make(Application_config)
	}

	return nil
}

func (config *Application_config) SaveAppConfig() error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("Engine: failed to marshal application config: %w", err)
	}

	if err := os.MkdirAll("Application", 0755); err != nil {
		return fmt.Errorf("Engine: failed to create application config directory: %w", err)
	}

	if err := os.WriteFile(app_config, data, 0644); err != nil {
		return fmt.Errorf("Engine: failed to save application config: %w", err)
	}

	return nil
}