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


func (engine *Engine) CreateApp(appName string, path string)string {
	currentApp := engine.GetAppFromPath(path)
	if currentApp == nil {
		return fmt.Sprintf("Engine: Invalid app path")
	}
	fmt.Println("Engine: App path is valid")

	if _, exists := currentApp.Children[appName]; exists {
		return fmt.Sprintf("Engine: App already exist")
	}

	application := engine.ApplicationObj.CreateApplication(appName)
	currentApp.Children[appName] = application

	if err := engine.ApplicationObj.SaveApplicationObj(); err != nil {
		return fmt.Sprintf("Engine: Failed to save application config: %v\n", err)
	}
	return fmt.Sprintf("Engine: Created and Saved the app")

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

func (engine *Engine) CreateAction(appName string, path string, actionName string, opts Actions.Action_opts) error {
	currentApp := engine.GetAppFromPath(path)
	if currentApp == nil {
		return fmt.Errorf("application path %q does not exist", path)
	}

	if _, exists := currentApp.Children[appName]; !exists {
		return fmt.Errorf("application %q does not exist at path %q", appName, path)
	}

	actionsPath := engine.GetActionPath(path, appName)
	engine.CreateActionPath(actionsPath)

	globalActionName := path + "_" + appName + "_" + actionName

	switch opts.Action_type {
	case "http":
		httpOpts, ok := opts.Opts.(Actions.HTTP_opts)
		if !ok {
			return fmt.Errorf("invalid HTTP action options")
		}

		action := engine.ActionsObj.Http.CreateAction(globalActionName, actionName, opts)
		action.Opts = httpOpts
		action.Path = actionsPath

		engine.ActionsObj.Http[globalActionName] = append(engine.ActionsObj.Http[globalActionName], action)

	case "ws":
		wsOpts, ok := opts.Opts.(Actions.WS_opts)
		if !ok {
			return fmt.Errorf("invalid WebSocket action options")
		}

		action := engine.ActionsObj.Ws.CreateAction(globalActionName, actionName, opts)
		action.Opts = wsOpts
		action.Path = actionsPath

		engine.ActionsObj.Ws[globalActionName] = append(engine.ActionsObj.Ws[globalActionName], action)

	default:
		return fmt.Errorf("unsupported action type %q", opts.Action_type)
	}

	if err := engine.ActionsObj.SaveActionsObj(); err != nil {
		return fmt.Errorf("failed to save actions: %w", err)
	}

	return nil
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


