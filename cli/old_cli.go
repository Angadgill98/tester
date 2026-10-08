package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"

// 	Actions "tester/Actions"
// 	Sequences "tester/Sequences"
// 	tester "tester/Tester"
// )

// type CLI struct {
// 	actions   *Actions.Actions_obj
// 	sequences *Sequences.Sequences_obj
// }

// func CreateCLI(actions *Actions.Actions_obj, sequences *Sequences.Sequences_obj) *CLI {
// 	return &CLI{
// 		actions:   actions,
// 		sequences: sequences,
// 	}
// }

// func (cli *CLI) Init() {
// 	fmt.Println("========== CLI START ==========")
// 	fmt.Println("Available commands:")
// 	fmt.Println("  create action <action_name>")
// 	fmt.Println("  create sequence <sequence_name>")
// 	fmt.Println("  delete action <action_name>")
// 	fmt.Println("  delete sequence <sequence_name>")
// 	fmt.Println("  sequence add-action <sequence_name>")
// 	fmt.Println("  execute action")
// 	fmt.Println("  execute sequence <sequence_name>")
// 	fmt.Println("  show")
// 	fmt.Println("  exit")
// 	fmt.Println()

// 	scanner := bufio.NewScanner(os.Stdin)

// 	for {
// 		fmt.Print("> ")

// 		if !scanner.Scan() {
// 			break
// 		}

// 		input := strings.TrimSpace(scanner.Text())

// 		if input == "" {
// 			continue
// 		}

// 		if input == "exit" {
// 			break
// 		}

// 		cli.ProcessCommand(input)
// 	}

// 	fmt.Println("========== CLI END ==========")
// }

// func (cli *CLI) ProcessCommand(input string) {
// 	parts := strings.Fields(input)

// 	if len(parts) == 0 {
// 		return
// 	}

// 	command := parts[0]
// 	remaining := strings.TrimSpace(strings.TrimPrefix(input, command))

// 	switch command {
// 	case "create":
// 		cli.ProcessCreate(remaining)

// 	case "delete":
// 		cli.ProcessDelete(remaining)

// 	case "sequence":
// 		cli.ProcessSequence(remaining)

// 	case "execute":
// 		cli.ProcessExecute(remaining)

// 	case "show":
// 		cli.ShowState()

// 	default:
// 		fmt.Println("Unknown command:", command)
// 	}
// }

// func (cli *CLI) ProcessCreate(input string) {
// 	parts := strings.Fields(input)

// 	if len(parts) == 0 {
// 		fmt.Println("Create what?")
// 		return
// 	}

// 	switch parts[0] {
// 	case "action":
// 		actionName := strings.TrimSpace(strings.TrimPrefix(input, "action"))

// 		if actionName == "" {
// 			fmt.Println("Usage: create action <action_name>")
// 			return
// 		}

// 		actionName = strings.ReplaceAll(actionName, " ", "_")

// 		var actionType string

// 		scanner := bufio.NewScanner(os.Stdin)

// 		validTypes := map[string]bool{
// 			"http": true,
// 			"ws":   true,
// 		}

// 		for {
// 			fmt.Print("Enter action type (http/ws or 'exit'): ")

// 			if !scanner.Scan() {
// 				return
// 			}

// 			actionType = strings.TrimSpace(scanner.Text())

// 			if actionType == "exit" {
// 				return
// 			}

// 			if validTypes[actionType] {
// 				break
// 			}

// 			fmt.Println("Invalid action type. Please try again.")
// 		}

// 		var opts Actions.Action_opts

// 		if actionType == "ws" {
// 			opts.Action_type = Actions.ActionTypeWS

// 			validWSEvents := map[string]bool{
// 				"send":    true,
// 				"connect": true,
// 				"recieve": true,
// 			}

// 			var wsEvent string

// 			for {
// 				fmt.Print("Enter WebSocket event (Send/Connect/Recieve): ")

// 				if !scanner.Scan() {
// 					return
// 				}

// 				wsEvent = strings.TrimSpace(scanner.Text())

// 				if validWSEvents[wsEvent] {
// 					break
// 				}

// 				fmt.Println("Invalid WebSocket event. Please try again.")
// 			}

// 			opts.Opts = Actions.WS_opts{
// 				Executor_type: wsEvent,
// 			}

	

// 		} else {
// 			opts.Action_type = Actions.ActionTypeHTTP

// 			opts.Opts = Actions.HTTP_opts{
// 				Executor_type: "req",
// 			}
// 		}

// 		created := cli.actions.CreateAction(actionName, opts)

// 		if created {
// 			fmt.Println("Action created successfully:", actionName)
// 		} else {
// 			fmt.Println("Action already exists")
// 		}

// 	case "sequence":
// 		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

// 		if sequenceName == "" {
// 			fmt.Println("Usage: create sequence <sequence_name>")
// 			return
// 		}

// 		sequenceName = strings.ReplaceAll(sequenceName, " ", "_")

// 		created := cli.sequences.CreateSequence(sequenceName)

// 		if created {
// 			fmt.Println("Sequence created successfully:", sequenceName)
// 		} else {
// 			fmt.Println("Sequence already exists")
// 		}

// 	default:
// 		fmt.Println("Unknown create command:", parts[0])
// 	}
// }

// func (cli *CLI) ProcessDelete(input string) {
// 	parts := strings.Fields(input)

// 	if len(parts) == 0 {
// 		fmt.Println("Usage: delete action <action_name>")
// 		fmt.Println("       delete sequence <sequence_name>")
// 		return
// 	}

// 	switch parts[0] {
// 	case "action":
// 		actionName := strings.TrimSpace(strings.TrimPrefix(input, "action"))

// 		if actionName == "" {
// 			fmt.Println("Usage: delete action <action_name>")
// 			return
// 		}

// 		actionName = strings.ReplaceAll(actionName, " ", "_")

// 		deleted := false

// 		if _, exists := cli.actions.Http.Http_actions[actionName]; exists {
// 			deleted = cli.actions.DeleteHttpAction(actionName)
// 		} else if _, exists := cli.actions.WS.WS_actions[actionName]; exists {
// 			deleted = cli.actions.DeleteWsAction(actionName)
// 		} else {
// 			fmt.Println("Action does not exist:", actionName)
// 			return
// 		}

// 		if deleted {
// 			fmt.Println("Action deleted successfully:", actionName)
// 		} else {
// 			fmt.Println("Failed to delete action:", actionName)
// 		}

// 	case "sequence":
// 		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

// 		if sequenceName == "" {
// 			fmt.Println("Usage: delete sequence <sequence_name>")
// 			return
// 		}

// 		sequenceName = strings.ReplaceAll(sequenceName, " ", "_")

// 		deleted := cli.sequences.DeleteSequence(sequenceName)

// 		if deleted {
// 			fmt.Println("Sequence deleted successfully:", sequenceName)
// 		} else {
// 			fmt.Println("Sequence does not exist:", sequenceName)
// 		}

// 	default:
// 		fmt.Println("Unknown delete command:", parts[0])
// 	}
// }

// func (cli *CLI) ProcessSequence(input string) {
// 	parts := strings.Fields(input)

// 	if len(parts) == 0 {
// 		fmt.Println("Usage: sequence add-action <sequence_name>")
// 		return
// 	}

// 	switch parts[0] {
// 	case "add-action":
// 		if len(parts) < 2 {
// 			fmt.Println("Usage: sequence add-action <sequence_name>")
// 			return
// 		}

// 		sequenceName := parts[1]

// 		_, exists := cli.sequences.Sequences[sequenceName]

// 		if !exists {
// 			fmt.Println("Sequence does not exist:", sequenceName)
// 			return
// 		}

// 		for {
// 			fmt.Println()
// 			fmt.Println("Available Actions:")
// 			fmt.Println("exit. Exit")
// 			fmt.Println("0. Back")

// 			type ActionEntry struct {
// 				Name      string
// 				ActionType Actions.ActionType
// 				Action     any
// 			}

// 			var availableActions []ActionEntry

// 			i := 1

// 			for actionName, action := range cli.actions.Http.Http_actions {
// 				fmt.Printf("%d. %s [HTTP]\n", i, actionName)

// 				availableActions = append(availableActions, ActionEntry{
// 					Name:       actionName,
// 					ActionType: Actions.ActionTypeHTTP,
// 					Action:     action,
// 				})

// 				i++
// 			}

// 			for actionName, action := range cli.actions.WS.WS_actions {
// 				fmt.Printf("%d. %s [WS]\n", i, actionName)

// 				availableActions = append(availableActions, ActionEntry{
// 					Name:       actionName,
// 					ActionType: Actions.ActionTypeWS,
// 					Action:     action,
// 				})

// 				i++
// 			}

// 			if len(availableActions) == 0 {
// 				fmt.Println("No actions available")
// 				return
// 			}

// 			var choice string

// 			fmt.Print("Choose action number: ")

// 			_, err := fmt.Scanln(&choice)

// 			if err != nil {
// 				fmt.Println("Invalid input")
// 				continue
// 			}

// 			if choice == "exit" {
// 				return
// 			}

// 			if choice == "0" {
// 				break
// 			}

// 			actionNumber, err := strconv.Atoi(choice)

// 			if err != nil {
// 				fmt.Println("Invalid input")
// 				continue
// 			}

// 			if actionNumber < 1 || actionNumber > len(availableActions) {
// 				fmt.Println("Invalid action number")
// 				continue
// 			}

// 			selected := availableActions[actionNumber-1]

// 			fmt.Println("Selected action:", selected.Name)
// 			fmt.Println("Adding it to the sequence:", sequenceName)

// 			added := cli.sequences.AddActionToSequence(sequenceName, selected.Action)

// 			if added {
// 				fmt.Println("Action added successfully:", selected.Name)
// 			} else {
// 				fmt.Println("Failed to add action:", selected.Name)
// 			}
// 		}

// 	default:
// 		fmt.Println("Unknown sequence command:", parts[0])
// 	}
// }

// func (cli *CLI) ShowState() {
// 	fmt.Println()
// 	fmt.Println("========== CURRENT STATE ==========")

// 	fmt.Println()
// 	fmt.Println("---------- HTTP ACTIONS ----------")
// 	fmt.Printf("Total HTTP Actions: %d\n", len(cli.actions.Http.Http_actions))

// 	for actionName, action := range cli.actions.Http.Http_actions {
// 		fmt.Println("Action:", actionName)
// 		fmt.Println("  Name:", action.Name)
// 		fmt.Println("  Path:", action.Path)
// 		fmt.Println("  Executor Type:", action.Executor_type)
// 		fmt.Println("  Executor Function:", action.Executor_function_name)
// 	}

// 	fmt.Println()
// 	fmt.Println("---------- HTTP ACTION CONFIG ----------")
// 	fmt.Printf("Total Config Entries: %d\n", len(cli.actions.Http.Config["http"]))

// 	for actionName, executor := range cli.actions.Http.Config["http"] {
// 		if executor != nil {
// 			fmt.Printf("Config: %s -> function installed\n", actionName)
// 		} else {
// 			fmt.Printf("Config: %s -> function NOT installed, need to restart\n", actionName)
// 		}
// 	}

// 	fmt.Println()
// 	fmt.Println("---------- WS ACTIONS ----------")
// 	fmt.Printf("Total WS Actions: %d\n", len(cli.actions.WS.WS_actions))

// 	for actionName, action := range cli.actions.WS.WS_actions {
// 		fmt.Println("Action:", actionName)
// 		fmt.Println("  Name:", action.Name)
// 		fmt.Println("  Path:", action.Path)
// 		fmt.Println("  Executor Type:", action.Exector_type)
// 		fmt.Println("  Executor Function:", action.Executor_function_name)
// 	}

// 	fmt.Println()
// 	fmt.Println("---------- WS ACTION CONFIG ----------")
// 	fmt.Printf("Total Config Entries: %d\n", len(cli.actions.WS.Config["ws"]))

// 	for actionName, executor := range cli.actions.WS.Config["ws"] {
// 		if executor != nil {
// 			fmt.Printf("Config: %s -> function installed\n", actionName)
// 		} else {
// 			fmt.Printf("Config: %s -> function NOT installed, need to restart\n", actionName)
// 		}
// 	}

// 	fmt.Println()
// 	fmt.Println("---------- SEQUENCES ----------")
// 	fmt.Printf("Total Sequences: %d\n", len(cli.sequences.Sequences))

// 	for sequenceName, sequenceActions := range cli.sequences.Sequences {
// 		fmt.Println("Sequence:", sequenceName)
// 		fmt.Printf("  Actions: %d\n", len(sequenceActions))

// 		for i, sequenceAction := range sequenceActions {
// 			switch sequenceAction.Action_type {
// 			case Actions.ActionTypeHTTP:
// 				if sequenceAction.HTTP_action != nil {
// 					fmt.Printf(
// 						"    %d. %s [HTTP]\n",
// 						i+1,
// 						sequenceAction.HTTP_action.Name,
// 					)
// 				}

// 			case Actions.ActionTypeWS:
// 				if sequenceAction.WS_action != nil {
// 					fmt.Printf(
// 						"    %d. %s [WS]\n",
// 						i+1,
// 						sequenceAction.WS_action.Name,
// 					)
// 				}
// 			}
// 		}
// 	}

// 	fmt.Println()
// 	fmt.Println("---------- SEQUENCE CONFIG ----------")
// 	fmt.Printf("Total Config Entries: %d\n", len(cli.sequences.Config))

// 	for sequenceName, sequenceActions := range cli.sequences.Config {
// 		fmt.Println("Sequence:", sequenceName)

// 		for i, sequenceAction := range sequenceActions {
// 			switch sequenceAction.Action_type {
// 			case Actions.ActionTypeHTTP:
// 				if sequenceAction.HTTP_action != nil {
// 					fmt.Printf(
// 						"  %d. %s [HTTP]\n",
// 						i+1,
// 						sequenceAction.HTTP_action.Name,
// 					)
// 				}

// 			case Actions.ActionTypeWS:
// 				if sequenceAction.WS_action != nil {
// 					fmt.Printf(
// 						"  %d. %s [WS]\n",
// 						i+1,
// 						sequenceAction.WS_action.Name,
// 					)
// 				}
// 			}
// 		}
// 	}

// 	fmt.Println()
// 	fmt.Println("===================================")
// 	fmt.Println()
// }

// func (cli *CLI) ProcessExecute(input string) {
// 	parts := strings.Fields(input)

// 	if len(parts) < 1 {
// 		fmt.Println("Usage: execute action")
// 		fmt.Println("       execute sequence <sequence_name>")
// 		return
// 	}

// 	switch parts[0] {
// 	case "action":
// 		fmt.Println("Starting Temp Sequence")

// 		tester := tester.CreateTesterObj()
// 		counter := 0

// 		for {
// 			fmt.Println()
// 			fmt.Println("========== EXECUTE ACTION ==========")
// 			fmt.Println("Available Actions:")
// 			fmt.Println("exit. Exit")
// 			fmt.Println("0. Back")

// 			type ActionEntry struct {
// 				Name   string
// 				Action any
// 			}

// 			var availableActions []ActionEntry

// 			i := 1

// 			for actionName, action := range cli.actions.Http.Http_actions {
// 				fmt.Printf("%d. %s [HTTP]\n", i, actionName)

// 				availableActions = append(availableActions, ActionEntry{
// 					Name:   actionName,
// 					Action: action,
// 				})

// 				i++
// 			}

// 			for actionName, action := range cli.actions.WS.WS_actions {
// 				fmt.Printf("%d. %s [WS]\n", i, actionName)

// 				availableActions = append(availableActions, ActionEntry{
// 					Name:   actionName,
// 					Action: action,
// 				})

// 				i++
// 			}

// 			if len(availableActions) == 0 {
// 				fmt.Println("No actions available")
// 				return
// 			}

// 			var choice string

// 			fmt.Print("Choose action number: ")

// 			_, err := fmt.Scanln(&choice)

// 			if err != nil {
// 				fmt.Println("Invalid input")
// 				continue
// 			}

// 			if choice == "exit" {
// 				return
// 			}

// 			if choice == "0" {
// 				break
// 			}

// 			actionNumber, err := strconv.Atoi(choice)

// 			if err != nil {
// 				fmt.Println("Invalid input")
// 				continue
// 			}

// 			if actionNumber < 1 || actionNumber > len(availableActions) {
// 				fmt.Println("Invalid action number")
// 				continue
// 			}

// 			selected := availableActions[actionNumber-1]

// 			fmt.Println("Selected action:", selected.Name)

// 			switch action := selected.Action.(type) {
// 			case Actions.Http_actions:
// 				fmt.Println("HTTP action selected:", action.Name)

// 				var res = tester.ExecuteHTTPAction(*cli.actions, action)

// 				fmt.Println("HTTP response received for action:", action.Name)
// 				fmt.Println("Status Code:", res.StatusCode)
// 				fmt.Println("Status:", res.Status)
// 				fmt.Println("Headers:", res.Headers)
// 				fmt.Println("Cookies:", res.Cookies)
// 				fmt.Println("Body:", res.Body)

// 				tester.Prev_res = *res
// 				tester.Sequence_state.State[action.Name] = *res
				
// 			case Actions.WsAction:
// 				fmt.Println("WS action selected:", action.Name)
// 				fmt.Println("WebSocket execution is not implemented yet.")
// 			}

// 			counter++
// 		}

// 	case "sequence":
// 		sequenceName := strings.TrimSpace(strings.TrimPrefix(input, "sequence"))

// 		sequence, exists := cli.sequences.Sequences[sequenceName]

// 		if !exists {
// 			fmt.Println("Sequence does not exist:", sequenceName)
// 			return
// 		}

// 		fmt.Println("Executing sequence:", sequenceName)

// 		actionTester := tester.CreateTesterObj()

// 		for i, sequenceAction := range sequence {
// 			switch sequenceAction.Action_type {
// 			case Actions.ActionTypeHTTP:
// 				if sequenceAction.HTTP_action == nil {
// 					fmt.Println("HTTP action is nil")
// 					return
// 				}

// 				action := sequenceAction.HTTP_action

// 				fmt.Printf(
// 					"Executing action %d: %s [HTTP]\n",
// 					i+1,
// 					action.Name,
// 				)

// 				/*
// 					Update this when Tester.ExecuteAction
// 					accepts the new Http_actions object.
// 				*/

// 			case Actions.ActionTypeWS:
// 				if sequenceAction.WS_action == nil {
// 					fmt.Println("WS action is nil")
// 					return
// 				}

// 				action := sequenceAction.WS_action

// 				fmt.Printf(
// 					"Executing action %d: %s [WS]\n",
// 					i+1,
// 					action.Name,
// 				)

// 				fmt.Println("WebSocket execution is not implemented yet.")
// 			}
// 		}

// 		actionTester.ShowState()

// 	default:
// 		fmt.Println("Unknown execute type:", parts[0])
// 	}
// }