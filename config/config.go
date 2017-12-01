package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// RawConfig is a type for unknowing configuration
type RawConfig interface{}

// GAkisitorConfig is the configuration that comes from loading
// the configuration file
type GAkisitorConfig struct {
	file     string // path to the file where this configuration was loaded from
	isLoaded bool   // loading state of this configuration

	System  SystemDefinition   `json:"system"`
	Modules []ModuleDefinition `json:"modules"`
	Log     []LogDefinition    `json:"log"`
}

// SystemDefinition define some information about the current system
// and network
type SystemDefinition struct {
	LinkID     string `json:"link_id"`
	DeviceName string `json:"device"`

	ServerIP   string `json:"ip"`
	ServerPort int    `json:"port"`
	ServerSSL  bool   `json:"ssl"`
}

// ModuleDefinition define module information like name or
// internal configuration
type ModuleDefinition struct {
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Config RawConfig `json:"config"`
}

// LogDefinition define log information, used during logger
// instantiation
type LogDefinition struct {
	File     string `json:"file"`
	Format   string `json:"format"`
	Encoding string `json:"encoding"`
	Level    string `json:"level"`
}

// FilePtr return a pointer of the configuration file path
func (p *GAkisitorConfig) FilePtr() *string {
	return &p.file
}

// IsLoaded return the current loading state
func (p *GAkisitorConfig) IsLoaded() bool {
	return p.isLoaded
}

// Load performs the loading file and unmarshal it into itself
func (p *GAkisitorConfig) Load() {
	p.isLoaded = false

	raw, err := ioutil.ReadFile(p.file)
	if err != nil {
		return
	}

	if err := json.Unmarshal(raw, p); err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("impossible to unmarshal the configuration file: %s", err.Error())
		log.SetOutput(os.Stdout)
		return
	}
	p.isLoaded = true
}

// WaitReconfiguration wait until the configuration file was
// modified. Next, this function reload the file
// TODO: Add timeout
func (p *GAkisitorConfig) WaitReconfiguration() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("impossible to load or monitor file '%s': '%s'", p.file, err.Error())
	}
	defer watcher.Close()

	err = watcher.Add(p.file)
	if err != nil {
		log.Fatalf("impossible to load or monitor file '%s': '%s'", p.file, err.Error())
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				p.Load()
				return
			}
		case err := <-watcher.Errors:
			log.Printf("error during file monitoring: '%s'", err.Error())
		}
	}
}
