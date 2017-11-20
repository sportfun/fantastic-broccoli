package properties

import (
	"encoding/json"
	"io/ioutil"
	"github.com/fsnotify/fsnotify"
	"log"
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
	log.Printf("monitoring")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("impossible to load or monitor file '%s': '%s'", properties.originPath, err)
	}
	defer watcher.Close()

	err = watcher.Add(properties.originPath)
	if err != nil {
		log.Fatalf("impossible to load or monitor file '%s': '%s'", properties.originPath, err)
	}

	for {
		log.Printf("monitoring")
		select {
		case event := <-watcher.Events:
			log.Printf("event: %s -> %s", event.Op, event.Name)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("reload file")
				*properties = *LoadFrom(properties.originPath)
				return
			}
		case err := <-watcher.Errors:
			log.Printf("error during file monitoring: '%s'", err)
		}
	}
}
