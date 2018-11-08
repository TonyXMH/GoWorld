package instance

import "../../../entity"

type SpaceDelegate struct {
	entity.DefaultSpaceDelegate
}

func (delegate *SpaceDelegate) OnSpaceCreated(space *entity.Space) {
	delegate.DefaultSpaceDelegate.OnSpaceCreated(space)
	N := 3
	for i := 0; i < N; i++ {
		space.CreateEntity("Monster")
	}
}
