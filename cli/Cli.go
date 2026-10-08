package main

import (
	"net"
	"time"
)


type CLI struct {
	unix_socket *net.Conn
}

func CreateCLI(socket *net.Conn ) *CLI {
	return &CLI{
		unix_socket: socket,
	}
}

func ConnectUnixSocketWithRetries() net.Conn {
	for i := 0; i < global.SocketRetries; i++ {
		conn, err := net.Dial("unix", "/tmp/testflow.sock")
		if err == nil {
			return conn
		}

		time.Sleep(time.Second)
	}

	return nil
}
