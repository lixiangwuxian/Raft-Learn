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

func (f *Follower) OnMsg(packet adapter.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntries {
		followerTimeout.Reset()
		data := adapter.ParseAppendEntries(packet.Data)
		if data.Term >= inform.CurrentTerm {
			inform.CurrentTerm = data.Term
			inform.KnownLeader = packet.SourceAddr
		} else {
			return
		}
		if inform.Store.LastTerm() == data.PrevLogTerm && inform.Store.LastIndex() == data.PrevLogIndex {
			inform.Store.Append(data.Entries)
			inform.Sender.AppendEntriesReply(packet.SourceAddr,
				adapter.AppendEntriesReply{Term: inform.CurrentTerm,
					Success:      true,
					CurrentIndex: inform.Store.LastIndex(),
				})
		} else {
			inform.Sender.AppendEntriesReply(packet.SourceAddr,
				adapter.AppendEntriesReply{Term: inform.CurrentTerm,
					Success:      false,
					CurrentIndex: inform.Store.LastIndex(),
				})
			return
		}
		inform.Sender.AppendEntriesReply(packet.SourceAddr,
			adapter.AppendEntriesReply{Term: inform.CurrentTerm,
				Success:      true,
				CurrentIndex: inform.Store.LastIndex(),
			})
	} else if packet.TypeOfMsg == constants.AppendEntriesReply {
		return
	} else if packet.TypeOfMsg == constants.RequestVote {
		if packet.Term > inform.CurrentTerm {
			followerTimeout.Reset()
			inform.Sender.RequestVoteReply(packet.SourceAddr,
				adapter.RequestVoteReply{
					Agree:  true,
					MyTerm: inform.CurrentTerm,
				},
				inform.CurrentTerm)
			inform.CurrentTerm = packet.Term
		} else {
			followerTimeout.Reset()
			inform.Sender.RequestVoteReply(packet.SourceAddr,
				adapter.RequestVoteReply{
					Agree:  false,
					MyTerm: inform.CurrentTerm,
				},
				inform.CurrentTerm)
		}
	} else if packet.TypeOfMsg == constants.RequestVoteReply {
		return
	}
}

func (f *Follower) Init(inform *structs.InformAndHandler, changeCallback func(constants.State)) {
	followerTimeout = timeout.NewTimerControl(time.Duration(inform.FollowerTimeout+rand.Intn(200)) * time.Millisecond)
	roleChangeCallback = changeCallback
	followerTimeout.Start(func() {
		roleChangeCallback(constants.Candidate)
	})
}

func (f *Follower) Clear() {
	if followerTimeout != nil {
		followerTimeout.Stop()
		followerTimeout = nil
	}
}
