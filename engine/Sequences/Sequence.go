package Sequences

import (
	"encoding/json"
	"fmt"
	"os"

	"engine/Global"

	"github.com/google/uuid"
)

var sequence_config_path = Global.GlobalConfig.SequenceConfigPath

type Sequence_config map[uuid.UUID]Sequence

type Sequence struct {
    ID   uuid.UUID   `json:"id"`
    Name string      `json:"name"`
    Http []uuid.UUID `json:"http"`
    Ws   []uuid.UUID `json:"ws"`
}

func CreateSequence(name string) *Sequence {
    return &Sequence{
        ID:   uuid.New(),
        Name: name,
        Http: []uuid.UUID{},
        Ws:   []uuid.UUID{},
    }
}

func CreateSequenceConfig() (*Sequence_config) {
	sequenceObj := make(Sequence_config)
	return &sequenceObj
}

func (sequenceObj *Sequence_config) SetUpSequenceConfig() error {
	data, err := os.ReadFile(sequence_config_path)
	if err != nil {
		return fmt.Errorf("failed to read sequence config: %w", err)
	}

	if err := json.Unmarshal(data, sequenceObj); err != nil {
		return fmt.Errorf("failed to parse sequence config: %w", err)
	}

	return nil
}

func (sequenceObj *Sequence_config) SaveSequenceConfig() error {
	if err := os.MkdirAll("Sequences", 0755); err != nil {
		return fmt.Errorf("failed to create sequence config directory: %w", err)
	}

	data, err := json.MarshalIndent(sequenceObj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to encode sequence config: %w", err)
	}

	if err := os.WriteFile(sequence_config_path, data, 0644); err != nil {
		return fmt.Errorf("failed to save sequence config: %w", err)
	}

	return nil
}