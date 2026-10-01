package app

import (
	actions "tester/Actions"
	sequences "tester/Sequences"
)

type App struct {
	Name      string
	Actions   map[string]actions.Action
	Sequences map[string]sequences.Sequence
}

type Apps struct {
	Apps map[string]App
}

func CreateAppObj(name string) App {
	return App{
		Name:      name,
		Actions:   make(map[string]actions.Action),
		Sequences: make(map[string]sequences.Sequence),
	}
}

func CreateAppsObj() Apps {
	return Apps{
		Apps: make(map[string]App),
	}
}

func (apps *Apps) AddApp(app App) bool {
	if _, exists := apps.Apps[app.Name]; exists {
		return false
	}

	apps.Apps[app.Name] = app

	return true
}