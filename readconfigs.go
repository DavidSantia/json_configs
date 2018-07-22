package json_configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

var Debug bool

// Read a single config file, return a struct, where 'data' is a pointer to that struct
func ReadConfigFile(data interface{}, filename string) (err error) {
	var b []byte
	var errList []string
	var config interface{}
	var k string
	var i int

	// Make sure data is a pointer to a struct
	k = reflect.TypeOf(data).Kind().String()
	if k != "ptr" {
		err = fmt.Errorf("ReadConfigFile: 'data' must be ptr, not %s", k)
		panic(err)
	}
	st := reflect.TypeOf(data).Elem()
	sv := reflect.ValueOf(data).Elem()

	_, err = ValidateFile(filename)
	if err != nil {
		return
	}

	if b, err = ioutil.ReadFile(filename); err != nil {
		err = fmt.Errorf("reading file: %v", err)
		return
	}

	// Parse config file into map[string]interface{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		err = fmt.Errorf("%v [%s]", err, filename)
		return
	}

	// Each file can contain a single element of type 'data'
	k = reflect.TypeOf(config).Kind().String()
	if k != "map" {
		err = fmt.Errorf("contains type %q, must be a JSON element [%s]", k, filename)
		return
	}

	if Debug {
		log.Printf("Parsing single element [%s]", filename)
	}
	parsed := Parsed{
		FileName:     filename,
		DistinctName: filepath.Base(filename),
		ElementMap:   config.(map[string]interface{}),
	}

	// store in resultMap
	parsedMap := make(ParsedMap)
	parsedMap["default"] = []Parsed{parsed}

	// Parse dataMap entries into data object (st, sv) fields
	parseConfig(st, sv, parsedMap, &errList)

	// Compile error list into an error message
	if len(errList) > 0 {
		if len(errList) == 1 {
			err = fmt.Errorf(errList[0])
		} else {
			for i = range errList {
				errList[i] = fmt.Sprintf("(#%d) %s", i+1, errList[i])
			}
			err = fmt.Errorf("multiple errors\n%s", strings.Join(errList, "\n"))
		}
	}
	return
}

// Return a map[id] of structs ('data' is a pointer to that struct), read from directory of JSON config files
func ReadConfigFiles(data interface{}, idName string, filenames ...string) (resultMap ResultMap, err error) {
	var b []byte
	var errList []string
	var config, v interface{}
	var file FileDetail
	var fileDetails []FileDetail
	var parsedMap ParsedMap
	var parsedArr []Parsed
	var k, elementId, idTag string
	var i int
	var ok, tag bool

	// Make sure data is a pointer to a struct
	k = reflect.TypeOf(data).Kind().String()
	if k != "ptr" {
		err = fmt.Errorf("ReadConfigFiles: 'data' must be ptr, not %s", k)
		panic(err)
	}
	st := reflect.TypeOf(data).Elem()
	sv := reflect.ValueOf(data).Elem()

	// Locate field idName and `json:"idTag"`
	for i = 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.Name == idName {
			idTag, tag = field.Tag.Lookup("json")
			ok = true
			break
		}
	}
	if !ok {
		err = fmt.Errorf("ReadConfigFiles: 'data' does not contain field %s", idName)
		panic(err)
	}

	// Validate file list
	fileDetails = DistinctFilenames(filenames, &errList)

	// Each file can contain a single element of type 'data', or an array of these elements
	parsedMap = make(ParsedMap)
	for _, file = range fileDetails {
		if b, err = ioutil.ReadFile(file.Name); err != nil {
			errList = append(errList, fmt.Sprintf("reading file: %v", err))
			continue
		}
		config = nil

		// Parse config file as map[string]interface{}, or slice of these
		err = json.Unmarshal(b, &config)
		if err != nil {
			if Debug {
				log.Printf("Parsing issue, skipping [%s]", file.Name)
			}
			errList = append(errList, fmt.Sprintf("%v, skipping [%s]", err, file.Name))
			continue
		}

		k = reflect.TypeOf(config).Kind().String()
		if k == "map" {
			if Debug {
				log.Printf("Parsing single element [%s]", file.Name)
			}
			parsed := Parsed{
				FileName:     file.Name,
				DistinctName: file.DistinctName,
				ElementMap:   config.(map[string]interface{}),
			}

			// Find element Id by json tag or field name
			v, _ = parsed.ElementMap[idTag]
			if v == nil && tag {
				v, _ = parsed.ElementMap[idName]
			}
			if v == nil {
				errList = append(errList, fmt.Sprintf("required id parameter %s not found, skipping [%s]",
					idName, file.Name))
				continue
			}
			elementId = fmt.Sprintf("%v", v)

			// Add parsed to array and store in resultMap
			parsedArr, ok = parsedMap[elementId]
			if ok {
				parsedArr = append(parsedArr, parsed)
			} else {
				parsedArr = []Parsed{parsed}
			}
			parsedMap[elementId] = parsedArr

		} else if k == "slice" {
			if Debug {
				log.Printf("Parsing %d elements [%s]", len(config.([]interface{})), file.Name)
			}
			// Config file contains an array of elements of type 'data'
			for i, v = range config.([]interface{}) {

				parsed := Parsed{
					FileName:     file.Name,
					DistinctName: file.DistinctName,
					Position:     i + 1,
					ElementMap:   v.(map[string]interface{}),
				}

				// Find element Id by json tag or field name
				v, _ = parsed.ElementMap[idTag]
				if v == nil && tag {
					v, _ = parsed.ElementMap[idName]
				}
				if v == nil {
					errList = append(errList, fmt.Sprintf("required id parameter %s not found, skipping [%s:elem#%d]",
						idName, file.Name, parsed.Position))
					continue
				}
				elementId = fmt.Sprintf("%v", v)

				// Add parsed to array and store in resultMap
				parsedArr, ok = parsedMap[elementId]
				if ok {
					parsedArr = append(parsedArr, parsed)
				} else {
					parsedArr = []Parsed{parsed}
				}
				parsedMap[elementId] = parsedArr
			}
		} else {
			errList = append(errList, fmt.Sprintf("parsing config: unrecognized JSON type %q [%s]", k, file.Name))
			continue
		}
	}

	// check for conflicting values and unused parameters
	validateParameters(st, parsedMap, &errList)

	// collapse each element into a single data object and load into resultMap
	resultMap = parseConfig(st, sv, parsedMap, &errList)

	// Compile error list into an error message
	if len(errList) > 0 {
		if len(errList) == 1 {
			err = fmt.Errorf(errList[0])
		} else {
			for i = range errList {
				errList[i] = fmt.Sprintf("(#%d) %s", i+1, errList[i])
			}
			err = fmt.Errorf("multiple errors\n%s", strings.Join(errList, "\n"))
		}
	}
	if Debug {
		log.Printf("Parsed %d distinct configurations", len(parsedMap))
	}
	return
}
