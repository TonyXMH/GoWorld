package entity

import (
	"github.com/TonyXMH/GoWorld/gwlog"
)

const (
	SPACE_ENTITY_TYPE = "__space__"
)

type Space struct {
	Entity
	entities EntitySet
}

func init() {
	RegisterEntity(SPACE_ENTITY_TYPE, &Space{})
}

func (space *Space) OnInit() {
	space.entities = EntitySet{}
}

func (space *Space) OnCreated() {
	gwlog.Debug("%s.OnCreated", space)
	space.Post(func() {
		spaceDelegate.OnSpaceCreated(space)
	})
}

func (space *Space) CreateEntity(typeName string) {
	//entityID:=createEntity(typeName,space)
	//gwlog.Info("%s.createEntity %s:%s",space,typeName,entityID)
	createEntity(typeName, space)
}

func (space *Space) enter(entity *Entity) {
	gwlog.Info("%s.enter<<<%s", space, entity)
	entity.space = space
	for other := range space.entities {
		entity.interest(other)
		other.interest(entity)
	}
	space.entities.Add(entity)
	entity.I.OnEnterSpace()
}

func (space *Space) leave(entity *Entity) {
	entity.space = nil
	space.entities.Del(entity)
	for other := range space.entities {
		entity.uninterest(other)
		other.uninterest(entity)
	}
}
