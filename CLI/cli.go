package CLI

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	Actions "tester/Actions"
	Sequences "tester/Sequences"
	tester "tester/Tester"
)

type CLI struct {
	actions *Actions.Actions_obj
	sequences *Sequences.Sequences_obj
}

func CreateCLI(actions *Actions.Actions_obj,sequences *Sequences.Sequences_obj) *CLI {
	return &CLI{
		actions: actions,
		sequences: sequences,
	}
}

func (cli *CLI) Init() {
	fmt.Println("========== CLI START ==========")
	fmt.Println("Available commands:")
	fmt.Println("  create action <action_name>")
	fmt.Println("  create sequence <sequence_name>")
	fmt.Println("  delete action <action_name>")
	fmt.Println("  sequence add-action <sequence_name>")
	fmt.Println("  execute actions")
	fmt.Println("  execute sequence <sequence_name>")
	fmt.Println("  show")
	fmt.Println("  exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "exit" {
			break
		}

		cli.ProcessCommand(input)
	}

	fmt.Println("========== CLI END ==========")
}

func (cli *CLI) ProcessCommand(input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return
	}

	command := parts[0]
	remaining := strings.TrimSpace(strings.TrimPrefix(input, command))

	switch command {
	case "create":
		cli.ProcessCreate(remaining)
	
	case "delete":
		cli.ProcessDelete(remaining)

	case "sequence":
		cli.ProcessSequence(remaining)

	case "execute":
		cli.ProcessExecute(remaining)

	case "show":
		cli.ShowState()

	default:
		fmt.Println("Unknown command:", command)
	}
}

func (cli *CLI) ProcessCreate(input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		fmt.Println("Create what?")
		return
	}

	switch parts[0] {
	case "action":
		actionName := strings.TrimSpace(strings.TrimPrefix(input, "action"))

		if actionName == "" {
			fmt.Println("Usage: create action <action_name>")
			return
		}

		actionName = strings.ReplaceAll(actionName, " ", "_")

		var actionType string

		scanner := bufio.NewScanner(os.Stdin)

		validTypes := map[string]bool{
			"http":   true,
		}

		for {
			fmt.Print("Enter action type (http/..) or 'exit': ")

			if !scanner.Scan() {
				return
			}

			actionType = strings.TrimSpace(scanner.Text())

			if actionType == "exit" {
				return
			}

			if validTypes[actionType] {
				break
			}

			fmt.Println("Invalid action type. Please try again.")
		}

		opts := Actions.Action_opts{
			Action_type: actionType,
		}

		created := cli.actions.CreateAction(actionName, opts)


		if created {
			fmt.Println("Action created successfully:", actionName)
		} else {
			fmt.Println("Action alredu exist")
		}

	case "sequence":
		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

		if sequenceName == "" {
			fmt.Println("Usage: create sequence <sequence_name>")
			return
		}

		sequenceName = strings.ReplaceAll(sequenceName, " ", "_")

		created := cli.sequences.CreateSequence(sequenceName)

		if created {
			fmt.Println("Sequence created successfully:", sequenceName)
		} else {
			fmt.Println("Sequence already exists")
		}

	default:
		fmt.Println("Unknown create command:", parts[0])
	}
}


func (cli *CLI) ProcessDelete(input string) {
	parts := strings.Fields(input)

	if len(parts) < 2 {
		fmt.Println("Usage: delete action <action_name>")
		return
	}

	switch parts[0] {
	case "action":
		actionName := strings.TrimSpace(strings.TrimPrefix(input, "action"))

		if actionName == "" {
			fmt.Println("Usage: delete action <action_name>")
			return
		}

		actionName = strings.ReplaceAll(actionName, " ", "_")

		deleted := cli.actions.DeleteAction(actionName)

		if deleted {
			fmt.Println("Action deleted successfully:", actionName)
		} else {
			fmt.Println("Action does not exist:", actionName)
		}

	case "sequence":
		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

		if sequenceName == "" {
			fmt.Println("Usage: delete sequence <sequence_name>")
			return
		}

		sequenceName = strings.ReplaceAll(sequenceName, " ", "_")

		deleted := cli.sequences.DeleteSequence(sequenceName)

		if deleted {
			fmt.Println("Sequence deleted successfully:", sequenceName)
		} else {
			fmt.Println("Sequence does not exist:", sequenceName)
		}

	default:
		fmt.Println("Unknown delete command:", parts[0])
	}
}



func (cli *CLI) ProcessSequence(input string) {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		fmt.Println("Usage: sequence add-action <sequence_name> <action_name>")
		return
	}

	switch parts[0] {
	case "add-action":
		if len(parts) < 2 {
			fmt.Println("Usage: sequence add-action <sequence_name>")
			return
		}

		sequenceName := parts[1]
		_, exists := cli.sequences.Sequences[sequenceName]

		if !exists {
			fmt.Println("Sequence does not exist:", sequenceName)
			return
		}

		for {
			fmt.Println()
			fmt.Println("Available Actions:")
			fmt.Println("exit. Exit")
			fmt.Println("0. Back")

			actionNames := make([]string, 0, len(cli.actions.Actions))

			i := 1
			for actionName := range cli.actions.Actions {
				fmt.Printf("%d. %s\n", i, actionName)
				actionNames = append(actionNames, actionName)
				i++
			}

			if len(actionNames) == 0 {
				fmt.Println("No actions available")
				return
			}

			var choice string

			fmt.Print("Choose action number: ")

			_, err := fmt.Scanln(&choice)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}

			if choice == "exit" {
				return
			}

			if choice == "0" {
				break
			}

			var actionNumber int

			_, err = fmt.Sscanf(choice, "%d", &actionNumber)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}

			if actionNumber < 1 || actionNumber > len(actionNames) {
				fmt.Println("Invalid action number")
				continue
			}

			actionName := actionNames[actionNumber-1]

			// Add action to sequence here
			action := cli.actions.Actions[actionName]

			fmt.Println("Selected action:", actionName)
			fmt.Println("Adding it to the sequence:", sequenceName)


	
			added := cli.sequences.AddActionToSequence(sequenceName, action)

			if added {
				fmt.Println("Action added successfully:", actionName)
			} else {
				fmt.Println("Failed to add action:", actionName)
			}


		}


		

	

		

	default:
		fmt.Println("Unknown sequence command:", parts[0])
	}
}


func (cli *CLI) ShowState() {
	fmt.Println()
	fmt.Println("========== CURRENT STATE ==========")

	fmt.Println()
	fmt.Println("---------- ACTIONS ----------")
	fmt.Printf("Total Actions: %d\n", len(cli.actions.Actions))

	for actionName, action := range cli.actions.Actions {
		fmt.Println("Action:", actionName)
		fmt.Println("  Name:", action.Name)
		fmt.Println("  Path:", action.Path)
	}

	fmt.Println()
	fmt.Println("---------- ACTION CONFIG ----------")
	fmt.Printf("Total Config Entries: %d\n", len(cli.actions.Config))

	for actionName, executor := range cli.actions.Config {
		if executor != nil {
			fmt.Printf("Config: %s -> function installed\n", actionName)
		} else {
			fmt.Printf("Config: %s -> function NOT installed, need to restart\n", actionName)
		}
	}

	fmt.Println()
	fmt.Println("---------- SEQUENCES ----------")
	fmt.Printf("Total Sequences: %d\n", len(cli.sequences.Sequences))

	for sequenceName, sequence := range cli.sequences.Sequences {
		fmt.Println("Sequence:", sequenceName)
		fmt.Printf("  Actions: %d\n", len(sequence.Actions_to_execute))

		for i, action := range sequence.Actions_to_execute {
			fmt.Printf("    %d. %s\n", i+1, action.Name)
		}
	}

	fmt.Println()
	fmt.Println("---------- SEQUENCE CONFIG ----------")
	fmt.Printf("Total Config Entries: %d\n", len(cli.sequences.Config))

	for sequenceName, actionNames := range cli.sequences.Config {
		fmt.Println("Sequence:", sequenceName)

		for i, actionName := range actionNames {
			fmt.Printf("  %d. %s\n", i+1, actionName)
		}
	}

	fmt.Println()
	fmt.Println("===================================")
	fmt.Println()
}


func (cli *CLI) ProcessExecute(input string) {
	parts := strings.Fields(input)

	if len(parts) < 1 {
		fmt.Println("Usage: execute action <action_name>")
		fmt.Println("       execute sequence <sequence_name>")
		return
	}

	switch parts[0] {
	case "action":
		fmt.Println("Starting Temp Sequence")
		var tester=tester.CreateTesterObj()
		var counter=0;
		for {
			fmt.Println()
			fmt.Println("========== EXECUTE ACTION ==========")
			fmt.Println("Available Actions:")
			fmt.Println("exit. Exit")
			fmt.Println("0. Back")

			actionNames := make([]string, 0, len(cli.actions.Actions))

			i := 1
			for actionName := range cli.actions.Actions {
				fmt.Printf("%d. %s\n", i, actionName)
				actionNames = append(actionNames, actionName)
				i++
			}

			if len(actionNames) == 0 {
				fmt.Println("No actions available")
				return
			}

			var choice string

			fmt.Print("Choose action number: ")

			_, err := fmt.Scanln(&choice)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}

			if choice == "exit" {
				return
			}

			if choice == "0" {
				break
			}

			var actionNumber int

			_, err = fmt.Sscanf(choice, "%d", &actionNumber)
			if err != nil {
				fmt.Println("Invalid input")
				continue
			}

			if actionNumber < 1 || actionNumber > len(actionNames) {
				fmt.Println("Invalid action number")
				continue
			}

			actionName := actionNames[actionNumber-1]
			action := cli.actions.Actions[actionName]

			fmt.Println("Selected action:", actionName)

			if action.Executor == nil {
				fmt.Println("Action is not executable. Restart the application first.")
				continue
			}

			// Execute action here
			fmt.Println("Executing:", actionName)
			var res=tester.ExecuteAction(action)
			tester.Sequence_state.State[action.Name+"_"+strconv.Itoa(counter)] = res

			tester.Prev_res=res

			tester.ShowState()

			counter++;
		}

	case "sequence":
		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

		sequence, exists := cli.sequences.Sequences[sequenceName]
		if !exists {
			fmt.Println("Sequence does not exist:", sequenceName)
			return
		}

		fmt.Println("Executing sequence:", sequenceName)
		var tester=tester.CreateTesterObj()


		for i, action := range sequence.Actions_to_execute {
			fmt.Printf("Executing action %d: %s\n", i+1, action.Name)

			if action.Executor == nil {
				fmt.Println("Action is not executable:", action.Name)
				fmt.Println("Restart the application to load the action")
				return
			}

			// execute action here
			var res=tester.ExecuteAction(action)
			tester.Sequence_state.State[action.Name+"_"+strconv.Itoa(i)] = res

			tester.Prev_res=res

		}

		tester.ShowState()

	default:
		fmt.Println("Unknown execute type:", parts[0])
	}
}