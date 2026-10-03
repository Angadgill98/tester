package sequences

import (
	actions "tester/Actions"
)

var Sequence_Config = map[string][]SequenceAction{
	"asd": {
		{
			Action_type: actions.ActionType("http"),
			HTTP_action: &actions.Http_actions{
				Name: "asd",
				Path: "data/actions/http/req",
				Executor_type: "req",
				Executor_function_name: "StartAction_asd",
			},
		},
		{
			Action_type: actions.ActionType("ws"),
			WS_action: &actions.WsAction{
				Name: "asd",
				Path: "data/actions/ws/send",
				Exector_type: "send",
				Executor_function_name: "StartAction_asd",
			},
		},
	},
}
