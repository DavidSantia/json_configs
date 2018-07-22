package json_configs

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Parse config dataMap entries corresponding to data object (st, sv) fields
func parseConfig(st reflect.Type, sv reflect.Value, parsedMap ParsedMap, errList *[]string) (resultMap ResultMap) {
	var err error
	var elementId, filename, tagName, paramValue, paramType string
	var v interface{}
	var parsedArr []Parsed
	var dur time.Duration
	var date time.Time
	var f float64
	var n int64
	var i int
	var ok bool

	// Map for json param names to tag names
	param2tag := make(map[string]string)
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)
		tagName, ok = param.Tag.Lookup("json")
		if ok {
			param2tag[param.Name] = tagName
		}
	}

	resultMap = make(ResultMap)
	for elementId, parsedArr = range parsedMap {

		// Make map of parameter names that will need to be reset after parsing element
		clearParamMap := make(map[string]bool)

		for _, parsed := range parsedArr {
			if parsed.Position == 0 {
				filename = parsed.DistinctName
			} else {
				filename = fmt.Sprintf("%s:elem#%d", parsed.DistinctName, parsed.Position)
			}

			// Iterate through element data fields, parse into correct type
			for i = 0; i < st.NumField(); i++ {
				param := st.Field(i)

				// lookup in ElementMap by tag name first, then by param name
				tagName, ok = param2tag[param.Name]
				if ok {
					v, ok = parsed.ElementMap[tagName]
				}
				if !ok {
					v, ok = parsed.ElementMap[param.Name]
				}

				if ok {
					paramValue = fmt.Sprintf("%v", v)
					paramType = param.Type.Name()
					clearParamMap[param.Name] = true

					switch paramType {
					case "string":
						sv.Field(i).SetString(paramValue)
					case "float64":
						f, err = strconv.ParseFloat(paramValue, 64)
						if err != nil {
							*errList = append(*errList,
								fmt.Sprintf("setting for %s invalid, parameter %s: float %s [%s]",
									elementId, param.Name, paramValue, filename))
							continue
						}
						sv.Field(i).SetFloat(f)
					case "int", "int64":
						n, err = strconv.ParseInt(paramValue, 10, 64)
						if err != nil {
							*errList = append(*errList,
								fmt.Sprintf("setting for %s invalid, parameter %s: integer %s [%s]",
									elementId, param.Name, paramValue, filename))
							continue
						}
						sv.Field(i).SetInt(n)
					case "bool":
						ok, err = strconv.ParseBool(paramValue)
						if err != nil {
							*errList = append(*errList,
								fmt.Sprintf("setting for %s invalid, parameter %s: boolean %s [%s]",
									elementId, param.Name, paramValue, filename))
							continue
						}
						sv.Field(i).SetBool(ok)
					case "Duration":
						dur, err = time.ParseDuration(paramValue)
						if err != nil {
							*errList = append(*errList,
								fmt.Sprintf("setting for %s invalid, parameter %s: duration %s [%s]",
									elementId, param.Name, paramValue, filename))
							continue
						}
						sv.Field(i).Set(reflect.ValueOf(dur))
					case "Time":
						date, err = time.Parse("2006-01-02T15:04:05Z", paramValue)
						if err != nil {
							date, err = time.Parse("2006-01-02", paramValue)
						}
						if err != nil {
							*errList = append(*errList,
								fmt.Sprintf("setting for %s invalid, parameter %s: date %s [%s]",
									elementId, param.Name, paramValue, filename))
							continue
						}
						sv.Field(i).Set(reflect.ValueOf(date))
					default:
						*errList = append(*errList,
							fmt.Sprintf("setting for %s invalid, parameter %s: unsupported type %s [%s]",
								elementId, param.Name, paramType, filename))
						continue
					}
				}
			}
		}
		// Store data object in resultMap
		resultMap[elementId] = sv.Interface()

		// Clear data object for next element Id
		clearConfig(st, sv, clearParamMap)
	}
	return
}

func clearConfig(st reflect.Type, sv reflect.Value, clearParamMap map[string]bool) (err error) {
	var paramType string
	var zeroD time.Duration
	var zeroT time.Time
	var i int
	var ok bool

	// Iterate through data fields, clear settings
	for i = 0; i < st.NumField(); i++ {
		param := st.Field(i)

		ok = clearParamMap[param.Name]
		if ok {
			paramType = param.Type.Name()

			switch paramType {
			case "string":
				sv.Field(i).SetString("")
			case "float64":
				sv.Field(i).SetFloat(0)
			case "int", "int64":
				sv.Field(i).SetInt(0)
			case "bool":
				sv.Field(i).SetBool(false)
			case "Duration":
				sv.Field(i).Set(reflect.ValueOf(zeroD))
			case "Time":
				sv.Field(i).Set(reflect.ValueOf(zeroT))
			default:
				err = fmt.Errorf("unsupported type %s [%s]", paramType, param.Name)
				return
			}
		}
	}

	return
}
