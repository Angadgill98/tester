package Application

import (
	"encoding/json"
	"engine/Global"
	"os"
)

type Application_obj map[string]*Application

func CreateApplicationObj() *Application_obj {
	applicationObj := make(Application_obj)
	return &applicationObj
}

type Application struct {
	Name     string                 `json:"name"`
	Children map[string]*Application `json:"children"`
}


func (*Application_obj)CreateApplication(name string) *Application {
	return &Application{
		Name:     name,
		Children: make(map[string]*Application),
	}
}



func (applicationObj *Application_obj) SetUpApplicationObj() {
	data, err := os.ReadFile(Global.GlobalConfig.ApplicationConfigPath)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, applicationObj); err != nil {
		return
	}
}


func (applicationObj *Application_obj) SaveApplicationObj() error {
	data, err := json.MarshalIndent(*applicationObj, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(Global.GlobalConfig.ApplicationConfigPath, data, 0644)
}