package role

import (
	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
)

type Role interface {
	OnMsg(packet adapter.Packet, inform *structs.Inform)
	Init(inform *structs.Inform, changeCallback func(constants.State))
	Clear()
}
