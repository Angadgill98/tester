package main

import (
	"engine/Actions"
	"engine/Application"
	"engine/Global"
	"engine/Sequences"
	"os"
	"path/filepath"
	"strings"
)




type Engine struct {
	
	Application *Application.Application
	Http_actions *Actions.HttpActionsConfig
	Ws_action *Actions.WsActionsConfig
	Sequences *Sequences.SequenceConfig
}


func CreateEngine(application *Application.Application, httpActions *Actions.HttpActionsConfig, wsActions *Actions.WsActionsConfig, sequences *Sequences.SequenceConfig) Engine {
	return Engine{
		Application:   application,
		Http_actions: httpActions,
		Ws_action:    wsActions,
		Sequences:    sequences,
	}
}


func (engine *Engine) CreateApp(appName string, path string) {
	currentApp := engine.GetAppFromPath(path)
	if currentApp == nil {
		return
	}

	if _, exists := currentApp.Children[appName]; exists {
		return
	}

	application := Application.CreateApplication(appName)
	currentApp.Children[appName] = application
}

func (engine *Engine) GetAppFromPath(path string) *Application.Application {
	currentApp := engine.Application

	for _, name := range engine.GetAppNamesFromPath(path) {
		child, exists := currentApp.Children[name]
		if !exists {
			return nil
		}

		currentApp = child
	}

	return currentApp
}
func (engine *Engine) GetAppNamesFromPath(path string) []string {
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "_")
}

func (engine *Engine) CreateAction(appName string, path string, actionName string,opts Actions.Action_opts) {
	currentApp := engine.GetAppFromPath(path)
	if currentApp == nil {
		return
	}

	if _, exists := currentApp.Children[appName]; !exists {
		return
	}

	var actions_path=engine.GetActionPath(path,appName)
	engine.CreateActionPath(actions_path)


	var global_actionName=path+"_"+appName+"_"+actionName


	switch opts.Action_type {
	case "http":
		httpOpts, ok := opts.Opts.(Actions.HTTP_opts)
		if !ok {
			return
		}

		action := engine.Http_actions.CreateAction(global_actionName, actionName, opts)
		action.Opts = httpOpts
		action.Path = actions_path
		action.CreateHttpFile()

	case "ws":
		wsOpts, ok := opts.Opts.(Actions.WS_opts)
		if !ok {
			return
		}

		action := engine.Ws_action.CreateAction(global_actionName, actionName, opts)
		action.Opts = wsOpts
		action.Path = actions_path
		action.CreateWsFile()
	}


	
	
}


func (engine *Engine) GetActionPath(path string, appName string) string {
	actionPath := Global.GlobalConfig.ActionStoragePath+ path + "_" + appName
	actionPath = strings.ReplaceAll(actionPath, "_", "/")

	return actionPath
}

func (engine *Engine) CreateActionPath(actionPath string) error {

	return os.MkdirAll(filepath.Dir(actionPath), 0755)
}