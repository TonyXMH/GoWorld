package game

import (
	//"github.com/TonyXMH/GoWorld/config"
	//"github.com/TonyXMH/GoWorld/gwlog"
	//"github.com/TonyXMH/GoWorld/timer"
	//"../../config"
	//"../../timer"
	"flag"
	//"os"
	//"time"
)

var (
	gameid int
)

func init() {
	parseArgs()
}

func parseArgs() {
	flag.IntVar(&gameid, "gid", 0, "set gameid")
	flag.Parse()
}

func Run(delegate IGameDelegate) {
	newGameService(gameid, delegate).run()
}
