package Actions

import (
	"encoding/json"
	"engine/Global"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Actions_opts struct {
	Action_type string    `json:"action_type"`
	HTTP        HTTP_opts `json:"http"`
	WS          WS_opts   `json:"ws"`
}

type Http_Actions_config map[uuid.UUID]*Http_action

type Http_action struct {
	Global_namme           string    `json:"global_name"`
	Name                   string    `json:"name"`
	Path                   string    `json:"path"`
	Executor_type          string    `json:"executor_type"`
	Executor_function_name string    `json:"executor_function_name"`
	Opts                   HTTP_opts `json:"opts"`
}

type HTTP_opts struct {
	Executor_type string            `json:"executor_type"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Query         []string          `json:"query"`
	Body          []byte            `json:"body"`
}

func CreateHttpAction(action_name string, opts HTTP_opts) *Http_action {
	return &Http_action{
		Global_namme: action_name,
		Name:         action_name,
		Opts:         opts,
	}
}

func CreateHttpActionsConfig() Http_Actions_config {
	return make(Http_Actions_config)
}


var http_actions_config_path = Global.GlobalConfig.Http_ActionsConfigPath

func (config *Http_Actions_config) RestoreHttpActionsConfig() error {
	data, err := os.ReadFile(http_actions_config_path)
	if err != nil {
		return fmt.Errorf("failed to read HTTP actions config: %w", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse HTTP actions config: %w", err)
	}

	return nil
}


func (config *Http_Actions_config) SaveHttpActionsConfig() error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("Engine: failed to marshal HTTP actions config: %w", err)
	}

	if err := os.MkdirAll("Actions", 0755); err != nil {
		return fmt.Errorf("Engine: failed to create actions config directory: %w", err)
	}

	if err := os.WriteFile(http_actions_config_path, data, 0644); err != nil {
		return fmt.Errorf("Engine: failed to save HTTP actions config: %w", err)
	}

	return nil
}


