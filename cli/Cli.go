package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CLI struct {
	unix_socket *net.Conn
}
var Commands = map[string]map[string][]string{
	"create": {
		"app":      {"name"},
		"http":     {"name"},
		"ws":       {"name"},
		"sequence": {"name"},
	},
	"cd": {
		"app":    {"name"},
		"action": {"type: http|ws", "name"},
	},
	"ls": {
		"app":       {},
		"actions":   {},
		"sequences": {},
	},
	"sequence": {
		"add":    {"sequence_name", "action_type: http|ws", "action_name", "position"},
		"remove": {"sequence_name", "position"},
	},
	"execute": {
		"action":   {"type: http|ws", "name"},
		"sequence": {"name"},
	},
	"exit": {},
	"help": {},
}
func CreateCLI(socket *net.Conn) *CLI {
	return &CLI{
		unix_socket: socket,
	}
}

func (cli *CLI) StartCLI() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			continue
		}

		command := input[0]
		args := input[1:]

		switch command {
		case "create":
			cli.ProcessCreate(args)
		case "cd":
			cli.ProcessCD(args)
		case "ls":
			cli.ProcessLS(args)
		case "sequence":
			cli.ProcessSequence(args)
		case "execute":
			cli.ProcessExecute(args)
		case "exit":
			return
		case "help":
			cli.ProcessHelp()
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}



func (cli *CLI) ProcessCreate(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: create <type> <name>")
		return
	}

	operation := args[0]
	params := args[1:]

	var data map[string]any

	switch operation {
	case "app":
		if len(params) != 1 {
			fmt.Println("Usage: create app <name>")
			return
		}

		data = map[string]any{"command": "create", "operation": "app", "name": params[0], "parent_app_id": global.CurrentAppID}

	case "http", "ws":
		if len(params) != 1 {
			fmt.Printf("Usage: create %s <name>\n", operation)
			return
		}

		data = map[string]any{"command": "create", "operation": "action", "type": operation, "name": params[0], "app_id": global.CurrentAppID, "opts": map[string]any{}}

	case "sequence":
		if len(params) != 1 {
			fmt.Println("Usage: create sequence <name>")
			return
		}

		data = map[string]any{"command": "create", "operation": "sequence", "name": params[0], "app_id": global.CurrentAppID}

	default:
		fmt.Printf("Invalid create type: %s\n", operation)
		fmt.Println("Valid types: app, http, ws, sequence")
		return
	}

	response, err := SendData(*cli.unix_socket, data)
	if err != nil {
		fmt.Printf("Failed to create %s: %v\n", operation, err)
		return
	}

	if !response.Status {
		fmt.Printf("Create failed: %s\n", response.Msg)
	}
	fmt.Printf("Engine response: %+v\n", response)
}

func (cli *CLI) ProcessCD(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: cd <app|action> ...")
		return
	}

	operation := args[0]
	params := args[1:]

	var data map[string]any

	switch operation {
	case "app":
		if len(params) != 1 {
			fmt.Println("Usage: cd app <name>")
			return
		}

		data = map[string]any{"command": "cd", "operation": "app", "parent_app_id": global.CurrentAppID, "child_app_name": params[0]}

	case "action":
		if len(params) != 2 {
			fmt.Println("Usage: cd action <http|ws> <name>")
			return
		}

		if params[0] != "http" && params[0] != "ws" {
			fmt.Printf("Invalid action type: %s\n", params[0])
			fmt.Println("Valid action types: http, ws")
			return
		}

		data = map[string]any{"command": "cd", "operation": "action", "app_id": global.CurrentAppID, "type": params[0], "action_name": params[1]}

		response, err := SendData(*cli.unix_socket, data)
		if err != nil {
			fmt.Printf("Failed to retrieve action: %v\n", err)
			return
		}

		if !response.Status {
			fmt.Println(response.Msg)
			return
		}

		editedActionType, editedOpts, err := EditActionOpts(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("CLI: Edited options:")
		fmt.Println(string(editedOpts))

		saveData := map[string]any{
			"command":   "cd",
			"operation": "save",
			"app_id":    global.CurrentAppID,
			"type":      editedActionType,
			"action_name": params[1],
			"opts":      json.RawMessage(editedOpts),
		}

		saveResponse, err := SendData(*cli.unix_socket, saveData)
		if err != nil {
			fmt.Printf("Failed to save action options: %v\n", err)
			return
		}

		if !saveResponse.Status {
			fmt.Println(saveResponse.Msg)
			return
		}

		fmt.Println(saveResponse.Msg)

		return

		
	default:
		fmt.Printf("Invalid cd type: %s\n", operation)
		fmt.Println("Valid types: app, action")
		return
	}

	response, err := SendData(*cli.unix_socket, data)
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	}

	if !response.Status {
		fmt.Printf("CD failed: %s\n", response.Msg)
	}
}

func (cli *CLI) ProcessLS(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ls <app|actions|sequences>")
		return
	}

	operation := args[0]
	var data map[string]any

	switch operation {
	case "app":
		data = map[string]any{"command": "ls", "operation": "", "app_id": global.CurrentAppID}
	case "actions", "sequences":
		data = map[string]any{"command": "ls", "operation": operation, "app_id": global.CurrentAppID}
	default:
		fmt.Printf("Invalid ls type: %s\n", operation)
		fmt.Println("Valid types: app, actions, sequences")
		return
	}

	response, err := SendData(*cli.unix_socket, data)
	if err != nil {
		fmt.Printf("Failed to list %s: %v\n", operation, err)
		return
	}

	if !response.Status {
		fmt.Printf("Failed to list %s: %s\n", operation, response.Msg)
	}
	fmt.Printf("Engine response: %+v\n", response)
}

func (cli *CLI) ProcessSequence(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: sequence <add|remove> ...")
		return
	}

	operation := args[0]
	params := args[1:]
	var data map[string]any

	switch operation {
	case "add":
		if len(params) != 4 {
			fmt.Println("Usage: sequence add <sequence_name> <http|ws> <action_name> <position>")
			return
		}

		position, err := strconv.Atoi(params[3])
		if err != nil {
			fmt.Printf("Invalid position: %s\n", params[3])
			return
		}

		if params[1] != "http" && params[1] != "ws" {
			fmt.Printf("Invalid action type: %s\n", params[1])
			fmt.Println("Valid action types: http, ws")
			return
		}

		data = map[string]any{"command": "sequence", "operation": "add", "app_id": global.CurrentAppID, "sequence_name": params[0], "type": params[1], "action_name": params[2], "position": position}

	case "remove":
		if len(params) != 2 {
			fmt.Println("Usage: sequence remove <sequence_name> <position>")
			return
		}

		position, err := strconv.Atoi(params[1])
		if err != nil {
			fmt.Printf("Invalid position: %s\n", params[1])
			return
		}

		data = map[string]any{"command": "sequence", "operation": "remove", "app_id": global.CurrentAppID, "sequence_name": params[0], "position": position}

	default:
		fmt.Printf("Invalid sequence operation: %s\n", operation)
		fmt.Println("Valid operations: add, remove")
		return
	}

	response, err := SendData(*cli.unix_socket, data)
	if err != nil {
		fmt.Printf("Failed to process sequence: %v\n", err)
		return
	}

	if !response.Status {
		fmt.Printf("Sequence operation failed: %s\n", response.Msg)
	}
}



func (cli *CLI) ProcessExecute(args []string) {
}



func (cli *CLI) ProcessHelp() {
	fmt.Println("Available commands:")

	for command, types := range Commands {
		if len(types) == 0 {
			fmt.Printf("  %s\n", command)
			continue
		}

		for typ, args := range types {
			fmt.Printf("  %s %s", command, typ)

			for _, arg := range args {
				fmt.Printf(" <%s>", arg)
			}

			fmt.Println()
		}
	}
}



func EditActionOpts(response Response) (string, []byte, error) {
	if !response.Status {
		return "", nil, fmt.Errorf("%s", response.Msg)
	}

	if len(response.Args) < 2 {
		return "", nil, fmt.Errorf("CLI: Invalid action options response")
	}

	actionType, ok := response.Args[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("CLI: Invalid action type")
	}

	optsJSON, err := json.MarshalIndent(response.Args[1], "", "    ")
	if err != nil {
		return "", nil, fmt.Errorf("CLI: Failed to encode action options: %w", err)
	}

	editedJSON, err := OpenEditor(optsJSON)
	if err != nil {
		return "", nil, err
	}

	if !json.Valid(editedJSON) {
		return "", nil, fmt.Errorf("CLI: Edited options contain invalid JSON")
	}

	return actionType, editedJSON, nil
}




func OpenEditor(data []byte) ([]byte, error) {
	file, err := os.CreateTemp("", "testflow-*.json")
	if err != nil {
		return nil, err
	}

	filePath := file.Name()
	defer os.Remove(filePath)

	if _, err := file.Write(data); err != nil {
		file.Close()
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	editor := global.Editor
	if editor == "" {
		return nil, fmt.Errorf("CLI: No editor configured in global.Editor")
	}

	if _, err := exec.LookPath(editor); err != nil {
		return nil, fmt.Errorf("CLI: Editor %q is not installed or available in PATH", editor)
	}

	switch editor {
	case "nano":
		err = ExecuteNano(filePath)
	case "vim":
		// err = ExecuteVim(filePath)
	case "code":
		cmd := exec.Command("code", "--wait", filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	default:
		cmd := exec.Command(editor, filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	if err != nil {
		return nil, fmt.Errorf("CLI: Failed to run editor: %w", err)
	}

	editedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if !json.Valid(editedData) {
		return nil, fmt.Errorf("CLI: Edited file contains invalid JSON")
	}

	return editedData, nil
}

func ExecuteNano(filePath string) error {
	cmd := exec.Command("nano", "-m", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}