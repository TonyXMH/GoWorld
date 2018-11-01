package main

import (
	"../../components/dispatcher"
	"flag"
	"fmt"
	"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/gwlog"
)

var configFile = "goworld.ini"

func dubuglog(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	gwlog.Debug("dispatcher:%s", s)
}

func parseArgs() {
	flag.Parse()
}

func main() {
	cfg := config.GetDispatcher()
	dispatcher := dispatcher.NewDispatcherService(cfg)
	dispatcher.Run()
}
