package entity

import (
	"../timer"
	"../uuid"
	"fmt"
	"github.com/TonyXMH/GoWorld/gwlog"
	"time"
)

const ENTITY_LENGTH = uuid.UUID_LENGTH

type EntityID string

func GenEntityID() EntityID {
	return EntityID(uuid.GenUUID())
}

type Entity struct {
	ID       EntityID
	TypeName string
	I        IEntity
	space    *Space
	aoi      AOI
	timers   map[*timer.Timer]struct{}
}

type IEntity interface {
	OnInit()
	OnCreated()
	OnDestroy()
}

func (e *Entity) String() string {
	return fmt.Sprintf("%s<%s>", e.TypeName, e.ID)
}

func (e *Entity) Destroy() {
	gwlog.Info("%s.Destroy.", e)
	if e.space != nil {
		e.space.leave(e)
	}
	e.clearTimers()
	e.I.OnDestroy()
	entityManager.del(e.ID)
}
func (e *Entity) interest(other *Entity) {
	e.aoi.interest(other)
}

func (e *Entity) uninterest(other *Entity) {
	e.aoi.uninterest(other)
}

func (e *Entity) Neighbors() EntitySet {
	return e.aoi.neighbors
}
func (e *Entity) AddCallback(d time.Duration, cb timer.CallbackFunc) {
	var t *timer.Timer
	t = timer.AddCallback(d, func() {
		delete(e.timers, t)
		cb()
	})
	e.timers[t] = struct{}{}
}

func (e *Entity) Post(cb timer.CallbackFunc) {
	e.AddCallback(0, cb)
}
func (e *Entity) AddTimer(d time.Duration, cb timer.CallbackFunc) {
	t := timer.AddTimer(d, cb)
	e.timers[t] = struct{}{}
}

func (e *Entity) clearTimers() {
	for t := range e.timers {
		t.Cancel()
	}
	e.timers = map[*timer.Timer]struct{}{}
}
func (e *Entity) OnInit() {
	gwlog.Warn("%s.OnInit not implemented", e)
}

func (e *Entity) OnCreated() {
	gwlog.Info("%s.OnCreated", e)
}

func (e *Entity) OnDestroy() {

}
