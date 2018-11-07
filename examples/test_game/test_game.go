package main

import (
	"../../../GoWorld"
	"../../components/game"
	"../../entity"
	"../../gwlog"
	"time"
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
	GoWorld.SetSpaceDelegate(&SpaceDelegate{})
	GoWorld.RegisterEntity("TestEntity", &TestEntity{})
	GoWorld.Run(gameid, &gameDelegate{})
}

func (game gameDelegate) OnReady() {
	game.GameDelegate.OnReady()
	//GoWorld.CreateEntity("TestEntity")
	GoWorld.CreateSpace()
}

type SpaceDelegate struct {
	entity.DefaultSpaceDelegate
}

func (delegate *SpaceDelegate) OnSpaceCreated(space *entity.Space) {
	delegate.DefaultSpaceDelegate.OnSpaceCreated(space)
	//space.CreateEntity("TestEntity")
	N := 3
	for i := 0; i < N; i++ {
		space.CreateEntity("TestEntity")
	}

}

func (e *TestEntity) OnCreated() {
	e.Entity.OnCreated()
	gwlog.Info("Creating callback ...")
	e.AddTimer(time.Second, func() {
		gwlog.Info("%s.Neighbors=%v", e, e.Neighbors())
	})
}
