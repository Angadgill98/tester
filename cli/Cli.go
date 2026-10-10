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
		fmt.Print("> ")

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
		
		current_path := strings.Join(global.AppPath, "_")
		payload := args[0:]
		payload = append(payload, current_path)

		command := CreateCommand("create", payload)
		cli.SendCommand(command)
	}
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
			commad:=CreateCommand("cd",args)
			cli.SendCommand(commad)

			break;
			

		default:
	}
}


func (cli *CLI) ProcessLs(args []string){
	if len(args) == 0 {
		commad:=CreateCommand("ls",args)
		cli.SendCommand(commad)

	}

	
}


func (cli *CLI) SendCommand(command []byte) error {
	_, err := (*cli.unix_socket).Write(command)
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	return nil
}



