package entity

import (
	"github.com/TonyXMH/GoWorld/gwlog"
	"reflect"
)

var (
	registeredEntityTypes = map[string]reflect.Type{}
	entityManager         = newEntityManager()
)

type EntityManager struct {
	entities map[EntityID]IEntity
}

func newEntityManager() *EntityManager {
	return &EntityManager{
		entities: map[EntityID]IEntity{},
	}
}

func (em *EntityManager) Put(entity *Entity) {
	em.entities[entity.ID] = entity
}

func (em *EntityManager) Get(id EntityID) IEntity {
	return em.entities[id]
}

func RegisterEntity(typename string, entityPtr IEntity) {

	if _, ok := registeredEntityTypes[typename]; ok {
		gwlog.Panicf("RegisterEntity:Entity type %s already registered", typename)
	}
	entityVal := reflect.Indirect(reflect.ValueOf(entityPtr))
	entityType := entityVal.Type()
	registeredEntityTypes[typename] = entityType
	gwlog.Debug(">>> RegisterEntity %s=>%s<<<", typename, entityType.Name())
}

func CreateEntity(typeName string) {
	gwlog.Debug("CreateEntity: %s", typeName)
	entityType, ok := registeredEntityTypes[typeName]
	if !ok {
		gwlog.Panicf("unknown entity type:%s", typeName)
	}
	entityID := GenEntityID()
	entityPtrVal := reflect.New(entityType)
	entity := reflect.Indirect(entityPtrVal).FieldByName("Entity").Addr().Interface().(*Entity)
	entity.ID = entityID
	entity.I = entityPtrVal.Interface().(IEntity)
	entityManager.Put(entity)
	entity.I.OnCreated()
}
