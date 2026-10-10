package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"engine/Global"
)

func StartSocket(engine *Engine) {
	socketPath := Global.GlobalConfig.SocketPath

	if err := os.MkdirAll(filepath.Dir(socketPath), 0755); err != nil {
		fmt.Printf("Failed to create socket directory: %v\n", err)
		return
	}

	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to remove existing socket: %v\n", err)
		return
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("Failed to bind socket: %v\n", err)
		return
	}
	defer listener.Close()
	defer os.Remove(socketPath)

	fmt.Printf("Listening on socket: %s\n", socketPath)

	if err := ListenForClients(listener,engine); err != nil {
		fmt.Printf("Socket listener stopped: %v\n", err)
	}


}

func ListenForClients(listener net.Listener,engine *Engine) error {
	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept client: %w", err)
	}

	HandleClient(conn,engine)
	return nil
}

type Command struct {
	Command string `json:"command"`
	Args    []any  `json:"args"`
}

func HandleClient(conn net.Conn, engine *Engine) {
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Failed to read from client: %v\n", err)
			}
			break
		}

		if n > 0 {
			var command Command

			if err := json.Unmarshal(buffer[:n], &command); err != nil {
				fmt.Printf("Failed to parse command JSON: %v\n", err)
				continue
			}

			fmt.Printf("Command: %s\n", command.Command)
			fmt.Printf("Args: %v\n", command.Args)
			HandleCommand(command, engine)
			// Process command using engine here.
		}
	}
}


func HandleCommand(command Command, engine *Engine) {
	switch command.Command {
	case "create":
		ProcessCreate(command.Args, engine)
	case "cd":
		ProcessCd(command.Args,engine)
	default:
		fmt.Println("Engine: Unknown command:", command.Command)
	}
}

func ProcessCreate(args []any, engine *Engine) {
	if len(args) == 0 {
		fmt.Println("Engine: Missing create operation")
		return
	}

	operation, ok := args[0].(string)
	if !ok {
		fmt.Println("Engine: Create operation must be a string")
		return
	}

	switch operation {
	case "app":
		if len(args) < 2 {
			fmt.Println("Engine: Missing app name")
			return
		}

		appName, ok := args[1].(string)
		if !ok {
			fmt.Println("Engine: App name must be a string")
			return
		}

		if len(args) < 3 {
			fmt.Println("Engine: Missing app path")
			return
		}

		appPath, ok := args[2].(string)
		if !ok {
			fmt.Println("Engine: App path must be a string")
			return
		}

		fmt.Printf("Engine: Creating app:%v with the path %v\n", appName,appPath)
		engine.CreateApp(appName, appPath)
	}
}

func ProcessCd(args []any,engine *Engine ){
	if len(args) == 0 {
		fmt.Println("Engine: Missing CD operation")
		return
	}

	operation, ok := args[0].(string)
	if !ok {
		fmt.Println("Engine: CD operation must be a string")
		return
	}

	switch operation{
		case "app":
			if len(args) < 2 {
				fmt.Println("Engine: Missing application name")
				return
			}

			appName, ok := args[1].(string)
			if !ok {
				fmt.Println("Engine: Application name must be a string")
				return
			}

			
			if len(args) < 3 {
				fmt.Println("Engine: Missing app path")
				return
			}

			appPath, ok := args[2].(string)
			if !ok {
				fmt.Println("Engine: App path must be a string")
				return
			}


			app := engine.GetAppFromPath(appPath + "_" + appName)
			if app == nil {
				fmt.Printf("Engine: Application %q does not exist at path %q\n", appName, appPath)
				return
			}

		default:
	}
}


func ProcessLs(args []any,engine *Engine ){
	if len(args) == 0 {
		
		return
	}

	
}



