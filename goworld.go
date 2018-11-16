package GoWorld

import (
	. "./common"
	"./components/game"
	"./entity"
)

func Run(gameid int, delegate game.IGameDelegate) {
	game.Run(delegate)
}

func RegisterEntity(typeName string, entityPtr entity.IEntity) {
	entity.RegisterEntity(typeName, entityPtr)
}

//func CreateEntity(typeName string) {
//	game.CreateEntity(typeName)
//}

func CreateSpace() {
	entity.CreateSpace()
}

func CreateEntity(typeName string) EntityID {
	return entity.CreateEntity(typeName)
}

func SetSpaceDelegate(delegate entity.ISpaceDelegate) {
	entity.SetSpaceDelegate(delegate)
}

func GetServiceProviders(serviceName string) []EntityID {
	return game.GetServiceProviders(serviceName)
}
