package GoWorld

import (
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

func SetSpaceDelegate(delegate entity.ISpaceDelegate) {
	entity.SetSpaceDelegate(delegate)
}
