package main

type Global struct {
	SocketRetries  int
	EnginePath     string
	Editor         string
	CurrentAppID   int64
}

var global = Global{
	SocketRetries: 3,
	EnginePath:    "../engine",
	Editor:        "code",
	CurrentAppID:  0,
}