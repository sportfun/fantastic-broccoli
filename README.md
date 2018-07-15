# UPDATE README
# GAkisitor [(SportsFun)](https://charlestati.github.io/eip-showcase/index.html)

[![Version](https://img.shields.io/badge/version-alpha-orange.svg)](https://github.com/sportfun/gakisitor/milestones)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/sportfun/gakisitor?status.svg)](https://godoc.org/github.com/sportfun/gakisitor)
[![Build Status](https://travis-ci.org/sportfun/gakisitor.svg?branch=master)](https://travis-ci.org/sportfun/gakisitor)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/9b630adf92f84adfaf89ceed51352304)](https://www.codacy.com/app/xunleii/gakisitor?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=sportfun/gakisitor&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/9b630adf92f84adfaf89ceed51352304)](https://www.codacy.com/app/xunleii/gakisitor?utm_source=github.com&utm_medium=referral&utm_content=sportfun/gakisitor&utm_campaign=Badge_Coverage)

GAkisitor is a Go service for connecting sports equipment to a SportsFun game session



## Features
 * Connect sport equipments to a SportsFun play session
 * Easily expandable thanks to the module mechanism
 * Automatic reload service when the configuration file is modified
 * Works on different embedded systems



## Module mechanism

### Module ... Vas is das ?

A module is a sort of package combining one (or more) metric sensor(s) with a plugin
 * The metric sensors get data from the player activity, like his speed.  
   However, to be compatible with the host, this physical part must respect some restriction (defined by the host's conceptor).
 * The plugin is a sort of driver, converting metric datas to usable data.

### How to develop a module ?

To develop a plugin, you must implement [Module](module/module.go).


#### 1. Module management

The first method called by the GAkisitor is `Start`. It implements all things that must be instantied at first.
```golang
Start(*NotificationQueue, log.Logger) error
```
> ##### Parameters
>  * `NotificationQueue` : instance used to notify information to the server
>  * `log.Logger` : logger instance, used for logging information
>
> ##### Error management
>  * If an error occurs, you need to return an `error` and set the `Panic` state ([see example](example/module_example.go#L78))


In contrast, the last method called before module shutdown is `Stop`. It close, defer, kill, eat, ... all values to stop cleanly the module. If a session is running, don't forget to close it (`StopSession`).
```golang
Stop() error
```
> ##### Error management
>  * If an error occurs, you need to return an `error` and set the `Panic` state.


#### 2. Module configuration

Just after calling `Start`, the module must be configured with `Configure`
```golang
Configure(properties.ModuleDefinition) error
```
> ##### Parameters
>  * `ModuleDefinition` : `struct` representing the module definition (module name, plugin path and configuration, throw an `interface{}`).
>
> ##### Error management
>  * If an error occurs, you need to return an `error` and set the `Panic` state.


#### 3. Session management

A session represents a game session. To reduce consumption (and because it's useless), the sensors must be active only during a session. Moreover, sensors data getting must be asynchronous and therefore started in a `goroutine`.

Like `Start`, to start a session, you must implement `StartSession`.
```golang
StartSession() error
```
> ##### Error management
>  * If an error occurs, you need to return an `error` and set the `Idle` state. Don't forget to close the goroutine if it launched. If the error is really critical, set the `Panic` state.

To stop a session, you must implement `StopSession`. You must close the goroutine here.
```golang
StopSession() error
```
> ##### Error management
>  * If an error occurs, you need to return an `error` and set the `Idle` state. Don't forget to close the goroutine if it launched. If the error is really critical, set the `Panic` state.


#### 4. Processing

To convert sensors data and send them to the server, a `Process` method is required.
```golang
Process() error
```
> ##### Error management
>  * If an error occurs, you need to return an `error`. If the error is really critical, set the `Panic` state.


#### 5. Properties

Some properties must be implemented to get some information about the module.
```golang
Name() string           // Return the module name
State() types.StateType // Return the current module state
```
