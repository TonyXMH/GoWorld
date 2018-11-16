package main

import (
	"../../../GoWorld"
	"../../components/game"
	"../../gwlog"
	"../../timer"
	. "./instance"
	"time"
)

var gameid = 0

func init() {

}

type gameDelegate struct {
	game.GameDelegate
}

func main() {
	GoWorld.SetSpaceDelegate(&SpaceDelegate{})
	GoWorld.RegisterEntity("Monster", &Monster{})
	GoWorld.RegisterEntity("OnlineService", &OnlineService{})
	GoWorld.RegisterEntity("Avator", &Avatar{})
	GoWorld.Run(gameid, &gameDelegate{})
}

func (game gameDelegate) OnReady() {
	game.GameDelegate.OnReady()
	//GoWorld.CreateEntity("Monster")
	GoWorld.CreateEntity("OnlineService")
	timer.AddCallback(time.Millisecond*1000, game.checkGameStarted)
}

func (game gameDelegate) checkGameStarted() {
	ok := game.isGameStarted()
	gwlog.Info("checkGameStarted:%v", ok)
	if ok {
		game.onGameStarted()
	} else {
		timer.AddCallback(time.Millisecond*1000, game.checkGameStarted)
	}
}

func (game gameDelegate) isGameStarted() bool {
	if len(GoWorld.GetServiceProviders("OnlineService")) == 0 {
		return false
	}
	return true
}

func (game gameDelegate) onGameStarted() {
	GoWorld.CreateSpace()
}
