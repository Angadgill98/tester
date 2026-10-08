package main




type Global struct {
	SocketRetries int
}

var global = Global{
	SocketRetries: 3,
}