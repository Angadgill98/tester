package Actions

import (
	http_req "tester/data/actions/http/req"
	ws_actions_send "tester/data/actions/ws/send"
	ws_actions_connect "tester/data/actions/ws/connect"
	ws_actions_receive "tester/data/actions/ws/recieve"
)

//doesnt allow me to define these apcges wihtout the mock 

func ignore(){
http_req.MockFunction()

ws_actions_send.MockFunction()

ws_actions_connect.MockFunction()

ws_actions_receive.MockFunction()
}
			

var ActionsConfig = map[string]map[string]Executor_function_signature{
	"http": {
	"asd": http_req.StartAction_asd,
	},
	"ws": {
	"asd": ws_actions_send.WS_send_asd,
	},
}