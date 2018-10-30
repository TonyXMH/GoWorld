package entity

import (
	"github.com/TonyXMH/GoWorld/gwlog"
	"testing"
)

type TestEntity struct {
	Entity
}

func TestRegisterEntity(t *testing.T) {
	RegisterEntity("TestEntity", &TestEntity{})
}

func TestGenEntityID(t *testing.T) {
	eid := GenEntityID()
	gwlog.Info("TestGenEntityID:%s", eid)
}

func TestEntityManager(t *testing.T) {

}
