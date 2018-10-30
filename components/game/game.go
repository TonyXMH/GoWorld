package game

import (
	//"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/gwlog"
	//"github.com/TonyXMH/GoWorld/timer"
	"../../config"
	"../../timer"
	"flag"
	"os"
	"time"
)

var (
	gameid       int
	gameDelegate IGameDelegate
)

func init() {
	parseArgs()
}

func parseArgs() {
	flag.IntVar(&gameid, "gid", 0, "set gameid")
	flag.Parse()
}

func Run(delegate IGameDelegate) {
	gameDelegate = delegate
	cfg := config.GetGame(gameid)
	gwlog.Info("Read game %d config:\n%s\n", gameid, config.DumpPretty(cfg))
	timer.AddCallback(0, func() {
		gameDelegate.OnReady()
	})
	for {
		timer.Tick()
		os.Stderr.Write([]byte{'.'})
		time.Sleep(time.Millisecond * 100)
	}
}
