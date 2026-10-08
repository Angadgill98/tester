package sequences

import (
	"fmt"
	"os"
	"strings"

	actions "tester/Actions"
)

type Sequences_obj struct {
	Sequences map[string][]SequenceAction
	Config    map[string][]SequenceAction
}

func CreateSequencesObj() Sequences_obj {
	return Sequences_obj{
		Sequences: Sequence_Config,
		Config:    Sequence_Config,
	}
}

type SequenceAction struct {
	Action_type actions.ActionType
	HTTP_action *actions.Http_actions
	WS_action   *actions.WsAction
}





func (this *Sequences_obj) SetupSequences(actions_obj actions.Actions_obj) {
	isok, err := this.CheckSequencesFS()

	if err != nil {
		return
	}

	if !isok {
		fmt.Println("Sequences FS not Setup")
		return
	}

	fmt.Println("Sequences FS Setup OK")


	fmt.Printf("Loaded %d sequences\n", len(this.Sequences))

	for sequenceName := range this.Sequences {
		fmt.Println("Sequence:", sequenceName)
	}
}

func (this *Sequences_obj) CheckSequencesFS() (bool, error) {
	_, err := os.Stat("Sequences/sequences_config.go")

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("sequences_config doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied for sequences_config")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	return true, nil
}



func (this *Sequences_obj) CreateSequence(sequenceName string) bool {
	if _, exists := this.Sequences[sequenceName]; exists {
		fmt.Println("Sequence already exists:", sequenceName)
		return false
	}

	this.Sequences[sequenceName] = []SequenceAction{}
	// this.Config[sequenceName] = []SequenceAction{}

	err := CreateCodeSequenceConfig(this.Sequences)

	if err != nil {
		fmt.Println("Failed to update sequences config:", err)
		return false
	}

	return true
}


var config_file_path = "Sequences/"


func CreateCodeSequenceConfig(sequences map[string][]SequenceAction) error {
	var builder strings.Builder

	builder.WriteString("package sequences\n\n")
	builder.WriteString("import (\n")
	builder.WriteString("\tactions \"tester/Actions\"\n")
	builder.WriteString(")\n\n")
	builder.WriteString("var Sequence_Config = map[string][]SequenceAction{\n")

	for sequenceName, sequenceActions := range sequences {
		builder.WriteString(fmt.Sprintf("\t%q: {\n", sequenceName))

		for _, sequenceAction := range sequenceActions {
			builder.WriteString("\t\t{\n")
			builder.WriteString(fmt.Sprintf("\t\t\tAction_type: actions.ActionType(%q),\n", sequenceAction.Action_type))

			if sequenceAction.HTTP_action != nil {
				action := sequenceAction.HTTP_action

				builder.WriteString("\t\t\tHTTP_action: &actions.Http_actions{\n")
				builder.WriteString(fmt.Sprintf("\t\t\t\tName: %q,\n", action.Name))
				builder.WriteString(fmt.Sprintf("\t\t\t\tPath: %q,\n", action.Path))
				builder.WriteString(fmt.Sprintf("\t\t\t\tExecutor_type: %q,\n", action.Executor_type))
				builder.WriteString(fmt.Sprintf("\t\t\t\tExecutor_function_name: %q,\n", action.Executor_function_name))
				builder.WriteString("\t\t\t},\n")
			}

			if sequenceAction.WS_action != nil {
				action := sequenceAction.WS_action

				builder.WriteString("\t\t\tWS_action: &actions.WsAction{\n")
				builder.WriteString(fmt.Sprintf("\t\t\t\tName: %q,\n", action.Name))
				builder.WriteString(fmt.Sprintf("\t\t\t\tPath: %q,\n", action.Path))
				builder.WriteString(fmt.Sprintf("\t\t\t\tExector_type: %q,\n", action.Exector_type))
				builder.WriteString(fmt.Sprintf("\t\t\t\tExecutor_function_name: %q,\n", action.Executor_function_name))
				builder.WriteString("\t\t\t},\n")
			}

			builder.WriteString("\t\t},\n")
		}

		builder.WriteString("\t},\n")
	}

	builder.WriteString("}\n")

	err := os.WriteFile(config_file_path+"sequences_config.go", []byte(builder.String()), 0644)

	if err != nil {
		return err
	}

	return nil
}




func (this *Sequences_obj) AddHTTPActionToSequence(sequenceName string, action actions.Http_actions) bool {
	_, sequenceExists := this.Sequences[sequenceName]

	if !sequenceExists {
		fmt.Println("Sequence does not exist:", sequenceName)
		return false
	}

	this.Sequences[sequenceName] = append(this.Sequences[sequenceName], SequenceAction{
		Action_type: actions.ActionTypeHTTP,
		HTTP_action: &action,
	})

	// this.Config[sequenceName] = append(this.Config[sequenceName], SequenceAction{
	// 	Action_type: actions.ActionTypeHTTP,
	// 	HTTP_action: &action,
	// })

	err := CreateCodeSequenceConfig(this.Sequences)

	if err != nil {
		fmt.Println("Failed to update sequence config:", err)
		return false
	}

	return true
}

func (this *Sequences_obj) AddWSActionToSequence(sequenceName string, action actions.WsAction) bool {
	_, sequenceExists := this.Sequences[sequenceName]

	if !sequenceExists {
		fmt.Println("Sequence does not exist:", sequenceName)
		return false
	}

	this.Sequences[sequenceName] = append(this.Sequences[sequenceName], SequenceAction{
		Action_type: actions.ActionTypeWS,
		WS_action:   &action,
	})

	// this.Config[sequenceName] = append(this.Config[sequenceName], SequenceAction{
	// 	Action_type: actions.ActionTypeWS,
	// 	WS_action:   &action,
	// })

	err := CreateCodeSequenceConfig(this.Sequences)

	if err != nil {
		fmt.Println("Failed to update sequence config:", err)
		return false
	}

	return true
}

func (this *Sequences_obj) AddActionToSequence(sequenceName string, action any) bool {
	switch typedAction := action.(type) {
	case actions.Http_actions:
		return this.AddHTTPActionToSequence(sequenceName, typedAction)

	case actions.WsAction:
		return this.AddWSActionToSequence(sequenceName, typedAction)

	default:
		fmt.Println("Invalid action type")
		return false
	}
}



func (this *Sequences_obj) DeleteSequence(sequenceName string) bool {
	_, exists := this.Sequences[sequenceName]

	if !exists {
		fmt.Println("Sequence does not exist:", sequenceName)
		return false
	}

	delete(this.Sequences, sequenceName)
	delete(this.Config, sequenceName)

	err := CreateCodeSequenceConfig(this.Sequences)

	if err != nil {
		fmt.Println("Failed to update sequences config:", err)
		return false
	}

	fmt.Println("Sequence deleted:", sequenceName)

	return true
}




