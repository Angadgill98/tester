package main

type Global struct {
	SocketRetries int
	EnginePath    string
	CurrentAppID	string
	CurrentAppPath	[]string
}

var global = Global{
	SocketRetries: 3,
	EnginePath:    "../engine",
	CurrentAppID: "",
	CurrentAppPath:	[]string{},

}