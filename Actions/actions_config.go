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
	"webhook_update_url": http_req.StartAction_webhook_update_url,
	"get_webhooks": http_req.StartAction_get_webhooks,
	"create_webhook": http_req.StartAction_create_webhook,
	"signin": http_req.StartAction_signin,
	"webhook_event": http_req.StartAction_webhook_event,
	"delete_webhook": http_req.StartAction_delete_webhook,
	"update_webhook_name": http_req.StartAction_update_webhook_name,
	},
	"ws": {

	},
}