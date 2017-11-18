package properties

import (
	"encoding/json"
	"io/ioutil"
)

func LoadFrom(path string) *Properties {
	var properties Properties
	properties.originPath = path
	properties.isLoaded = false

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return &properties
	}

	json.Unmarshal(raw, &properties)
	properties.isLoaded = true
	return &properties
}

func WaitReconfiguration(properties *Properties) {
	// Check if file is edited (https://godoc.org/github.com/fsnotify/fsnotify)
	*properties = *LoadFrom(properties.originPath)
}
