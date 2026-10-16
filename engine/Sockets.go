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

func StartSocket() {
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

	if err := ListenForClients(listener); err != nil {
		fmt.Printf("Socket listener stopped: %v\n", err)
	}


}

func ListenForClients(listener net.Listener) error {
	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept client: %w", err)
	}

	HandleClient(conn)
	return nil
}


type Command struct {
	Command string `json:"command"`
	Args    []any  `json:"args"`
}

func HandleClient(conn net.Conn) {
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
			response:=HandleCommand(command)
			// Process command using engine here.
			

			if err := SendResponse(conn, response); err != nil {
				fmt.Printf("Failed to send response: %v\n", err)
			}
		}
	}
}


func HandleCommand(command Command) Response {
	switch command.Command {
	case "create":
		response:=ProcessCreate(command.Args)
		return response
	case "sequence":
		response:=ProcessSequence(command.Args)
		return response
	case "execute":
		response:=ProcessExecute(command.Args)
		return response
	case "cd":
		response:=ProcessCd(command.Args)
		return response
	case "ls":
		response:=ProcessLs(command.Args)
		return response
	default:
		return Response{
			Status: false,
			Msg:    fmt.Sprintf("CLI: Unknown command: %v", command.Command),
			Args:   []any{},
		}

	}
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


func ProcessCreate(args []any) Response {
	operation, err := GetFirstArg(args)
	if err != nil {
		return *err
	}

	switch operation {
	case "app":

	case "action":

	case "sequence":

	default:
		return Response{
			Status: false,
			Msg:    "Engine :Invalid operation.",
		}
	}

	return Response{Status: true}
}

func ProcessSequence(args []any) Response {
	operation, err := GetFirstArg(args)
	if err != nil {
		return *err
	}

	switch operation {
	case "add":
	case "remove":
	case "reorder":
	case "exchange":
	default:
		return Response{
			Status: false,
			Msg:    "Invalid operation. Expected add, remove, reorder, or exchange",
		}
	}

	return Response{Status: true}
}

func ProcessCd(args []any) Response {
	operation, err := GetFirstArg(args)
	if err != nil {
		return *err
	}

	switch operation {
	case "app":
	case "action":
	default:
		return Response{
			Status: false,
			Msg:    "Invalid operation. Expected app or action",
		}
	}

	return Response{Status: true}
}

func ProcessLs(args []any) Response {
	operation, err := GetFirstArg(args)
	if err != nil {
		//here hsow h te child apps 
	}

	switch operation {
	case "actions":
	case "sequences":
	default:
		return Response{
			Status: false,
			Msg:    "Invalid operation. Expected app, actions, or sequences",
		}
	}

	return Response{Status: true}
}

func ProcessExecute(args []any) Response {
	operation, err := GetFirstArg(args)
	if err != nil {
		return *err
	}

	switch operation {
	case "action":
	case "sequence":
	default:
		return Response{
			Status: false,
			Msg:    "Invalid operation. Expected action or sequence",
		}
	}

	return Response{Status: true}
}

func GetFirstArg(args []any) (any, *Response) {
	if len(args) == 0 {
		return nil, &Response{
			Status: false,
			Msg:    "Engine: Missing operation",
			Args:   nil,
		}
	}

	return args[0], nil
}