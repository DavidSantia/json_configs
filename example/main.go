package main

import (
	"fmt"
	"os"

	"github.com/DavidSantia/json_configs"
)

func main() {
	var err error

	data := json_configs.NewData()

	err = data.GetCommandline()
	if err != nil {
		fmt.Printf("Command error: %v\n", err)
		os.Exit(1)
	}

	err = data.ReadConfigs()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	for name, device := range data.Configs {
		fmt.Printf("Device %s: %+v\n", name, device)
	}
	os.Exit(0)
}
