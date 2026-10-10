package main

type Global struct {
	SocketRetries int
	EnginePath    string
	GlobalAppName string
	AppPath       []string
}

var global = Global{
	SocketRetries: 3,
	EnginePath:    "../engine",
	GlobalAppName: "GLOBAL",
	AppPath:       []string{"GLOBAL"},
}