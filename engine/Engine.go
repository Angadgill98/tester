package main

import (
	"engine/Actions"
	"engine/Application"
	"engine/Global"
	"engine/Sequences"
	"fmt"
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
		fmt.Println("Engine: Invalid app path")
		return
	}
	fmt.Println("Engine: App path is valid")

	if _, exists := currentApp.Children[appName]; exists {
		fmt.Println("Engine: App already exist")

		return
	}

	application := engine.ApplicationObj.CreateApplication(appName)
	currentApp.Children[appName] = application

	if err := engine.ApplicationObj.SaveApplicationObj(); err != nil {
		fmt.Printf("Engine: Failed to save application config: %v\n", err)
		return
	}
	fmt.Println("Engine: Created and Saved the app")

}

func (engine *Engine) GetAppFromPath(path string) *Application.Application {
	currentApp := (*engine.ApplicationObj)["GLOBAL"]

	for _, name := range engine.GetAppNamesFromPath(path) {
		child, exists := currentApp.Children[name]
		if !exists {
			child, exists = (*engine.ApplicationObj)[name]
		}

		if !exists {
			fmt.Printf("Engine: Failed to get the app:%v on the path:%v\n", name, path)
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

		engine.ActionsObj.Http[global_actionName] = append(engine.ActionsObj.Http[global_actionName], action)

	case "ws":
		wsOpts, ok := opts.Opts.(Actions.WS_opts)
		if !ok {
			return
		}

		action := engine.ActionsObj.Ws.CreateAction(global_actionName, actionName, opts)
		action.Opts = wsOpts
		action.Path = actions_path
	}


	if err := engine.ActionsObj.SaveActionsObj(); err != nil {
		fmt.Printf("Engine: Failed to save actions: %v\n", err)
		return
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

func (engine *Engine) GetAppChildren(path string) ([]string, error) {
	app := engine.GetAppFromPath(path)
	if app == nil {
		return nil, fmt.Errorf("application path %q does not exist", path)
	}

	children := make([]string, 0, len(app.Children))
	for name := range app.Children {
		children = append(children, name)
	}

	return children, nil
}


