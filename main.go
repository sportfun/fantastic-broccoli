package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	fluentd "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
	"github.com/sportfun/gakisitor/profile"
	"github.com/takama/daemon"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	ServiceName = "gakisitor"
	ServiceDesc = "Sportsfun user action acquisition service"
)

var (
	Gakisitor = struct {
		daemon.Daemon
		profile.Profile
		*scheduler
	}{
		scheduler: &scheduler{
			workers: map[string]*worker{},
		},
	}

	debugEnabled bool
	profilePath  string
)

func init() {
	flag.BoolVar(&debugEnabled, "debug", false, "Enable debug logging")
	flag.StringVar(&profilePath, "conf", "/etc/gakisitor/profile.json", "Path to profile file")
}

func exitIfFail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func command(cmd string) (string, error) {
	switch cmd {
	case "install":
		return Gakisitor.Install()
	case "remove":
		return Gakisitor.Remove()
	case "start":
		return Gakisitor.Start()
	case "stop":
		return Gakisitor.Stop()
	case "status":
		return Gakisitor.Status()
	default:
		return fmt.Sprintf("Usage: %s install | remove | start | stop | status", os.Args[0]), nil
	}
}
func prepare() {
	// Configure logger
	log.SetLevel(log.InfoLevel)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	switch Gakisitor.Log.Format {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "fluentd":
		log.SetFormatter(&fluentd.FluentdFormatter{})
	case "system":
		log.SetFormatter(&prefixed.TextFormatter{})
	default:
		log.SetFormatter(&log.JSONFormatter{})
	}

	switch Gakisitor.Log.Path {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "":
		out, err := os.Open("/var/log/gakisitor.log")
		exitIfFail(err)
		log.SetOutput(out)
	default:
		out, err := os.Open(Gakisitor.Log.Path)
		exitIfFail(err)
		log.SetOutput(out)
	}
}

func main() {
	var err error
	var cancel context.CancelFunc
	var ctx context.Context

	flag.Parse()

	// Create daemon
	Gakisitor.Daemon, err = daemon.New(ServiceName, ServiceDesc)
	exitIfFail(err)

	// If received any kind of command, manage it
	if len(flag.Args()) > 1 {
		status, err := command(flag.Args()[1])
		switch {
		case err != nil:
			exitIfFail(err)
		case status != "":
			fmt.Fprintln(os.Stdout, status)
			return
		}
	}

	// Prepare Gakisitor environment
	exitIfFail(Gakisitor.Load(profilePath))
	prepare()

	// Subscribe if the conf file was altered
	Gakisitor.SubscribeAlteration(func(profile *profile.Profile, err error) {
		// If error, stop all
		if err != nil {
			log.Errorf("Failed with file events: %s", err)
			cancel()
			return
		}

		// Prepare new Gakisitor environment
		prepare()
		cancel()
	})

	// Start Gakisitor
	for {
		ctx, cancel = context.WithCancel(context.Background())
		Gakisitor.scheduler = &scheduler{
			workers:             Gakisitor.scheduler.workers,
			bus:                 bus.New(),
			ctx:                 ctx,
			deadSig:             make(chan string),
			workerRetryMax:      int32(Gakisitor.Scheduler.Worker.Retry),
			workerRetryInterval: time.Millisecond * time.Duration(Gakisitor.Scheduler.Worker.Interval),
		}

		restart := Gakisitor.scheduler.Run()
		if !<-restart {
			log.Infof("Stop service '%s'", ServiceDesc)
			break
		}
	}
}
