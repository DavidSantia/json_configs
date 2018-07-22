package main

import (
	"flag"
	"fmt"
	"github.com/DavidSantia/json_configs"
	"os"
	"path/filepath"
)

// Example configuration
type Device struct {
	Name       string `json:"name"`
	Accessory  string `json:"accessory"`
	DeviceType string `json:"deviceType"`
	DeviceID   string `json:"device_id"`
	Host       string `json:"host"`
	OffValue   string `json:"offValue"`
	OnValue    string `json:"onValue"`
	Password   string `json:"password"`
	Port       string `json:"port"`
	Username   string `json:"username"`
}

var ConfigDir string
var ConfigMap = make(map[string]Device)

func main() {
	var err error
	var device Device
	var filenames []string
	var resultMap json_configs.ResultMap

	err = getCommandline()
	if err != nil {
		fmt.Printf("Command error: %v\n", err)
		os.Exit(1)
	}

	// Get files
	fmt.Printf("== Reading config files in %s ==\n", ConfigDir)
	if filenames, err = filepath.Glob(ConfigDir + "/*.json"); err != nil {
		fmt.Printf("Reading directory: %v\n", err)
		return
	}
	if len(filenames) == 0 {
		fmt.Println("No config files found")
		return
	}

	// Parse files into resultMap, using field "Name" as map key and device as each element
	resultMap, err = json_configs.ReadConfigFiles(&device, "Name", filenames...)
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
	} else {
		fmt.Println("No errors")
	}

	// Display results
	fmt.Printf("== Results %d elements ==\n", len(resultMap))
	for name, v := range resultMap {
		fmt.Printf("Device %s: %+v\n", name, v)
	}
	os.Exit(0)
}

func getCommandline() (err error) {
	var dirInfo os.FileInfo

	// Command-line arguments
	flag.StringVar(&ConfigDir, "c", "", "Directory of JSON configs")
	flag.Parse()

	// Validate
	if len(ConfigDir) == 0 {
		return fmt.Errorf("option -c <directory> required")
	}
	dirInfo, err = os.Stat(ConfigDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("-c %s %v", ConfigDir, err)
	} else if !dirInfo.IsDir() {
		return fmt.Errorf("-c %s is not a directory", ConfigDir)
	}
	return
}
