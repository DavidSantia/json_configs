# json_configs
Flexible JSON config file example, using Go language reflect package

### Purpose
The package demonstrates how you can configure an application using one or more JSON files. 
For example, we are configuring a Fan device across two files.

*Fan.json*
```json
{
  "Fan": {
    "device_id": "A2",
    "onValue": "low",
    "offValue": "off",
    "deviceType": "fanlinc"
  }
}	
```

*Credentials.json*
```json
{
  "Fan": {
    "accessory": "Insteon",
    "host": "192.168.0.10",
    "port": "21000",
    "username": "home1",
    "password": "welcome1"
  }
}
```
### Example Code
A sample program is provided in [example/main.go](https://github.com/DavidSantia/json_configs/blob/master/example/main.go)
```
Usage of ./example:
  -c string
    	Directory of JSON configs
```
This code sample

* Reads the command-line using the *GetCommandline()* method
* This specifies a directory contianing JSON configuration files
* It then opens the JSON files, and then parses them uisng *ReadConfigs()*

```go
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
```

### Running the Example
Go into the *example* subdirectory, build the executable, and run as follows:
```sh
$ cd example
$ go build
````
Use the *config* directory to see how *Fan.json* and *Credentials.json* are combined for the Fan settings:
```
$ ./example -c ../config
Reading 3 config files [credentials.json fan.json lamp.json]
Device Fan: {Name:Fan Accessory:Insteon DeviceType:fanlinc DeviceID:A2 Host:192.168.0.10 OffValue:off OnValue:low Password:welcome1 Port:21000 Username:home1}
Device Lamp: {Name:Lamp Accessory:Insteon DeviceType:lightBulb DeviceID:C12 Host:192.168.0.10 OffValue: OnValue: Password:welcome1 Port:21000 Username:home1}
```
The reflect package is used to ensure parameters are consistent across multiple files.  Try the *config_err* directory to see this:
```
$ ./example -c ../config_err
Reading 2 config files [fan.json fan_err.json]
Issue with config: DeviceID in fan_err.json: value "A3" differs from previous "A2" for Fan
Config error: 1 errors found reading configuration files
```
### Customizing
The structs in the file *datatypes.go* are just for example.  You can modify the Device struct to represent whatever you want to configure for your application.

```go
// Example configuration
type Device struct {
	Name       string
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
```

