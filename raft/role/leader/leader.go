package leader

import (
	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
)

var roleCallback func(constants.State)

type Leader struct {
}

func (l Leader) OnMsg(packet adapter.Packet, inform *structs.Inform) {
	if packet.TypeOfMsg == constants.AppendEntriesReply {
	}
}

func (l Leader) Init(inform *structs.Inform, changeCallback func(constants.State)) {
	roleCallback = changeCallback
}

func (l Leader) Clear() {
}
