package instance

import (
	"../../../entity"
)

type Avatar struct {
	entity.Entity
}

func (e *Avatar) OnCreated() {
	e.Entity.OnCreated()
}
