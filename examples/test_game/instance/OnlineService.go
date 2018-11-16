package instance

import (
	"../../../entity"
	"github.com/TonyXMH/GoWorld/gwlog"
)

type OnlineService struct {
	entity.Entity
}

func (s *OnlineService) OnCreated() {
	gwlog.Info("Registering OnlineService ... ")
	s.DeclareService("OnlineService")
}
