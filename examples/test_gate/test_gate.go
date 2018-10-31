package main

import (
	"../../components/gate"
	"flag"
	"github.com/TonyXMH/GoWorld/config"
)

var gateid int

func parseArgs() {
	flag.IntVar(&gateid, "gid", 1, "set gateid") //测试需要与config.ini配套gate1必须填1
	flag.Parse()
}
func main() {
	parseArgs()
	gatecfg := config.GetGate(gateid)
	gate.NewGateServer(gatecfg).Run()
}
