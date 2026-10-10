package main

import (
	"engine/Actions"
	"engine/Application"
	"engine/Sequences"
	"fmt"

	"github.com/google/uuid"
)



type Engine struct {
	App_config      *Application.Application_config
	Http_config     *Actions.Http_Actions_config
	Ws_config       *Actions.WS_Actions_config
	Sequence_config *Sequences.Sequence_config
}

func CreateEngine(appConfig *Application.Application_config, httpConfig *Actions.Http_Actions_config, wsConfig *Actions.WS_Actions_config, sequenceConfig *Sequences.Sequence_config) Engine {
	return Engine{
		App_config:      appConfig,
		Http_config:     httpConfig,
		Ws_config:       wsConfig,
		Sequence_config: sequenceConfig,
	}
}


func (engine *Engine) CreateApp(appName string, parent_appid uuid.UUID) error {
	parentApp, exists := (*engine.App_config)[parent_appid]
	if !exists {
		return fmt.Errorf("Engine: Parent application with ID %s does not exist", parent_appid)
	}

	if _, exists := parentApp.Children[appName]; exists {
		return fmt.Errorf("Engine: Application %q already exists under parent %q", appName, parentApp.Name)
	}

	app := Application.CreateApplication(appName,parent_appid)

	(*engine.App_config)[app.Application_id] = app
	parentApp.Children[appName] = app.Application_id

	if err := engine.App_config.SaveAppConfig(); err != nil {
		return fmt.Errorf("Engine: Failed to save application config: %w", err)
	}

	return nil
}


func (engine *Engine) CreateAction(appid uuid.UUID, actionName string, opts any) error {
	app, exists := (*engine.App_config)[appid]
	if !exists {
		return fmt.Errorf("Engine: Application with ID %s does not exist", appid)
	}

	if _, exists := app.Actions[actionName]; exists {
		return fmt.Errorf("Engine: Action %q already exists in application %q", actionName, app.Name)
	}

	switch actionOpts := opts.(type) {
	case Actions.HTTP_opts:
		actionID := uuid.New()
		_ = actionOpts

		(*engine.Http_config)[actionID] = Actions.CreateHttpAction(actionName, actionOpts)
		app.Actions[actionName] = actionID

		if err := engine.Http_config.SaveHttpActionsConfig(); err != nil {
			return fmt.Errorf("Engine: Failed to save HTTP actions config: %w", err)
		}

	case Actions.WS_opts:
		actionID := uuid.New()

		(*engine.Ws_config)[actionID] = Actions.CreateWsAction(actionName, actionOpts)
		app.Actions[actionName] = actionID

		if err := engine.Ws_config.SaveWSActionsConfig(); err != nil {
			return fmt.Errorf("Engine: Failed to save WebSocket actions config: %w", err)
		}

	default:
		return fmt.Errorf("Engine: Unsupported action options type %T; expected Actions.HTTP_opts or Actions.WS_opts", opts)
	}

	if err := engine.App_config.SaveAppConfig(); err != nil {
		return fmt.Errorf("Engine: Failed to save application config: %w", err)
	}

	return nil
}




