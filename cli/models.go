package main


type Actions_opts struct {
	Action_type string    `json:"action_type"`
	HTTP        HTTP_opts `json:"http"`
	WS          WS_opts   `json:"ws"`
}


type HTTP_opts struct {
	Executor_type string            `json:"executor_type"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Query         []string          `json:"query"`
	Body          []byte            `json:"body"`
}

type WS_opts struct {
	Executor_type string `json:"executor_type"`
}