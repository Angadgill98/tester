package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
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
	Args    []string `json:"args"`
}


func (cli *CLI)StartCLILoop() {
	fmt.Println("========== CLI START ==========")
	fmt.Println("Available commands:")
	fmt.Println("exit")
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s > ", strings.Join(global.AppPath, "/"))

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

	case "help":
		ProcessHelp(args)

	case "cd":
		cli.ProcessCd(args)

	case "ls":
		cli.ProcessLs(args)

	case "pwd":
		fmt.Println(strings.Join(global.AppPath, "/"))

	default:
		fmt.Println("Unknown operation:", operation)
	}
}

func CreateCommand(command string, args []string) []byte {
	commandObj := Command{
		Command: command,
		Args:    args,
	}

	data, err := json.Marshal(commandObj)
	if err != nil {
		fmt.Println("Failed to serialize command:", err)
		return nil
	}

	return data
}

func (cli *CLI)ProcessCreate(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "app":
		
		response:=cli.SendPipeline("create",args)

		fmt.Println(response.Response)

	case "action":
		response:=cli.SendPipeline("create",args)
		

		fmt.Println(response.Response)

	}
}

type Response struct {
	Status   bool `json:"status"`
	Response any  `json:"response"`
}

func (cli *CLI)SendPipeline(com string,args []string)Response{
	current_path := strings.Join(global.AppPath, "_")
	payload := args[0:]
	payload = append(payload, current_path)

	command := CreateCommand(com, payload)
	var response Response

	data, err := cli.SendCommand(command)
	if err != nil {
		response.Status=false
		response.Response=fmt.Sprintf("CLI: Failed to send command: %v\n", err)
		return response
	}


	if err := json.Unmarshal(data, &response); err != nil {
		response.Status=false
		response.Response=fmt.Sprintf("CLI: Failed to parse response: %v\n", err)
		return response
	}
	return response
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



func (cli *CLI) ProcessCd(args []string){
	if len(args) == 0 {
		return
	}

	switch args[0] {
		case "app":

			response:=cli.SendPipeline("cd",args)

			if !response.Status {
				fmt.Println("CLI: Failed , Reponse is",response.Response)
				return
			}

			appPathItems, ok := response.Response.([]any)
			if !ok {
				fmt.Println("CLI: Invalid application path response")
				return
			}

			appPath := make([]string, 0, len(appPathItems))
			for _, item := range appPathItems {
				appPath = append(appPath, item.(string))
			}
			if !ok {
				fmt.Println("CLI: Invalid application path response")
				fmt.Println("Reponse is",response.Response)

				return
			}

			global.AppPath = appPath
			// fmt.Println("Reponse is",response.Response)
		
		case "..":
		if len(global.AppPath) > 1 {
			global.AppPath = global.AppPath[:len(global.AppPath)-1]
		}

		default:
			fmt.Println("CLI: Unknown CD operation:", args[0])
	}
}


func (cli *CLI) ProcessLs(args []string) {
	response:=cli.SendPipeline("ls",args)

	if !response.Status {
		fmt.Println(response.Response)
		return
	}

	switch items := response.Response.(type) {
	case []any:
		for _, item := range items {
			fmt.Println(item)
		}
	default:
		fmt.Println(response.Response)
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





