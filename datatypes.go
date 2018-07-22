package json_configs

// Parse JSON config files into a data object using reflect package
// - Can configure an application using one or more JSON files
// - For example, put general settings in one file, credentials in a second file.

// Parsed{} contains each JSON element, remembering where it was found
// - DistinctName is shortest unique name across all filenames
// - Position is element # within the file: 0 if single element, 1...N if array of N elements
type Parsed struct {
	FileName     string
	DistinctName string
	Position     int
	ElementMap   ElementMap
}

// ElementMap is the JSON element parsed into a key-value map
type ElementMap map[string]interface{}

// When mutiple files are parsed, a field in each element is specified as the Id
// - This element Id is used as the ParsedMap key (so becomes a required field)
type ParsedMap map[string][]Parsed

// The Parsed map form is then collapsed into a single data object result per Id
type ResultMap map[string]interface{}
