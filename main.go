package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	fluentd "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"github.com/sportfun/gakisitor/event/bus"
	"github.com/sportfun/gakisitor/profile"
	"github.com/takama/daemon"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	serviceName = "gakisitor"
	serviceDesc = "Sportsfun user action acquisition service"
)

var (
	// Gakisitor is the main instance of the service
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
	logrus.SetLevel(logrus.InfoLevel)
	if debugEnabled {
		logrus.SetLevel(logrus.DebugLevel)
	}

	switch Gakisitor.Log.Format {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "fluentd":
		logrus.SetFormatter(&fluentd.FluentdFormatter{})
	case "system":
		logrus.SetFormatter(&prefixed.TextFormatter{})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	switch Gakisitor.Log.Path {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	case "":
		out, err := os.OpenFile("/var/logrus/gakisitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		exitIfFail(err)
		logrus.SetOutput(out)
	default:
		out, err := os.OpenFile(Gakisitor.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		exitIfFail(err)
		logrus.SetOutput(out)
	}
}

func main() {
	var err error
	var cancel context.CancelFunc
	var ctx context.Context

	flag.Parse()

	// Create daemon
	Gakisitor.Daemon, err = daemon.New(serviceName, serviceDesc)
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
			logrus.Errorf("Failed with file events: %s", err)
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
			logrus.Infof("Stop service '%s'", serviceDesc)
			break
		}
	}
	cancel()
}
