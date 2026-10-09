package Actions

import (
	"fmt"
	"strconv"
	"strings"
)









func (this *Http_Action_obj)CreateAction(global_name string,action_name string, action_opts Action_opts) Http_action {
	httpOpts, ok := action_opts.Opts.(HTTP_opts)
	if !ok {
		fmt.Println("Invalid HTTP options")
		return Http_action{}
	}

	var http_action = Http_action{
		Global_namme: 			global_name,
		Name:                   action_name,
		Executor_type:          httpOpts.Executor_type,
		Path:                   "",
		Executor_function_name: "",
		Opts: HTTP_opts{},
	}

	return http_action

}




func BytesToGoLiteral(body []byte) string {
	values := make([]string, len(body))
	for i, b := range body {
		values[i] = strconv.Itoa(int(b))
	}

	return "[]byte{" + strings.Join(values, ", ") + "}"
}


