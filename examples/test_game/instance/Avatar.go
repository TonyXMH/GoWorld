package instance

import (
	"../../../entity"
	"../../../gwlog"
	"../../../../GoWorld"
)

type Avatar struct {
	entity.Entity
}

func (a *Avatar) OnCreated() {
	a.Entity.OnCreated()
	OnlineServiceEid := GoWorld.GetServiceProviders("OnlineService")[0]
	gwlog.Debug("Found OnlineService:%s", OnlineServiceEid)
	a.Call(OnlineServiceEid, "CheckIn", a.ID)
}
