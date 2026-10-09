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


func StartCLILoop() {
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

		HandleOperation(operation, args)
	}
}

func SplitCommand(command string) []string {
	return strings.Fields(command)
}

func HandleOperation(operation string, args []string) {
	switch operation {
	case "create":
		ProcessCreate(args)

	case "help":
		ProcessHelp()
	default:
		fmt.Println("Unknown operation:", operation)
	}
}

func ProcessCreate(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "app":
		app := args[0]
		payload := args[1:]


		

		data := CreateCommand(app, payload)

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