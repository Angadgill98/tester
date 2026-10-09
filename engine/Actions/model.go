package Actions


type Executor_function_signature any


type Action_opts struct {
	Action_type string `json:"action_type"`
	Opts        any    `json:"opts"`
}

type Http_action struct {
	Global_namme           string    `json:"global_name"`
	Name                   string    `json:"name"`
	Path                   string    `json:"path"`
	Executor_type          string    `json:"executor_type"`
	Executor_function_name string    `json:"executor_function_name"`
	Opts                   HTTP_opts `json:"opts"`
}

type HTTP_opts struct {
	Executor_type string            `json:"executor_type"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Query         []string          `json:"query"`
	Body          []byte            `json:"body"`
}

type WsAction struct {
	Global_name            string  `json:"global_name"`
	Name                   string  `json:"name"`
	Path                   string  `json:"path"`
	Executor_type          string  `json:"executor_type"`
	Executor_function_name string  `json:"executor_function_name"`
	Opts                   WS_opts `json:"opts"`
}

type WS_opts struct {
	Executor_type string `json:"executor_type"`
}