package instance

import (
	"../../../entity"
	"github.com/TonyXMH/GoWorld/gwlog"
	"time"
)

type OnlineService struct {
	entity.Entity
}

func (s *OnlineService) OnCreated() {
	s.AddCallback(time.Second*3, func() {
		gwlog.Info("Registering OnlineService...")
		s.RegisterService("OnlineService")
	})
}
