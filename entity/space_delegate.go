package entity

import "github.com/TonyXMH/GoWorld/gwlog"

var (
	spaceDelegate ISpaceDelegate = &DefaultSpaceDelegate{}
)

func SetSpaceDelegate(delegate ISpaceDelegate) {
	spaceDelegate = delegate
}

type ISpaceDelegate interface {
	OnSpaceCreated(space *Space)
}

type DefaultSpaceDelegate struct {
}

func (delegate *DefaultSpaceDelegate) OnSpaceCreated(space *Space) {
	gwlog.Debug("OnSpaceCreated: %s", space)
}
