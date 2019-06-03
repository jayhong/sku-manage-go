package log

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var debug bool
var logPath string
var hook string
var hookAddress string
var serviceHost string

func init() {
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.StringVar(&logPath, "log", "./log/", "log path")
	flag.StringVar(&hook, "hook", "lfshook", "logrus hook")
	flag.StringVar(&hookAddress, "hook_addr", "lfshook", "logrus hook address")
	flag.StringVar(&serviceHost, "s_host", "", "service host")
}

func InitLog(service string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetOutput(ioutil.Discard)
	}

	logrus.Infof("format: %v, %v, %v, %v, %v, %v", debug, logPath, hook, hookAddress, serviceHost, service)

	if hk := buildLogHook(hook, buildConfig(options(service))); hk != nil {
		logrus.AddHook(hk)
	}
}

func buildLogHook(name string, cfg *hookConfig) logrus.Hook {
	switch name {
	case "lfshook":
		return lfsHook(cfg)
	case "elastic":
		return elasticHook(cfg)
	case "logstash":
		return logstashHook(cfg)
	}
	return nil
}

func options(service string) []HookOption {
	return []HookOption{
		PathOption(logPath),
		ServiceOption(service),
		AddressOption(hookAddress),
		HostOption(serviceHost),
	}
}

func buildConfig(options []HookOption) (result *hookConfig) {
	result = &hookConfig{}
	for _, op := range options {
		op(result)
	}
	return
}
