package Actions

import (
	"fmt"
)







var Ws_executor_types = []string{
	"send",
	"connect",
	"recieve",
}



func (this *Ws_Action_obj)CreateAction(global_name string,action_name string, action_opts Action_opts) WsAction {
	httpOpts, ok := action_opts.Opts.(HTTP_opts)
	if !ok {
		fmt.Println("Invalid HTTP options")
		return nil
	}

	var ws_action = WsAction{
		Global_name: 			global_name,
		Name:                   action_name,
		Executor_type:          httpOpts.Executor_type,
		Path:                   "",
		Executor_function_name: "",
		Opts: WS_opts{},
	}



	return ws_action

}