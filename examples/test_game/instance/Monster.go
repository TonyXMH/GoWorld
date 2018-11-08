package instance

import "../../../entity"

type Monster struct {
	entity.Entity
}

func (e *Monster) OnCreated() {
	e.Entity.OnCreated()
}
