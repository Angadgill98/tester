package Actions

import (
	"encoding/json"
	"engine/Global"
	"fmt"
	"os"

	"github.com/google/uuid"
)


type WS_Actions_config map[uuid.UUID]*WsAction

type WsAction struct {
	Global_name            string  `json:"global_name"`
	Name                   string  `json:"name"`
	Path                   string  `json:"path"`
	Executor_type          string  `json:"executor_type"`
	Executor_function_name string  `json:"executor_function_name"`
	Opts                   WS_opts `json:"opts"`
}

type WS_opts struct {
	Executor_type string `json:"executor_type"`
}

func CreateWsAction(action_name string, opts WS_opts)* WsAction {
	return &WsAction{
		Global_name: action_name,
		Name:        action_name,
		Opts:        opts,
	}
}

func CreateWSActionsConfig() WS_Actions_config {
	return make(WS_Actions_config)
}
var ws_actions_config_path = Global.GlobalConfig.Ws_ActionsConfigPath

func (config *WS_Actions_config) RestoreWSActionsConfig() error {
	data, err := os.ReadFile(ws_actions_config_path)
	if err != nil {
		return fmt.Errorf("failed to read WebSocket actions config: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse WebSocket actions config: %w", err)
	}

	return nil
}

func (config *WS_Actions_config) SaveWSActionsConfig() error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("Engine: failed to marshal WebSocket actions config: %w", err)
	}

	if err := os.MkdirAll("Actions", 0755); err != nil {
		return fmt.Errorf("Engine: failed to create actions config directory: %w", err)
	}

	if err := os.WriteFile(ws_actions_config_path, data, 0644); err != nil {
		return fmt.Errorf("Engine: failed to save WebSocket actions config: %w", err)
	}

	return nil
}