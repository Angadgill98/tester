package main

type Global struct {
	SocketRetries  int
	EnginePath     string
	Editor         string
	CurrentAppID   string
	CurrentAppPath []string
}

var global = Global{
	SocketRetries:  3,
	EnginePath:     "../engine",
	Editor:         "code",
	CurrentAppID:   "",
	CurrentAppPath: []string{},
}