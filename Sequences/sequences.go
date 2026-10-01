package sequences

import (
	"fmt"
	"os"
	"strings"

	actions "tester/Actions"
)

type Sequences_obj struct {
	Sequences map[string]Sequence
	Config map[string][]string
}

func CreateSequencesObj() Sequences_obj {
	return Sequences_obj{
		Sequences: make(map[string]Sequence),
		Config:Sequence_Config,

	}
}

var SequenceFS_path = "data/sequences/"

type Sequence struct {
	Sequence_name     string
	Actions_to_execute []actions.Action
	path             string
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

	this.LoadSequences(actions_obj)

	fmt.Printf("Loaded %d sequences\n", len(this.Sequences))

	for sequenceName := range this.Sequences {
		fmt.Println("Sequence:", sequenceName)
	}
}

func (this *Sequences_obj) CheckSequencesFS() (bool, error) {
	
	_, err := os.Stat("Sequences/"+"sequences_config.go")
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


func (this *Sequences_obj) LoadSequences(actions_obj actions.Actions_obj) {
	for sequenceName, actionNames := range this.Config {
		var sequenceActions []actions.Action

		for _, actionName := range actionNames {
			action, exists := actions_obj.Actions[actionName]
			if !exists {
				fmt.Println("Action does not exist:", actionName)
				continue
			}

			sequenceActions = append(sequenceActions, action)
		}

		sequence := Sequence{
			Sequence_name:       sequenceName,
			Actions_to_execute: sequenceActions,
			path:                 SequenceFS_path,
		}

		this.Sequences[sequenceName] = sequence
	}
}



func (this *Sequences_obj) CreateSequence(SequenceName string) bool {
	if _, exists := this.Sequences[SequenceName]; exists {
		fmt.Println("Sequence already exists:", SequenceName)
		return false
	}

	Sequence := Sequence{
		Sequence_name:       SequenceName,
		Actions_to_execute: []actions.Action{},
		path:               SequenceFS_path,
	}

	

	err := this.UpdateConfig(SequenceName)
	if err != nil {
		fmt.Println("Failed to update sequences config:", err)
		return false
	}

	this.Sequences[SequenceName] = Sequence
	this.Config[SequenceName]=make([]string, 0)

	return true
}



var config_file_path="Sequences/"
func (this *Sequences_obj) UpdateConfig(sequenceName string) error {
	configPath := config_file_path + "/sequences_config.go"

	code, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	entry := fmt.Sprintf("\t\"%s\": []string{},\n", sequenceName)

	configCode := string(code)

	insertPosition := strings.LastIndex(configCode, "}")
	if insertPosition == -1 {
		return fmt.Errorf("invalid sequences_config.go")
	}

	configCode = configCode[:insertPosition] + entry + configCode[insertPosition:]

	err = os.WriteFile(configPath, []byte(configCode), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Sequences config updated:", sequenceName)

	return nil
}



func (this *Sequences_obj) AddActionToSequence(sequenceName string, action actions.Action) bool {
	sequence, sequenceExists := this.Sequences[sequenceName]

	if !sequenceExists {
		fmt.Println("Sequence does not exist:", sequenceName)
		return false
	}

	sequence.Actions_to_execute = append(sequence.Actions_to_execute, action)
	this.Sequences[sequenceName] = sequence

	err := this.UpdateSequenceConfig(sequenceName, action.Name)
	if err != nil {
		fmt.Println("Failed to update sequence config:", err)
		return false
	}

	return true
}

func (this *Sequences_obj) UpdateSequenceConfig(sequenceName string, actionName string) error {
	configPath := config_file_path + "/sequences_config.go"

	code, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	configCode := string(code)

	this.Config[sequenceName] = append(this.Config[sequenceName], actionName)

	var builder strings.Builder

	builder.WriteString("var Sequence_Config = map[string][]string{\n")

	for name, actionNames := range this.Config {
		builder.WriteString(fmt.Sprintf("\t%q: {\n", name))

		for _, actionName := range actionNames {
			builder.WriteString(fmt.Sprintf("\t\t%q,\n", actionName))
		}

		builder.WriteString("\t},\n")
	}

	builder.WriteString("}\n")

	varStart := strings.Index(configCode, "var Sequence_Config")
	if varStart == -1 {
		return fmt.Errorf("Sequence_Config not found")
	}

	newConfigCode := configCode[:varStart] + builder.String()

	err = os.WriteFile(configPath, []byte(newConfigCode), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Sequence config updated:", sequenceName, "->", actionName)

	return nil
}







func (this *Sequences_obj) DeleteSequence(sequenceName string) bool {
	_, exists := this.Sequences[sequenceName]

	if !exists {
		fmt.Println("Sequence does not exist:", sequenceName)
		return false
	}

	err := this.DeleteSequenceConfig(sequenceName)
	if err != nil {
		fmt.Println("Failed to remove sequence from config:", err)
		return false
	}

	delete(this.Sequences, sequenceName)
	delete(this.Config, sequenceName)

	fmt.Println("Sequence deleted:", sequenceName)

	return true
}

func (this *Sequences_obj) DeleteSequenceConfig(sequenceName string) error {
	configPath := config_file_path + "/sequences_config.go"

	code, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	configCode := string(code)

	entry := fmt.Sprintf("\t%q: {\n", sequenceName)

	start := strings.Index(configCode, entry)
	if start == -1 {
		return fmt.Errorf("sequence not found in config: %s", sequenceName)
	}

	end := strings.Index(configCode[start:], "\t},")
	if end == -1 {
		return fmt.Errorf("invalid sequence config entry: %s", sequenceName)
	}

	end = start + end + len("\t},")

	configCode = configCode[:start] + configCode[end:]

	err = os.WriteFile(configPath, []byte(configCode), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Sequence removed from config:", sequenceName)

	return nil
}