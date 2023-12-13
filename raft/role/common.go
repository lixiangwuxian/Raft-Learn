package role

import (
	"lxtend.com/m/constants"
	"lxtend.com/m/packages"
	"lxtend.com/m/structs"
)

type Role interface {
	OnMsg(packet packages.Packet, inform *structs.InformAndHandler)
	Init(inform *structs.InformAndHandler, changeCallback func(constants.State))
	Clear()
}
