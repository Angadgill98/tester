package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"engine/Actions"
	"engine/Global"
	"engine/Sequences"

	"github.com/google/uuid"
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

func ListenForClients(listener net.Listener, engine *Engine) error {
	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept client: %w", err)
	}

	if err := SendGlobalAppID(conn, engine); err != nil {
		conn.Close()
		return fmt.Errorf("failed to send GLOBAL application ID: %w", err)
	}

	HandleClient(conn, engine)
	return nil
}

func SendGlobalAppID(conn net.Conn, engine *Engine) error {
	for id, app := range *engine.App_config {
		if app.Name == "GLOBAL" {
			return SendResponse(conn, Response{
				Status: true,
				Msg:    "GLOBAL application ID",
				Args:   []any{app.Name, id.String()},
			})
		}
	}

	return fmt.Errorf("GLOBAL application does not exist")
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

			fmt.Printf("Engine: Command: %s\n", command.Command)
			fmt.Printf("Engine: Args: %v\n", command.Args)
			response:=HandleCommand(command, engine)
			// Process command using engine here.
			

			if err := SendResponse(conn, response); err != nil {
				fmt.Printf("Failed to send response: %v\n", err)
			}
		}
	}
}


func HandleCommand(command Command, engine *Engine) Response {
	switch command.Command {
	case "create":
		response:=ProcessCreate(command.Args, engine)
		return response
	case "cd":
		response:=ProcessCd(command.Args,engine)
		return response
	case "ls":
		response:=ProcessLs(command.Args,engine)
		return response
	default:
		return Response{
			Status: false,
			Msg:    fmt.Sprintf("CLI: Unknown command: %v", command.Command),
			Args:   []any{},
		}

	}
}

type Response struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Args   []any  `json:"args"`
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


func ProcessCreate(args []any, engine *Engine) Response {
	if len(args) == 0 {
		return Response{Status: false, Msg: "Engine: Missing create operation", Args: []any{}}
	}

	operation, ok := args[0].(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: Create operation must be a string", Args: []any{}}
	}

	switch operation {
		case "app":
			if len(args) < 3 {
				return Response{Status: false, Msg: "Engine: Expected app name and parent application UUID", Args: []any{}}
			}

			appName, ok := args[1].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: App name must be a string", Args: []any{}}
			}

			parentIDString, ok := args[2].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: Parent application ID must be a UUID string", Args: []any{}}
			}

			parentID, err := uuid.Parse(parentIDString)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid parent application UUID: %v", err), Args: []any{}}
			}

			if err := engine.CreateApp(appName, parentID); err != nil {
				return Response{Status: false, Msg: err.Error(), Args: []any{}}
			}

			return Response{Status: true, Msg: fmt.Sprintf("Engine: Application %q created successfully", appName), Args: []any{}}

		case "action":
			if len(args) < 4 {
				return Response{Status: false, Msg: "Engine: Expected action name, application UUID, and action options", Args: []any{}}
			}

			actionName, ok := args[1].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: Action name must be a string", Args: []any{}}
			}

			appIDString, ok := args[2].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: Application ID must be a UUID string", Args: []any{}}
			}

			appID, err := uuid.Parse(appIDString)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid application UUID: %v", err), Args: []any{}}
			}

			optsJSON, err := json.Marshal(args[3])
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Failed to encode action options: %v", err), Args: []any{}}
			}

			var rawOpts struct {
				ActionType string          `json:"action_type"`
				HTTP       json.RawMessage `json:"http"`
				WS         json.RawMessage `json:"ws"`
			}

			if err := json.Unmarshal(optsJSON, &rawOpts); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Failed to decode action options: %v", err), Args: []any{}}
			}

			var opts any

			switch rawOpts.ActionType {
			case "http":
				var httpOpts Actions.HTTP_opts
				if err := json.Unmarshal(rawOpts.HTTP, &httpOpts); err != nil {
					return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid HTTP action options: %v", err), Args: []any{}}
				}
				opts = httpOpts

			case "ws":
				var wsOpts Actions.WS_opts
				if err := json.Unmarshal(rawOpts.WS, &wsOpts); err != nil {
					return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid WebSocket action options: %v", err), Args: []any{}}
				}
				opts = wsOpts

			default:
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Unsupported action type %q", rawOpts.ActionType), Args: []any{}}
			}

			if err := engine.CreateAction(appID, actionName, opts); err != nil {
				return Response{Status: false, Msg: err.Error(), Args: []any{}}
			}

			return Response{Status: true, Msg: fmt.Sprintf("Engine: Action %q created successfully", actionName), Args: []any{}}
		
		case "sequence":
			if len(args) < 3 {
				return Response{Status: false, Msg: "Engine: Expected sequence name and application UUID", Args: []any{}}
			}

			sequenceName, ok := args[1].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: Sequence name must be a string", Args: []any{}}
			}

			appIDString, ok := args[2].(string)
			if !ok {
				return Response{Status: false, Msg: "Engine: Application ID must be a UUID string", Args: []any{}}
			}

			appID, err := uuid.Parse(appIDString)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid application UUID: %v", err), Args: []any{}}
			}

			sequence := Sequences.CreateSequence(sequenceName)

			(*engine.Sequence_config)[sequence.ID] = *sequence

			app, exists := (*engine.App_config)[appID]
			if !exists {
				delete(*engine.Sequence_config, sequence.ID)
				return Response{Status: false, Msg: "Engine: Application not found", Args: []any{}}
			}

			app.Sequences[sequenceName] = sequence.ID

			if err := engine.Sequence_config.SaveSequenceConfig(); err != nil {
				delete(*engine.Sequence_config, sequence.ID)
				delete(app.Sequences, sequenceName)
				return Response{Status: false, Msg: err.Error(), Args: []any{}}
			}

			if err := engine.App_config.SaveAppConfig(); err != nil {
				delete(*engine.Sequence_config, sequence.ID)
				delete(app.Sequences, sequenceName)
				return Response{Status: false, Msg: err.Error(), Args: []any{}}
			}

			return Response{Status: true, Msg: fmt.Sprintf("Engine: Sequence %q created successfully", sequenceName), Args: []any{}}
			

		
		default:
			return Response{Status: false, Msg: "Engine: Unknown create command", Args: []any{}}
	}
}

func ProcessCd(args []any, engine *Engine) Response {
	if len(args) < 2 {
		return Response{Status: false, Msg: "Engine: Expected application name and parent application UUID", Args: []any{}}
	}

	operation, ok := args[0].(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: CD operation must be a string", Args: []any{}}
	}

	switch operation {
	case "app":
		appName, ok := args[1].(string)
		if !ok {
			return Response{Status: false, Msg: "Engine: Application name must be a string", Args: []any{}}
		}

		if len(args) < 3 {
			return Response{Status: false, Msg: "Engine: Missing parent application UUID", Args: []any{}}
		}

		parentIDString, ok := args[2].(string)
		if !ok {
			return Response{Status: false, Msg: "Engine: Parent application ID must be a UUID string", Args: []any{}}
		}

		parentID, err := uuid.Parse(parentIDString)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid parent application UUID: %v", err), Args: []any{}}
		}

		parentApp, exists := (*engine.App_config)[parentID]
		if !exists {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: Parent application with ID %q does not exist", parentIDString), Args: []any{}}
		}

		childID, exists := parentApp.Children[appName]
		if !exists {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: Application %q does not exist under parent %q", appName, parentApp.Name), Args: []any{}}
		}

		if _, exists := (*engine.App_config)[childID]; !exists {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: Application %q references a missing application ID %q", appName, childID), Args: []any{}}
		}

		return Response{Status: true, Msg: "Engine: Application selected successfully", Args: []any{appName, childID.String()}}

	default:
		return Response{Status: false, Msg: "Engine: Unknown cd command", Args: []any{}}
	}
}

func ProcessLs(args []any, engine *Engine) Response {
	if len(args) == 0 {
		return Response{Status: false, Msg: "Engine: Missing application UUID", Args: []any{}}
	}

	appIDString, ok := args[0].(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: Application ID must be a UUID string", Args: []any{}}
	}

	appID, err := uuid.Parse(appIDString)
	if err != nil {
		return Response{Status: false, Msg: fmt.Sprintf("Engine: Invalid application UUID: %v", err), Args: []any{}}
	}

	app, exists := (*engine.App_config)[appID]
	if !exists {
		return Response{Status: false, Msg: fmt.Sprintf("Engine: Application with ID %q does not exist", appIDString), Args: []any{}}
	}

	children := make([]any, 0, len(app.Children))
	for name := range app.Children {
		children = append(children, name)
	}

	return Response{Status: true, Msg: "Engine: Children retrieved successfully", Args: children}
}


