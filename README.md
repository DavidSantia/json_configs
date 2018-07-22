# json_configs
Flexible JSON config file example, using Go language reflect package

### Purpose
The package demonstrates how you can configure an application using one or more JSON files. 
For example, we are configuring a Fan device as shown below.

*fan.json*
```json
{
  "name": "Fan",
  "device_id": "A2",
  "onValue": "low",
  "offValue": "off",
  "deviceType": "fanlinc"
}
```
*credentials.json*
```json
[
  {
    "name": "Fan",
    "accessory": "Insteon",
    "host": "192.168.0.10",
    "port": "21000",
    "username": "home1",
    "password": "welcome1"
  },
  {
    "name": "Lamp",
    "accessory": "Insteon",
    "host": "192.168.0.10",
    "port": "21000",
    "username": "home1",
    "password": "welcome1"
  }
]
```
Notice the first file has generic device information for a Fan, in a single JSON element.
The second file contains an an array of two elements, configuring the network details for the Fan and another device.

Both formats (single elements and arrays) are handled to provide flexibility for configuring settings.

### Example Code
A sample program is provided in [example/config_device.go](https://github.com/DavidSantia/json_configs/blob/master/example/config_device.go)

* Reads the command-line with *getCommandline()* to set the directory to read
* It lists all *.json files in that directory
* It then opens and parses those files uisng *ReadConfigFiles()*

```go
package main

import (
	"fmt"
	"github.com/DavidSantia/json_configs"
	"path/filepath"
)

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

func main() {
	var ConfigDir string = "../config"
	var device Device

	// Get files
	fmt.Printf("== Reading config files in %s ==\n", ConfigDir)
	filenames, err := filepath.Glob(ConfigDir + "/*.json")
	if err != nil {
		fmt.Printf("reading config directory: %v\n", err)
		return
	}
	if len(filenames) == 0 {
		fmt.Println("no config files found")
		return
	}

	// Parse files into resultMap, using field "Name" as map key and device as each element
	resultMap, err := json_configs.ReadConfigFiles(&device, "Name", filenames...)
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
	}

	// Display results
	for name, v := range resultMap {
		fmt.Printf("Device %s: %+v\n", name, v)
	}
}
```

### Running the Example
Go into the *example* subdirectory, build the executable, and run as follows:
```sh
cd example
go build config_device.go
./config_device -h
Usage of ./config_device:
  -c string
    	Directory of JSON configs
```
Use the *config* directory to see how *Fan.json* and *Credentials.json* are combined to form the Fan settings:
```
./config_device -c ../config
== Reading config files in ../config ==
No errors
== Results 2 elements ==
Device Fan: {Name:Fan Accessory:Insteon DeviceType:fanlinc DeviceID:A2 Host:192.168.0.10 OffValue:off OnValue:low Password:welcome1 Port:21000 Username:home1}
Device Lamp: {Name:Lamp Accessory:Insteon DeviceType:lightBulb DeviceID:C12 Host:192.168.0.10 OffValue: OnValue: Password:welcome1 Port:21000 Username:home1}
```
Settings are also checked for consistency across multiple files.  Specify the *config_err* directory to see this:
```
== Reading config files in ../config_err ==
Config error: multiple errors
(#1) required id parameter Name not found, skipping [../config_err/fan_nameless.json]
(#2) invalid character '}' looking for beginning of object key string, skipping [../config_err/fan_bad.json]
(#3) settings for Fan conflict, parameter DeviceID: "A2" [fan.json:elem#1] != "A3" [fan_err.json] != "A0" [fan_extra.json]
(#4) unused setting for Fan, parameter color [fan_extra.json]
== Results 1 elements ==
Device Fan: {Name:Fan Accessory:Insteon DeviceType:fanlinc DeviceID:A0 Host:192.168.0.10 OffValue:off OnValue:low Password:welcome1 Port:21000 Username:home1}
```

### Validating Filenames
A second sample program [example/distinct_files.go](https://github.com/DavidSantia/json_configs/blob/master/example/distinct_files.go)
is provided to illustrate detailed error checking on a list of filenames.  It uses the function *DistinctFilenames* to
make sure files are specified only once, are accesible, and are valid files.
This function also constructs distinct filenames for messaging.  To see this, run the second example.
```sh
== Input file names ==
1. ../config/credentials.json
2. ../config/fan.json	
3. ../config/lamp.json
4. ../config_err/fan.json
5. ../config
6. ../config/nada.json
7. ../../json_configs/config/fan.json
== Results ==
multiple errors
(#1) invalid, no such file [../config/fan.json	]
(#2) invalid, filename is a directory [../config]
(#3) invalid, no such file [../config/nada.json]
== Output 4 items ==
• Filename: ../config/lamp.json
   Distict: lamp.json
• Filename: ../config_err/fan.json
   Distict: config_err/fan.json
• Filename: ../../json_configs/config/fan.json
   Distict: config/fan.json
• Filename: ../config/credentials.json
   Distict: credentials.json
```

### Customizing
The Device struct in the first sample program is just for example.
You can specify any struct for whatever you want to configure in your application.
* The functions *ReadConfigFile* and *ReadConfigFiles* have an interface as the first argument,
so you can pass a pointer to any struct.
* If you are using *ReadConfigFiles*, also specify the name of the field that will
serve as the Id for each element cofigured.
