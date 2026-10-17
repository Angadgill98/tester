package main

import "encoding/json"


type Actions_opts struct {
	Action_type string    `json:"action_type"`
	HTTP        HTTP_opts `json:"http"`
	WS          WS_opts   `json:"ws"`
}


type HTTP_opts struct {
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Query         []string          `json:"query"`
	Body          json.RawMessage   `json:"body"`
}

type WS_opts struct {
	Executor_type string `json:"executor_type"`
}


type Application struct {
	AppId       int64
	AppName     string
	ParentAppId int64
	Actions     []int64
	Sequences   []int64
	Children    []int64
}


type Response struct {
	Status bool  `json:"status"`
	Msg    string `json:"msg"`
	Args   []any  `json:"args"`
}

