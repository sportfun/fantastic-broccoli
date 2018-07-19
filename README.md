# GAkisitor [(SportsFun)](https://charlestati.github.io/eip-showcase/index.html)

[![Version](https://img.shields.io/badge/version-v2.0.0-green.svg)](https://github.com/sportfun/gakisitor/releases/edit/v2.0)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/sportfun/gakisitor?status.svg)](https://godoc.org/github.com/sportfun/gakisitor)
[![Build Status](https://travis-ci.org/sportfun/gakisitor.svg?branch=master)](https://travis-ci.org/sportfun/gakisitor)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/9b630adf92f84adfaf89ceed51352304)](https://www.codacy.com/app/xunleii/gakisitor?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=sportfun/gakisitor&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/9b630adf92f84adfaf89ceed51352304)](https://www.codacy.com/app/xunleii/gakisitor?utm_source=github.com&utm_medium=referral&utm_content=sportfun/gakisitor&utm_campaign=Badge_Coverage)

GAkisitor is a Go service for connecting sports equipment to a SportsFun game session

## Features

- Connect sport equipments to a SportsFun play session
- Easily expandable through plugins
- Automatic reload service when the configuration file is modified
- Works on different embedded systems *(but only Unix)*

## Modules and plugins

### Module, plugins ... Need some defintions

#### Plugins

A plugin is, hum ... here is the definition
> In computing, a plug-in (or plugin, add-in, addin, add-on, addon, or extension) is a software component that adds a specific feature to an existing computer program
>
> -- <cite>[Wikipedia](https://en.wikipedia.org/wiki/Plug-in_(computing))</cite>

This awesome technology use the *"simple"* Golang plugin library to easily expend the Gakisitor functionnalities.dfsfsdf

I will describe you how to develop one in the next chapter "How to make a plugin in less than 54 steps ?".

#### Modules

This time, I will use *my* defintion of the ***Module***.  
A module is a sort of package combining one (or more) real sensor(s) with a plugin

- The metric sensors get data from the player activity, like his speed.  
   However, to be compatible with the host, this physical part must respect some restriction (defined by the host's conceptor).
- The plugin is the driver in charge to convert metrics data into usable data.

### How to make a plugin in less than 54 steps ?

#### Step 1 :: Architecture

A plugin is represented by a simple structure:

```go
type Plugin struct {
  // The plugin name. It will be used by the server/game
  // engine to know which plugin the data comes from.
  Name string

  // Start the plugin instance with the plugin profile and channels. You
  // MUST check the profile before starting the process.
  //
  // For more information about plugin, see the package description.
  // For more information about plugin channels, see the Chan structure above.
  Instance func(ctx context.Context, profile profile.Plugin, channels Chan) error
}
```
The `Plugin.Name` is used to know where the metric (sensor data) comes from.

> Be careful, the name must be unique

#### Step 2 :: Instance

To be able to stop & start the module, we need an *Instanciator*; a function that creates a live instance of your plugin. For that, this *instanciator* takes 3 arguments:

- A context. This context MUST BE used by your plugin to allow to stop it properly by the Gakisitor. If you don't know how it works, see [this article](https://blog.golang.org/context).
- A profile. This profile contains the user configuration, writed into the configuration file. A buit-in function is available to get properties: `Profile.AccessTo(paths ...interface{})`. With this tool, you can easily access to the required property only with its path.
- A channel list. It contains all channel used by your plugin to communicate with our acquisitor.
  - `Chan.Data` is where you send the metrics acquired by the sensior. **It only takes JSON serializable data.**
  - `Chan.Status` is where you sent the plugin status. It will used to know if your plugin is running or not.
  - `Chan.Instruction` contains the instructions sent by the Gakisitor. Currently, only theses three instructions are provided
    - `StatusPluginInstruction` = send a the current status
    - `StartSessionInstruction` = start a game session (you MUST retrieve user input during this session)
    - `StopSessionInstruction` = stop the game session (you MUST stop your retrieving user input)

#### Step 3 :: Test your plugin

To know if your plugin is compliant with this system, a test tools is provided:  

```go
  PluginValidityCheckert(*testing.T, *plugin.Plugin, PluginTestDesc)
```

It checks if your plugin works like expected when we send an instruction.

> The `PluginTestDesc` contains a custom profile and a value checker to detect if your returned data is what you expect.

#### Step 54 :: Compile it

Of course, to work as a Golang plugin (and loaded by us), you need to compile this `struct` with the flags `--buildmode=plugin`. Currently, this is not available on Windows, but you can use Docker to do that ([library/golang](https://hub.docker.com/_/golang/)).

## Configuration

### Example better than precept

```
{
  "link_id": "6e920d2c-28c8-4667-8b5c-14769ef023e2",  // unique identifier of the gakisitor instance
  "scheduler": {          // scheduler configuration
    "timing": {
      "ttl": 50000,       // Time to live before a worker is claim as dead
      "ttw": 1500,        // Time to wait before a worker will be respawn
      "ttr": 5000         // Time to respawn, avoid infinite respawn
    }
  },
  "network": {                    // network configuration
    "host_address": "127.0.0.1",  // server address
    "port": 8080,                 // server port
    "ssl": false                  // enable SSL (no custom certificates allowed)
  },
  "modules": [                        // plugin configurations
    {                                 // simple plugin conf
      "name": "example",              // plugin name
      "path": "example.gkplugin.so",  // path where the plugin binary is located
      "config": {                     // raw plugin configuration (here for a RPM simulator)
        "rpm.min": 0,                 // minimum RPM value
        "rpm.max": 1200.0,            // maximum RPM value
        "rpm.step": 250,              // RPM step between two value
        "rpm.precision": 1000         // RPM precision
      }
    }
  ]
}
```
