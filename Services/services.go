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

	wsPath := actionsPath + "/ws"

	err = os.MkdirAll(wsPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("WebSocket actions directory setup successfully")

	sendPath := wsPath + "/send"

	err = os.MkdirAll(sendPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("WebSocket actions directory setup successfully")

	conPath := wsPath + "/connect"

	err = os.MkdirAll(conPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("WebSocket actions directory setup successfully")

	recPath := wsPath + "/recieve"

	err = os.MkdirAll(recPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("WebSocket actions directory setup successfully")

	httpPath := actionsPath + "/http"

	err = os.MkdirAll(httpPath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("http directory setup successfully")

	reqTypePath := httpPath + "/req"

	err = os.MkdirAll(reqTypePath, 0755)
	if err != nil {
		return err
	}

	fmt.Println("req directory setup successfully")

	err = createTmpGoFile(reqTypePath, "http_req")
	if err != nil {
		return err
	}

	err = createTmpGoFile(sendPath, "ws_actions_send")
	if err != nil {
		return err
	}

	err = createTmpGoFile(conPath, "ws_actions_connect")
	if err != nil {
		return err
	}

	err = createTmpGoFile(recPath, "ws_actions_receive")
	if err != nil {
		return err
	}

	return nil
}
func createTmpGoFile(path string, packageName string) error {
	tmpPath := path + "/tmp.go"

	code := fmt.Sprintf(`package %s

func MockFunction() {
}
`, packageName)

	err := os.WriteFile(tmpPath, []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating tmp.go:", err)
		return err
	}

	fmt.Println("Created:", tmpPath)

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