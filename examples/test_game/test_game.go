package main

import (
	"../../../GoWorld"
	"../../components/game"
	"../../entity"
)

var gameid = 0

type TestEntity struct {
	entity.Entity
}

func init() {

}

type gameDelegate struct {
	game.GameDelegate
}

func main() {
	GoWorld.RegisterEntity("TestEntity", &TestEntity{})
	GoWorld.Run(gameid, &gameDelegate{})
}

func (game gameDelegate) OnReady() {
	game.GameDelegate.OnReady()
	GoWorld.CreateEntity("TestEntity")
}
