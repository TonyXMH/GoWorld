package entity

import (
	. "../common"
	"../components/dispatcher/dispatcher_client"
	"../timer"
	"github.com/TonyXMH/GoWorld/gwlog"
	"reflect"
)

var (
	registeredEntityTypes = map[string]reflect.Type{}
	entityManager         = newEntityManager()
)

type EntityManager struct {
	entities EntityMap
}

func newEntityManager() *EntityManager {
	return &EntityManager{
		entities: EntityMap{},
	}
}

func (em *EntityManager) put(entity *Entity) {
	em.entities.Add(entity)
}

func (em *EntityManager) del(id EntityID) {
	em.entities.Del(id)
}

func (em *EntityManager) get(id EntityID) IEntity {
	return em.entities.Get(id)
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

func createEntity(typeName string, space *Space) EntityID {
	gwlog.Debug("CreateEntity: %s in space %s", typeName, space)
	entityType, ok := registeredEntityTypes[typeName]
	if !ok {
		gwlog.Panicf("unknown entity type:%s", typeName)
	}
	entityID := GenEntityID()
	entityPtrVal := reflect.New(entityType)
	entity := reflect.Indirect(entityPtrVal).FieldByName("Entity").Addr().Interface().(*Entity)
	entity.ID = entityID
	entity.I = entityPtrVal.Interface().(IEntity)
	entity.TypeName = typeName
	entity.timers = map[*timer.Timer]struct{}{}
	initAOI(&entity.aoi)
	entity.I.OnInit()
	entityManager.put(entity)
	entity.I.OnCreated()
	if space != nil {
		space.enter(entity)
	}
	return entityID
}

func CreateEntity(typeName string) EntityID {
	return createEntity(typeName, nil)
}

func call(id EntityID, method string, args []interface{}) {
	dispatcher_client.GetDispatcherClientForSend().SendCallEntityMethod(id, method)
}
