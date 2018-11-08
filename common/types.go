package common

import "../uuid"
const ENTITYID_LENGTH = uuid.UUID_LENGTH

type EntityID string

func GenEntityID() EntityID {
	return EntityID(uuid.GenUUID())
}
