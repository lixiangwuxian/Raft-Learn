package leader

import (
	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
)

type Leader struct {
}

func (l Leader) OnMsg(packet adapter.Packet, inform *structs.Inform) constants.State {
	return constants.Leader
}

func (l Leader) Init(inform *structs.Inform, changeCallback func(constants.State)) {

}

func (l Leader) Clear() {

}
