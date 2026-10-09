package Actions

import (
	"encoding/json"
	"engine/Global"
	"fmt"
	"os"
)

type Actions_obj struct{
	Http Http_Action_obj
	Ws Ws_Action_obj
}


func CreateActionsObj() (*Actions_obj, error) {
	actionsObj := &Actions_obj{}

	actionsObj.Http = *CreateHttpActionObj()
	if err := ReadHTTPActionsConfig( &actionsObj.Http); err != nil {
		return nil, err
	}

	actionsObj.Ws = *CreateWsActionObj()
	if err := ReadWSActionsConfig( &actionsObj.Ws); err != nil {
		return nil, err
	}

	return actionsObj, nil
}




type Http_Action_obj map[string][]Http_action

func CreateHttpActionObj() *Http_Action_obj {
	actionObj := make(Http_Action_obj)
	return &actionObj
}

func ReadHTTPActionsConfig(actionObj *Http_Action_obj) error {
	data, err := os.ReadFile(Global.GlobalConfig.Http_ActionsConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read HTTP actions file: %w", err)
	}

	if err := json.Unmarshal(data, actionObj); err != nil {
		return fmt.Errorf("failed to parse HTTP actions file: %w", err)
	}

	return nil
}




type Ws_Action_obj map[string][]WsAction

func CreateWsActionObj() *Ws_Action_obj {
	actionObj := make(Ws_Action_obj)
	return &actionObj
}

func ReadWSActionsConfig( actionObj *Ws_Action_obj) error {
	data, err := os.ReadFile(Global.GlobalConfig.Ws_ActionsConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read WebSocket actions file: %w", err)
	}

	if err := json.Unmarshal(data, actionObj); err != nil {
		return fmt.Errorf("failed to parse WebSocket actions file: %w", err)
	}

	return nil
}


