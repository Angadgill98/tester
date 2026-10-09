package Sequences

import (
	"encoding/json"
	"engine/Actions"
	"engine/Global"
	"os"
)

type SequenceObj map[string]Sequence

type Sequence struct {
	Name string                `json:"name"`
	Http []Actions.Http_action `json:"http"`
	Ws   []Actions.WsAction    `json:"ws"`
}

func CreateSequenceObj() (*SequenceObj,error) {
	sequenceObj := make(SequenceObj)
	err:=sequenceObj.SetUpSequenceObj()
	return &sequenceObj,err
}

func (sequenceObj *SequenceObj) SetUpSequenceObj() error {
	data, err := os.ReadFile(Global.GlobalConfig.SequenceConfigPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, sequenceObj); err != nil {
		return err
	}

	return nil
}

