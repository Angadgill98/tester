package actions

import (
	"net/http"
	"tester/data"
	data_actions "tester/data/actions"
)

func asd(){
	if data_actions.MockExecutor{}
}
			

var ActionsConfig = map[string]func(*http.Client, *data.ActionResponse, *data.SequnceState) *data.ActionResponse{
	
	"as": data_actions.StartAction_as,
}