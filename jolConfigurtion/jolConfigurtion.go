package jolConfigurtion

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Configurtion struct {
	Popup     string   `json:"popup"`
	Fields    []string `json:"fields"`
	UrlAddr   string   `json:"url"`
	Extension string   `json:"extension"`
	Agent     string   `json:"agent"`
	Token     string   `json:"token"`
	Debug     bool     `json:"debug"`
}

type App struct {
	Config Configurtion
	Wg     *sync.WaitGroup
}

func LoadConfigFile() (Configurtion, error) {

	appDir := getAppPath()
	fmt.Println(appDir)
	//configPath := filepath.Join(appDir, "config/config.json")
	configPath := "D:\\Work\\Dropbox\\GO\\src\\CloudAsteriskAMI\\AgentDesktop\\config\\config.json"
	var config Configurtion
	Configurtion, err := os.ReadFile(configPath)

	if err != nil {
		return config, err
	}
	err = json.Unmarshal(Configurtion, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func getAppPath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// مسیر پوشه برنامه
	appDir := filepath.Dir(exePath)
	return appDir
}
