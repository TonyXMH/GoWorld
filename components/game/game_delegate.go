package game

import "github.com/TonyXMH/GoWorld/gwlog"

type IGameDelegate interface {
	OnReady()
}

type GameDelegate struct {
}

func (gd *GameDelegate) OnReady() {
	gwlog.Info("game %d is ready.", gameid)
}
