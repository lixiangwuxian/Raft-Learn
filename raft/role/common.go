package role

import (
	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
)

type Role interface {
	OnMsg(packet adapter.Packet, inform *structs.InformAndHandler)
	Init(inform *structs.InformAndHandler, changeCallback func(constants.State))
	Clear()
}
