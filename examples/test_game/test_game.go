package main

import (
	"../../../GoWorld"
	"../../components/game"
	. "./instance"
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
	GoWorld.RegisterEntity("OnlineService",&OnlineService{})
	GoWorld.RegisterEntity("Avator",&Avatar{})
	GoWorld.Run(gameid, &gameDelegate{})
}

func (game gameDelegate) OnReady() {
	game.GameDelegate.OnReady()
	//GoWorld.CreateEntity("Monster")
	GoWorld.CreateEntity("OnlineService")
	GoWorld.CreateSpace()
}
