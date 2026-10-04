package sequences

import (
	actions "tester/Actions"
)

var Sequence_Config = map[string][]SequenceAction{
	"asd": {
		{
			Action_type: actions.ActionType("ws"),
			WS_action: &actions.WsAction{
				Name: "StartAction_rev",
				Path: "data/actions/ws/recieve/",
				Exector_type: "recieve",
				Executor_function_name: "StartAction_rev",
			},
		},
		{
			Action_type: actions.ActionType("ws"),
			WS_action: &actions.WsAction{
				Name: "StartAction_send",
				Path: "data/actions/ws/send/",
				Exector_type: "send",
				Executor_function_name: "StartAction_send",
			},
		},
		{
			Action_type: actions.ActionType("ws"),
			WS_action: &actions.WsAction{
				Name: "StartAction_con2",
				Path: "data/actions/ws/connect/",
				Exector_type: "connect",
				Executor_function_name: "StartAction_con2",
			},
		},
	},
}
