// Package profile provides types used to configure the Gakisitor and its plugins.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

// Raw type represents any type used to configure a plugin.
type Raw interface{}

// Profile represents the configuration of the Gakisitor. It contains all
// required information about network, scheduler and plugins.
type Profile struct {
	file     string // path of the current profile instance
	isLoaded bool   // state of the profile. true if the profile is loaded, else false

	LinkID string `json:"link_id"` // unique link id, used to identify the Gakisitor

	// Scheduler configuration.
	Scheduler struct {
		// information about worker failure
		Worker struct {
			Retry    int `json:"retry"`    // how many time the scheduler restart a worker with a 'failed' tag before shutdown
			Interval int `json:"interval"` // in which interval a worker which restart is tagged has failed (in ms)
		} `json:"worker"`
	} `json:"scheduler"`

	// Network configuration
	Network struct {
		HostAddress string `json:"host_address"` // host address (IPv4 / IPv6)
		Port        int    `json:"port"`         // host port
		EnableSsl   bool   `json:"enable_ssl"`   // enable SSL (if required)
	} `json:"network"`

	// Plugins configuration
	Plugins []Plugin `json:"plugins"`

	// Log configuration
	Log struct {
		Format string `json:"format"`
		Path   string `json:"path"`
	}
}

// Plugin describes the plugin profile.
type Plugin struct {
	Name   string `json:"name"`   // plugin name
	Path   string `json:"path"`   // plugin library path
	Config Raw    `json:"config"` // plugin configuration
}

// Errors which can be occur in AccessTo function.
var (
	ErrEmptyAccessPath   = errors.New("empty access path")
	ErrInvalidAccessPath = errors.New("invalid access path")
	ErrInvalidIndexType  = errors.New("invalid index path type (must be a string or an int)")
	ErrOutOfBoundIndex   = errors.New("out of bound index path")
)

// Load loads the profile from a file. The optional parameter change the
// internal profile file path, if it already set.
func (profile *Profile) Load(file ...string) error {
	profile.isLoaded = false

	if len(file) > 0 {
		profile.file = file[0]
	}

	raw, err := ioutil.ReadFile(profile.file)
	if err != nil {
		return fmt.Errorf("impossible to read the profile file: %s", err.Error())
	}

	if err := json.Unmarshal(raw, profile); err != nil {
		return fmt.Errorf("impossible to unmarshal the profile file: %s", err.Error())
	}
	profile.isLoaded = true

	return nil
}

// SubscribeAlteration subscribes an handler, called when the profile file was altered.
func (profile *Profile) SubscribeAlteration(handler func(profile *Profile, err error)) (*fsnotify.Watcher, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler can't be nil")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("impossible to load or monitor file '%s': %s", profile.file, err.Error())
	}

	err = watcher.Add(profile.file)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("impossible to load or monitor file '%s': %s", profile.file, err.Error())
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := profile.Load(profile.file); err != nil {
						handler(profile, err)
					} else {
						handler(profile, nil)
					}
				}
			case err, open := <-watcher.Errors:
				if !open {
					return
				}
				handler(nil, err)
			}
		}
	}()
	return watcher, nil
}

// AccessTo allows to access to an item into the plugin profile raw.
func (profile *Plugin) AccessTo(paths ...interface{}) (interface{}, error) {
	if len(paths) == 0 {
		return nil, ErrEmptyAccessPath
	}

	var currentNode = profile.Config
	for _, path := range paths {
		switch idx := path.(type) {
		case string:
			node, isObj := currentNode.(map[string]interface{})
			if !isObj {
				return currentNode, ErrInvalidAccessPath
			}

			var exists bool
			currentNode, exists = node[idx]
			if !exists {
				return node, ErrInvalidAccessPath
			}

		case int:
			node, isObj := currentNode.([]interface{})
			if !isObj {
				return currentNode, ErrInvalidAccessPath
			}

			if idx >= len(node) {
				return currentNode, ErrOutOfBoundIndex
			}
			currentNode = node[idx]

		default:
			return currentNode, ErrInvalidIndexType
		}

	}

	return currentNode, nil
}
