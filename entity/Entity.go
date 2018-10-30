package entity

import (
	"../uuid"
	"fmt"
	"github.com/TonyXMH/GoWorld/gwlog"
)

type EntityID string

func GenEntityID() EntityID {
	return EntityID(uuid.GenUUID())
}

type Entity struct {
	ID EntityID
	I  IEntity
}

type IEntity interface {
	OnCreated()
}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity<%s>", e.ID)
}

func (e *Entity) OnCreated() {
	gwlog.Debug("%s.OnCreated", e)
}
