package json_configs

func NewData() (data *ExampleData) {
	data = &ExampleData{}
	data.Configs = make(ConfigMap)
	return
}

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

type ConfigMap map[string]Device

type ExampleData struct {
	ConfigDir string
	Configs   ConfigMap
}
