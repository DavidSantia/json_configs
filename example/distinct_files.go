package main

import (
	"fmt"
	"github.com/DavidSantia/json_configs"
)

func main() {
	var errList []string

	var filenames = []string{
		"../config/credentials.json",
		"../config/fan.json	",
		"../config/lamp.json",
		"../config_err/fan.json",
		"../config",
		"../config/nada.json",
		"../../json_configs/config/fan.json",
	}

	// Get files
	fmt.Println("== Input file names ==")
	for i, filename := range filenames {
		fmt.Printf("%d. %s\n", i+1, filename)
	}

	fileDetails := json_configs.DistinctFilenames(filenames, &errList)

	// Compile error list into an error message
	fmt.Println("== Results ==")
	if len(errList) > 0 {
		if len(errList) == 1 {
			fmt.Println(errList[0])
		} else {
			fmt.Println("multiple errors")
			for i := range errList {
				fmt.Printf("(#%d) %s\n", i+1, errList[i])
			}
		}
	} else {
		fmt.Println("no errors")
	}

	fmt.Printf("== Output %d items ==\n", len(fileDetails))

	for _, file := range fileDetails {
		fmt.Printf("• Filename: %s\n", file.Name)
		fmt.Printf("   Distict: %s\n", file.DistinctName)
	}
}
