package follower

import (
	"math/rand"
	"time"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var followerTimeout *timeout.TimerTrigger
var roleChangeCallback func(constants.State)

type Follower struct {
}

func (f Follower) OnMsg(packet adapter.Packet, inform *structs.Inform) constants.State {
	return constants.Follower
}

func (f Follower) Init(inform *structs.Inform, changeCallback func(constants.State)) {
	followerTimeout = timeout.NewTimerControl(time.Duration(inform.FollowerTimeout+rand.Intn(200)) * time.Millisecond)
	roleChangeCallback = changeCallback
	followerTimeout.Start(func() {
		roleChangeCallback(constants.Candidate)
	})
}

func (f Follower) Clear() {
	if followerTimeout != nil {
		followerTimeout.Stop()
		followerTimeout = nil
	}
}
