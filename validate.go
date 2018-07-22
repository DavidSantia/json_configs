package json_configs

import (
	"fmt"
	"reflect"
	"strings"
)

// Check data object (st) fields for any conflicting result map values
func validateParameters(st reflect.Type, parsedMap ParsedMap, errList *[]string) {
	var filenames []string
	var elementId, key, filename, tagName, paramName, paramValue string
	var parsedArr []Parsed
	var paramValuesMap map[string][]string
	var v interface{}
	var i int
	var ok bool

	// Maps for json tag names <-> param names, and valid param names
	tag2param := make(map[string]string)
	param2tag := make(map[string]string)
	validparam := make(map[string]bool)
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)
		tagName, ok = param.Tag.Lookup("json")
		if ok {
			tag2param[tagName] = param.Name
			param2tag[param.Name] = tagName
			validparam[param.Name] = true
		}
	}

	// Validate parameters for each element
	for elementId, parsedArr = range parsedMap {

		// make list of values for each parameter, with filenames found in
		elementParamValuesMap := make(map[string]map[string][]string)

		// make map of parameter names, with filenames found in, to see what isn't used
		unusedParamMap := make(map[string][]string)

		// check across all files parsed
		for _, parsed := range parsedArr {
			if parsed.Position == 0 {
				filename = parsed.DistinctName
			} else {
				filename = fmt.Sprintf("%s:elem#%d", parsed.DistinctName, parsed.Position)
			}

			// load elementParamValuesMap to identify possible conflicting values
			for i = 0; i < st.NumField(); i++ {
				param := st.Field(i)
				// lookup in config dataMap by tag name first, then by param name
				tagName, ok = param2tag[param.Name]
				if ok {
					v, ok = parsed.ElementMap[tagName]
				}
				if !ok {
					v, ok = parsed.ElementMap[param.Name]
				}
				if ok {
					paramValue = fmt.Sprintf("%v", v)

					// make list all filenames for each parameter value
					paramValuesMap, ok = elementParamValuesMap[param.Name]
					if !ok {
						paramValuesMap = make(map[string][]string)
						filenames = []string{filename}
					} else {
						filenames, ok = paramValuesMap[paramValue]
						if ok {
							filenames = append(filenames, filename)
						} else {
							filenames = []string{filename}
						}
					}
					paramValuesMap[paramValue] = filenames
					elementParamValuesMap[param.Name] = paramValuesMap
				}
			}

			// load unusedParamMap to identify possible unused parameters
			for key, _ = range parsed.ElementMap {
				// lookup element key by tag name first
				paramName, ok = tag2param[key]
				if !ok {
					paramName = key
				}
				if !validparam[paramName] {
					filenames, ok = unusedParamMap[paramName]
					if ok {
						filenames = append(filenames, filename)
					} else {
						filenames = []string{filename}
					}
					unusedParamMap[paramName] = filenames
				}
			}
		}

		// List errors for conflicting values
		for paramName, paramValuesMap = range elementParamValuesMap {
			// if there are more than one value, it means settings conflict
			if len(paramValuesMap) > 1 {
				var conflicts []string
				for paramValue, filenames = range paramValuesMap {
					conflicts = append(conflicts, fmt.Sprintf("%q [%s]",
						paramValue, strings.Join(filenames, ",")))
				}
				*errList = append(*errList, fmt.Sprintf("settings for %s conflict, parameter %s: %s",
					elementId, paramName, strings.Join(conflicts, " != ")))
			}
		}

		// List  errors for unused parameters
		for paramName, filenames = range unusedParamMap {
			if len(filenames) == 1 {
				*errList = append(*errList, fmt.Sprintf("unused setting for %s, parameter %s [%s]",
					elementId, paramName, filenames[0]))
			} else {
				*errList = append(*errList, fmt.Sprintf("unused settings for %s, parameter %s: %d occurences [%s]",
					elementId, paramName, len(filenames), strings.Join(filenames, ",")))
			}
		}
	}
	return
}
