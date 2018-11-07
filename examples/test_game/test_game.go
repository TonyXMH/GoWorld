package main

import (
	"../../../GoWorld"
	"../../components/game"
	"../../entity"
	"../../gwlog"
	"time"
)

var gameid = 0

type Monster struct {
	entity.Entity
}

func init() {

}

type gameDelegate struct {
	game.GameDelegate
}

func main() {
	GoWorld.SetSpaceDelegate(&SpaceDelegate{})
	GoWorld.RegisterEntity("Monster", &Monster{})
	GoWorld.Run(gameid, &gameDelegate{})
}

func (game gameDelegate) OnReady() {
	game.GameDelegate.OnReady()
	//GoWorld.CreateEntity("Monster")
	GoWorld.CreateSpace()
}

type SpaceDelegate struct {
	entity.DefaultSpaceDelegate
}

func (delegate *SpaceDelegate) OnSpaceCreated(space *entity.Space) {
	delegate.DefaultSpaceDelegate.OnSpaceCreated(space)
	//space.CreateEntity("Monster")
	N := 3
	for i := 0; i < N; i++ {
		space.CreateEntity("Monster")
	}

}

func (e *Monster) OnCreated() {
	e.Entity.OnCreated()
	gwlog.Info("Creating callback ...")
	e.AddTimer(time.Second, func() {
		gwlog.Info("%s.Neighbors=%v", e, e.Neighbors())
		for _other := range e.Neighbors() {
			if _other.TypeName != "Monster" {
				continue
			}
			other := _other.I.(*Monster)
			gwlog.Info("%s is a neighbor of %s", other, e)
		}
	})
}
