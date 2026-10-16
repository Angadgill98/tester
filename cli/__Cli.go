package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)


type CLI struct {
	unix_socket *net.Conn
}

func CreateCLI(socket *net.Conn ) *CLI {
	return &CLI{
		unix_socket: socket,
	}
}
type Command struct {
	Command string   `json:"command"`
	Args    []any `json:"args"`
}


func (cli *CLI)StartCLILoop() {
	fmt.Println("========== CLI START ==========")
	fmt.Println("Available commands:")
	fmt.Println("exit")
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s > ", strings.Join(global.CurrentAppPath, "/"))

		if !scanner.Scan() {
			return
		}

		command := scanner.Text()
		payload := SplitCommand(command)

		if len(payload) == 0 {
			continue
		}

		operation := payload[0]
		args := payload[1:]

		cli.HandleOperation(operation, args)
	}
}

func SplitCommand(command string) []string {
	return strings.Fields(command)
}

func (cli *CLI)HandleOperation(operation string, args []string) {
	switch operation {
	case "create":
		cli.ProcessCreate(args)

	case "sequence":
		cli.ProcessSequence(args)

	case "execute":
		cli.ProcessExecute(args)

	case "help":
		ProcessHelp(args)

	case "cd":
		cli.ProcessCd(args)

	case "ls":
		cli.ProcessLs(args)

	case "pwd":
		fmt.Println(strings.Join(global.CurrentAppPath, "/"))

	default:
		fmt.Println("CLI: Unknown operation:", operation)
	}
}

func (cli *CLI)ProcessCreate(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
		case "app":
			if len(args) < 2 {
				fmt.Println("CLI: Application name is missing")
				return
			}

			appName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
			currentAppID := global.CurrentAppID
			payload := []any{"app",appName, currentAppID}

			response := cli.SendPayload("create", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}else{
				fmt.Println(response.Msg)
			}

			for _, item := range response.Args {
				fmt.Println(item)
			}

		case "action":
			if len(args) < 2 {
				fmt.Println("CLI: Action name is missing")
				return
			}

			actionName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))

			var actionType string

			for {
				fmt.Print("Enter action type (http/ws, 0 to exit): ")

				fmt.Scanln(&actionType)

				if actionType == "0" {
					return
				}

				actionType = strings.ToLower(strings.TrimSpace(actionType))

				if actionType == "http" || actionType == "ws" {
					break
				}

				fmt.Println("CLI: Invalid action type. Enter http, ws, or 0 to exit.")
			}

			var opts Actions_opts

			switch strings.ToLower(actionType) {
			case "http":
				opts = Actions_opts{
					Action_type: "http",
					HTTP:        HTTP_opts{},
				}

			case "ws":
				opts = Actions_opts{
					Action_type: "ws",
					WS:          WS_opts{},
				}

			default:
				fmt.Println("CLI: Invalid action type. Use http or ws")
				return
			}

			currentAppID := global.CurrentAppID
			payload := []any{"action", actionName, currentAppID, opts}

			response := cli.SendPayload("create", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}

			fmt.Println(response.Msg)

			for _, item := range response.Args {
				fmt.Println(item)
			}
		
		case "sequence":
			if len(args) < 2 {
				fmt.Println("CLI: Sequence name is missing")
				return
			}

			sequenceName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
			currentAppID := global.CurrentAppID
			payload := []any{"sequence", sequenceName, currentAppID}

			response := cli.SendPayload("create", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}

			fmt.Println(response.Msg)

			for _, item := range response.Args {
				fmt.Println(item)
			}

		default:
			fmt.Println("CLI: Invalid create command. Use app, action, or sequence.")

	}
}

func (cli *CLI) ProcessSequence(args []string) {
	if len(args) < 1 {
		fmt.Println("CLI: Missing sequence operation")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			fmt.Println("CLI: Expected sequence name and action name")
			return
		}

		sequenceName := args[1]
		actionName := args[2]
		currentAppID := global.CurrentAppID

		response := cli.SendPayload("sequence", []any{"add", sequenceName, actionName, currentAppID})
		if !response.Status {
			fmt.Println(response.Msg)
			return
		}
		fmt.Println(response.Msg)

	case "remove":
		if len(args) < 3 {
			fmt.Println("CLI: Expected sequence name and action name")
			return
		}

		sequenceName := args[1]
		actionName := args[2]
		currentAppID := global.CurrentAppID

		response := cli.SendPayload("sequence", []any{"remove", sequenceName, actionName, currentAppID})
		if !response.Status {
			fmt.Println(response.Msg)
			return
		}
		fmt.Println(response.Msg)

	default:
		fmt.Println("CLI: Invalid sequence operation. Use add or remove.")
	}
}

func (cli *CLI) ProcessExecute(args []string) {
	if len(args) < 2 {
		fmt.Println("CLI: Invalid execute arguments")
		return
	}

	operation := args[0]

	switch operation {
	case "action":
		if len(args) < 2 {
			fmt.Println("CLI: Action name is missing")
			return
		}

		actionName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
		appID := global.CurrentAppID

		payload := []any{"action", actionName, appID}
		response := cli.SendPayload("execute", payload)

		if !response.Status {
			fmt.Println(response.Msg)
			return
		}

		fmt.Println(response.Msg)
		for _, arg := range response.Args {
			fmt.Printf("%+v\n", arg)
		}

	case "sequence":
		if len(args) < 2 {
			fmt.Println("CLI: Sequence name is missing")
			return
		}

		sequenceName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
		appID := global.CurrentAppID

		payload := []any{"sequence", sequenceName, appID}
		response := cli.SendPayload("execute", payload)

		if !response.Status {
			fmt.Println(response.Msg)
			return
		}

		fmt.Println(response.Msg)
		for _, arg := range response.Args {
			fmt.Printf("%+v\n", arg)
		}

	case "interactive":
		fmt.Println("CLI: Starting temporary sequence")

		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print("Enter action name (0 to exit): ")

			if !scanner.Scan() {
				break
			}

			input := strings.TrimSpace(scanner.Text())

			if input == "0" {
				break
			}

			args := strings.Fields(input)
			if len(args) != 1 {
				fmt.Println("CLI: Enter exactly one action name")
				continue
			}

			actionName := args[0]
			appID := global.CurrentAppID

			payload := []any{"action", actionName, appID}
			response := cli.SendPayload("execute", payload)

			if !response.Status {
				fmt.Println(response.Msg)
				continue
			}

			fmt.Println(response.Msg)
			for _, arg := range response.Args {
				fmt.Printf("%+v\n", arg)
			}
		}

		fmt.Println("CLI: Temporary sequence completed")
	default:
		fmt.Printf("CLI: Unknown execute operation %q\n", operation)
	}
}

type Response struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Args   []any  `json:"args"`
}




var commandHelp = map[string]any{
	"create": map[string]any{
		"syntax": "create <type> <args...>",
		"app":    "create app <name>",
	},
	"help": map[string]any{
		"syntax": "help [command]",
	},
	"exit": map[string]any{
		"syntax": "exit",
	},
}
func ProcessHelp(args []string) {
	if len(args) == 0 {
		fmt.Println("Available commands:")

		for command, value := range commandHelp {
			data := value.(map[string]any)
			fmt.Println(command, "-", data["syntax"])
		}

		return
	}

	commandData, exists := commandHelp[args[0]]
	if !exists {
		fmt.Println("Unknown command:", args[0])
		return
	}

	data := commandData.(map[string]any)

	if len(args) == 1 {
		fmt.Println(data["syntax"])
		return
	}

	subCommand, exists := data[args[1]]
	if !exists {
		fmt.Println("Unknown option:", args[1])
		return
	}

	fmt.Println(subCommand)
}


func (cli *CLI) SendPayload(command string, payload []any) Response {
	request := Command{
		Command: command,
		Args:    payload,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return Response{Status: false, Msg: fmt.Sprintf("CLI: Failed to encode request: %v", err)}
	}

	data = append(data, '\n')

	if _, err := (*cli.unix_socket).Write(data); err != nil {
		return Response{Status: false, Msg: fmt.Sprintf("CLI: Failed to send request: %v", err)}
	}

	var response Response
	if err := json.NewDecoder(*cli.unix_socket).Decode(&response); err != nil {
		return Response{Status: false, Msg: fmt.Sprintf("CLI: Failed to receive response: %v", err)}
	}

	return response
}

func (cli *CLI) ProcessCd(args []string){
	if len(args) == 0 {
		return
	}

	switch args[0] {
		case "app":
			if len(args) < 2 {
				fmt.Println("CLI: Application name is missing")
				return
			}

			appName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
			currentAppID := global.CurrentAppID
			payload := []any{"app",appName, currentAppID}

			response := cli.SendPayload("cd", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}

			fmt.Println(response.Msg)

			if len(response.Args) < 2 {
				fmt.Println("CLI: Invalid response arguments")
				return
			}

			appName, ok := response.Args[0].(string)
			if !ok {
				fmt.Println("CLI: Invalid application name")
				return
			}

			appID, ok := response.Args[1].(string)
			if !ok {
				fmt.Println("CLI: Invalid application ID")
				return
			}

			global.CurrentAppPath = append(global.CurrentAppPath, appName)
			global.CurrentAppID = appID

		case "action":
			if len(args) < 2 {
				fmt.Println("CLI: Action name is missing")
				return
			}

			actionName := strings.TrimSpace(strings.ReplaceAll(strings.Join(args[1:], " "), " ", "_"))
			currentAppID := global.CurrentAppID

			payload := []any{"action",actionName, currentAppID}

			response := cli.SendPayload("cd", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}

			actionType, editedOpts, err := EditActionOpts(response)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("CLI: Edited options:")
			fmt.Println(string(editedOpts))

			savePayload := []any{"save", actionName, currentAppID, actionType, json.RawMessage(editedOpts)}
			saveResponse := cli.SendPayload("cd", savePayload)
			if !saveResponse.Status {
				fmt.Println(saveResponse.Msg)
				return
			}

			fmt.Println(saveResponse.Msg)

		case "..":
			currentAppID := global.CurrentAppID
			payload := []any{currentAppID}

			response := cli.SendPayload("..", payload)
			if !response.Status {
				fmt.Println(response.Msg)
				return
			}else{
				fmt.Println(response.Msg)
			}

			for _, item := range response.Args {
				fmt.Println(item)
			}

		default:
			fmt.Println("CLI: Unknown CD operation:", args[0])
	}
}


func (cli *CLI) ProcessLs(args []string) {
	payload := []any{global.CurrentAppID}
	response := cli.SendPayload("ls", payload)

	if !response.Status {
		fmt.Println(response.Msg)
		return
	}else{
		fmt.Println(response.Msg)
	}

	for _, item := range response.Args {
		fmt.Println(item)
	}
}



func (cli *CLI) SendCommand(command []byte) ([]byte, error) {
	_, err := (*cli.unix_socket).Write(command)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	buffer := make([]byte, 4096)
	n, err := (*cli.unix_socket).Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to receive response: %w", err)
	}

	return buffer[:n], nil
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