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
	
	ApplicationObj *Application.Application_obj
	ActionsObj *Actions.Actions_obj
	SequencesObj *Sequences.SequenceObj
}


func CreateEngine(applicationObj *Application.Application_obj,actions_obj *Actions.Actions_obj, sequences *Sequences.SequenceObj) Engine {
	return Engine{
		ApplicationObj:   applicationObj,
		ActionsObj: actions_obj,
		SequencesObj:    sequences,
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

	application := engine.ApplicationObj.CreateApplication(appName)
	currentApp.Children[appName] = application
}

func (engine *Engine) GetAppFromPath(path string) *Application.Application {
	currentApp := (*engine.ApplicationObj)["GLOBAL"]

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

		action := engine.ActionsObj.Http.CreateAction(global_actionName, actionName, opts)
		action.Opts = httpOpts
		action.Path = actions_path

	case "ws":
		wsOpts, ok := opts.Opts.(Actions.WS_opts)
		if !ok {
			return
		}

		action := engine.ActionsObj.Ws.CreateAction(global_actionName, actionName, opts)
		action.Opts = wsOpts
		action.Path = actions_path
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