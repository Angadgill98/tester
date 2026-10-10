package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"engine/Actions"
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
			status,response:=HandleCommand(command, engine)
			// Process command using engine here.
			responseObj := Response{
				Status:   status,
				Response: response,
			}

			if err := SendResponse(conn, responseObj); err != nil {
				fmt.Printf("Failed to send response: %v\n", err)
			}
		}
	}
}


func HandleCommand(command Command, engine *Engine) (bool,any) {
	switch command.Command {
	case "create":
		status,response:=ProcessCreate(command.Args, engine)
		return status,response
	case "cd":
		status,response:=ProcessCd(command.Args,engine)
		return status,response
	case "ls":
		status,response:=ProcessLs(command.Args,engine)
		return status,response
	default:
		response:=fmt.Sprintf("CLI: Unknown command: %v", command.Command)
		return false,response

	}
}

type Response struct {
	Status   bool `json:"status"`
	Response any    `json:"response"`
}

func SendResponse(conn net.Conn, responseObj Response) error {
	data, err := json.Marshal(responseObj)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}


func ProcessCreate(args []any, engine *Engine) (bool, any) {
	if len(args) == 0 {
		return false, "Engine: Missing create operation"
	}

	operation, ok := args[0].(string)
	if !ok {
		return false, "Engine: Create operation must be a string"
	}

	switch operation {
	case "app":
		if len(args) < 2 {
			return false, "Engine: Missing app name"
		}

		appName, ok := args[1].(string)
		if !ok {
			return false, "Engine: App name must be a string"
		}

		if len(args) < 3 {
			return false, "Engine: Missing app path"
		}

		appPath, ok := args[2].(string)
		if !ok {
			return false, "Engine: App path must be a string"
		}

		fmt.Printf("Engine: Creating app: %v with the path %v\n", appName, appPath)

		response := engine.CreateApp(appName, appPath)
		if strings.HasPrefix(response, "Engine: ") {
			return !strings.Contains(response, "Failed"), response
		}

		return true, response
	case "action":
		if len(args) < 2 {
			return false, "Engine: Missing action name"
		}

		actionName, ok := args[1].(string)
		if !ok {
			return false, "Engine: Action name must be a string"
		}

		if len(args) < 3 {
			return false, "Engine: Missing app path"
		}

		appPath, ok := args[2].(string)
		if !ok {
			return false, "Engine: App path must be a string"
		}

		apps := engine.GetAppNamesFromPath(appPath)
		if len(apps) == 0 {
			return false, fmt.Sprintf("Engine: Application path %q does not exist", appPath)
		}

		if len(args) < 4 {
			return false, "Engine: Missing action options"
		}

		opts, ok := args[len(args)-1].(Actions.Action_opts)
		if !ok {
			return false, "Engine: Action options must be of type Actions.Action_opts"
		}

		appName := apps[len(apps)-1]
		parentPath := strings.TrimSuffix(appPath, "_"+appName)

		if err := engine.CreateAction(appName, parentPath, actionName, opts); err != nil {
			return false, fmt.Sprintf("Engine: Failed to create action: %v", err)
		}

		return true, fmt.Sprintf("Engine: Action %q created successfully", actionName)

	default:
		return false, "Engine: Unknown create command"
	}
}

func ProcessCd(args []any, engine *Engine) (bool, any) {
	if len(args) == 0 {
		return false, "Engine: Missing CD operation"
	}

	operation, ok := args[0].(string)
	if !ok {
		return false, "Engine: CD operation must be a string"
	}

	switch operation {
	case "app":
		if len(args) < 2 {
			return false, "Engine: Missing application name"
		}

		appName, ok := args[1].(string)
		if !ok {
			return false, "Engine: Application name must be a string"
		}

		if len(args) < 3 {
			return false, "Engine: Missing app path"
		}

		appPath, ok := args[2].(string)
		if !ok {
			return false, "Engine: App path must be a string"
		}

		app := engine.GetAppFromPath(appPath + "_" + appName)
		if app == nil {
			return false, fmt.Sprintf("Engine: Application %q does not exist at path %q", appName, appPath)
		}

		return true, append(strings.Split(appPath, "_"), appName)

	
	default:
		return false, "Engine: Unknown cd command"
	}
}

func ProcessLs(args []any, engine *Engine) (bool, any) {
	if len(args) == 0 {
		return false, "Engine: Missing application path"
	}

	appPath, ok := args[0].(string)
	if !ok {
		return false, "Engine: Application path must be a string"
	}

	children, err := engine.GetAppChildren(appPath)
	if err != nil {
		return false, fmt.Sprintf("Engine: Failed to list application children: %v", err)
	}

	return true, children
}



