package Hooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"engine/Global"
)

type Hook struct {
	Hook_name string
	Hook_type string
	Executor  any
}

type Hook_obj map[string]*Hook

func CreateHooksObj() *Hook_obj {
	hooksObj := make(Hook_obj)
	return &hooksObj
}

func (hooksObj *Hook_obj) CreateHook(hook_name string, hook_type string) (*Hook, error) {
	if _, exists := (*hooksObj)[hook_name]; exists {
		return nil, fmt.Errorf("hook %q already exists", hook_name)
	}

	hook := &Hook{
		Hook_name: hook_name,
		Hook_type: hook_type,
		Executor:  nil,
	}

	(*hooksObj)[hook_name] = hook

	switch hook.Hook_type {
		case "before_req":

		case "after_req":

		default:
			return nil,fmt.Errorf("unknown hook type %q", hook.Hook_type)
		}

	return hook, nil
}

func (hook *Hook) CreateBeforeReqHookContent() string {
	return fmt.Sprintf(`package hooks

import (
	"net/http"
	"engine/data"
)

func %s(req *http.Request, prevReq *data.ActionResponse, sequences *data.SequnceState) error {
	return nil
}
`, hook.Hook_name)
}

func (hook *Hook) CreateAfterReqHookContent() string {
	return fmt.Sprintf(`package hooks

import (
	"net/http"
	"engine/data"
)

func %s(req *http.Request, response *data.ActionResponse, sequences *data.SequnceState) error {
	return nil
}
`, hook.Hook_name)
}

func (hook *Hook) GenerateHookFile(content string) error {
	if err := os.MkdirAll(Global.GlobalConfig.HooksPath, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	path := filepath.Join(Global.GlobalConfig.HooksPath, hook.Hook_name+".go")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	return nil
}

func (hooksObj *Hook_obj) SetUpHooksObj() error {
	data, err := os.ReadFile(Global.GlobalConfig.HooksPath)
	if err != nil {
		return fmt.Errorf("failed to read hooks config: %w", err)
	}

	if err := json.Unmarshal(data, hooksObj); err != nil {
		return fmt.Errorf("failed to parse hooks config: %w", err)
	}

	return nil
}







func (hooksObj *Hook_obj) ValidateHooks() {
	for hookName, hook := range *hooksObj {
		executor, exists := Hook_registry[hookName]
		if !exists {
			delete(*hooksObj, hookName)
			continue
		}

		switch hook.Hook_type {
		case "before_req":
			if _, ok := executor.(func(*http.Request, *Global.ActionResponse, *Global.SequnceState) error); !ok {
				delete(*hooksObj, hookName)
				continue
			}
		case "after_req":
			if _, ok := executor.(func(*http.Request, *Global.ActionResponse, *Global.SequnceState) error); !ok {
				delete(*hooksObj, hookName)
				continue
			}
		default:
			delete(*hooksObj, hookName)
			continue
		}

		hook.Executor = executor
	}
}