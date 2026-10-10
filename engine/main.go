package main

import (
	"engine/Actions"
	"engine/Application"
	"engine/Sequences"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	app_config := Application.CreateApplicationConfig()
	err :=app_config.RestoreAppConfig()
	if err != nil {
		fmt.Printf("Engine: Failed to create application config: %v\n", err)
		os.Exit(1)
	}

	Globall_Init(app_config)

	if err := app_config.RestoreAppConfig(); err != nil {
		fmt.Printf("Engine: Failed to restore application config: %v\n", err)
		os.Exit(1)
	}

	http_config := Actions.CreateHttpActionsConfig()
	if err := http_config.RestoreHttpActionsConfig(); err != nil {
		fmt.Printf("Engine: Failed to restore HTTP actions config: %v\n", err)
		os.Exit(1)
	}

	ws_config := Actions.CreateWSActionsConfig()
	if err := ws_config.RestoreWSActionsConfig(); err != nil {
		fmt.Printf("Engine: Failed to restore WebSocket actions config: %v\n", err)
		os.Exit(1)
	}

	sequences_config := Sequences.CreateSequenceConfig()
	err =sequences_config.SetUpSequenceConfig()
	if err != nil {
		fmt.Printf("Engine: Failed to create or restore sequence config: %v\n", err)
		os.Exit(1)
	}

	engine:=CreateEngine(app_config, &http_config, &ws_config, sequences_config)

	StartSocket(&engine)
}

func Globall_Init(app_config *Application.Application_config){
	globalExists := false
	for _, app := range *app_config {
		if app.Name == "GLOBAL" {
			globalExists = true
			break
		}
	}

	if !globalExists {
		globalApp := Application.CreateApplication("GLOBAL",uuid.Nil)
		(*app_config)[globalApp.Application_id] = globalApp

		if err := app_config.SaveAppConfig(); err != nil {
			fmt.Printf("Engine: Failed to save GLOBAL application: %v\n", err)
			os.Exit(1)
		}
	}
}