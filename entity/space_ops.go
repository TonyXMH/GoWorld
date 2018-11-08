package entity

import (
	. "../common"
)

func CreateSpace() EntityID {
	return createEntity(SPACE_ENTITY_TYPE, nil)
}
