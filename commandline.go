package json_configs

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func (data *ExampleData) validateFlags() (err error) {
	var absolutePath string
	var dirInfo os.FileInfo

	// Validate PaymentConfig
	if len(data.ConfigDir) == 0 {
		return fmt.Errorf("option -c <directory> required")
	}
	if absolutePath, err = filepath.Abs(data.ConfigDir); err != nil {
		return fmt.Errorf("-c %s %v", data.ConfigDir, err)
	}
	dirInfo, err = os.Stat(absolutePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("-c %s %v", data.ConfigDir, err)
	} else if !dirInfo.IsDir() {
		return fmt.Errorf("-c %s is not a directory", data.ConfigDir)
	}
	data.ConfigDir = absolutePath

	return
}

func (data *ExampleData) GetCommandline() (err error) {

	// Command-line arguments
	flag.StringVar(&data.ConfigDir, "c", "", "Directory of JSON configs")

	// Parse commandline flag arguments
	flag.Parse()

	// Validate
	err = data.validateFlags()
	if err != nil {
		return
	}

	return
}
