package services

import (
	"fmt"
	"os"
)



type Services struct {
	File_system File_system
}
func CreateServices()Services{
	return Services{
		File_system:File_system{},
	}
}









type File_system struct{}
var file_path="data/"
// var file_name="data/storage.log"

func (fs *File_system) SetUpFiles() error {
    fmt.Println("Skipping Setting up Sequence Files")

    // fmt.Println("Setting up Sequence Files")

    // err := fs.SetUpSequnceFiles()
    // if err != nil {
    //     return err
    // }

    fmt.Println("Setting up Action Files")

    err := fs.SetupActionFiles()
    if err != nil {
        return err
    }

    fmt.Println("Setting up Models File")

    err = fs.CreateModelsFile()
    if err != nil {
        return err
    }


    return nil
}

func (fs *File_system) SetupActionFiles() error {
	actionsPath := file_path + "/actions"

	err := os.MkdirAll(actionsPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("Actions directory setup successfully")

	tmpFilePath := actionsPath + "/tmp.go"

	code := `package data_actions

var MockExecutor = true
`

	err = os.WriteFile(tmpFilePath, []byte(code), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Temporary action file created successfully")

	return nil
}

func (fs *File_system) SetUpSequnceFiles() error {
	sequencesPath := file_path + "/sequences"

	err := os.MkdirAll(sequencesPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("Sequences directory setup successfully")

	tmpFilePath := sequencesPath + "/tmp.go"

	code := `package data_sequences

var MockExecutor = true
`

	err = os.WriteFile(tmpFilePath, []byte(code), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Temporary sequence file created successfully")

	return nil
}

func (fs *File_system) CreateModelsFile() error {
    modelsPath := file_path + "models.go"

    modelsCode := `package data

import (
	"net/http"
)

type SequnceState struct {
	State map[string]ActionResponse
	CustomData map[string]any
}

type ActionResponse struct {
	Body       string
	StatusCode int
	Status     string
	Headers    http.Header
	Cookies    []*http.Cookie
	Data       map[string]any
}
`

    file, err := os.Create(modelsPath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(modelsCode)
    if err != nil {
        return err
    }

    fmt.Println("Created:", modelsPath)

    return nil
}