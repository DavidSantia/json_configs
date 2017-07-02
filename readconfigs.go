package json_configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

func (data *ExampleData) ReadConfigs() (err error) {
	var matches, matches_shortname []string
	var filename string
	var i, errCount int

	if matches, err = filepath.Glob(data.ConfigDir + "/*.json"); err != nil {
		return
	}
	if len(matches) == 0 {
		err = fmt.Errorf("No config files found")
		return
	}

	for _, filename = range matches {
		matches_shortname = append(matches_shortname, filepath.Base(filename))
	}
	fmt.Printf("Reading %d config files %v\n", len(matches), matches_shortname)

	var b []byte
	for i, filename = range matches {
		if b, err = ioutil.ReadFile(filename); err != nil {
			fmt.Printf("Issue reading config: %v\n", err)
			errCount++
			continue
		}
		configMap := make(ConfigMap)
		if err = json.Unmarshal(b, &configMap); err != nil {
			fmt.Printf("Issue parsing config: %v\n", err)
			errCount++
			continue
		}
		if err = data.insertMap(configMap, matches_shortname[i]); err != nil {
			fmt.Printf("Issue with config: %v\n", err)
			errCount++
			continue
		}
	}

	if errCount > 0 {
		err = fmt.Errorf("%d errors found reading configuration files", errCount)
	} else {
		err = nil
	}
	return
}

// Allow a place to be configured across multiple .json files
// Dis-allow conflicting settings across files

func (data *ExampleData) insertMap(configMap ConfigMap, filename string) (err error) {
	var curr Device
	var name, tag string

	// Check for possible duplicate config
	for name, curr = range configMap {

		prev, ok := data.Configs[name]

		// New item, create using this Config
		if !ok {
			curr.Name = name
			data.Configs[name] = curr
			continue
		}

		st := reflect.TypeOf(prev)
		prevV := reflect.ValueOf(&prev).Elem()
		currV := reflect.ValueOf(curr)

		// Iterate through Config fields, validating if tagged `json:`
		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			if tag, ok = field.Tag.Lookup("json"); ok {
				if tag != "" {
					prev := prevV.Field(i).String()
					curr := currV.Field(i).String()
					if len(curr) != 0 {
						if len(prev) != 0 {
							// If previously set, does it differ?
							if prev != curr {
								err = fmt.Errorf("%s in %s: value %q differs from previous %q for %s",
									field.Name, filename, curr, prev, name)
								return
							}
						} else {
							// New value, set
							prevV.Field(i).SetString(curr)
						}
					}
				}
			}
		}
		data.Configs[name] = prev
	}
	return
}

