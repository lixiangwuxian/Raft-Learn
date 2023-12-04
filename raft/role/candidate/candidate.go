package candidate

import (
	"time"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var candidateTimeout *timeout.TimerTrigger
var roleCallback func(constants.State)

type Candidate struct {
	tickets int
}

func (c *Candidate) OnMsg(packet adapter.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntriesReply {
		return //ignore, this should be handled by leader
	} else if packet.TypeOfMsg == constants.RequestVoteReply {
		msgData := adapter.ParseRequestVoteReply(packet.Data)
		if msgData.Agree {
			c.tickets++
			if c.tickets >= len(inform.KnownNodes)/2+1 {
				c.Clear()
				roleCallback(constants.Leader)
			}
		} else {
			if msgData.MyTerm > inform.CurrentTerm {
				c.Clear()
				roleCallback(constants.Follower)
			}
		}
	} else if packet.TypeOfMsg == constants.AppendEntries {
		msgData := adapter.ParseAppendEntries(packet.Data)
		if msgData.Term >= inform.CurrentTerm {
			inform.CurrentTerm = msgData.Term
			inform.KnownLeader = packet.SourceAddr
			c.Clear()
			roleCallback(constants.Follower)
		} else {
			inform.Sender.AppendEntriesReply(packet.SourceAddr, adapter.AppendEntriesReply{
				Term:         inform.CurrentTerm,
				Success:      false,
				CurrentIndex: inform.Store.LastIndex(),
			})
		}
	} else if packet.TypeOfMsg == constants.RequestVote {
		msgData := adapter.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			inform.Sender.RequestVoteReply(packet.SourceAddr, adapter.RequestVoteReply{Agree: true, MyTerm: inform.CurrentTerm}, inform.CurrentTerm)
			c.Clear()
			roleCallback(constants.Follower)
		}
	}
}

func (c *Candidate) Init(inform *structs.InformAndHandler, changeCallback func(constants.State)) {
	candidateTimeout = timeout.NewTimerControl(time.Duration(inform.CandidateTimeout) * time.Millisecond)
	c.tickets = 1
	roleCallback = changeCallback
	candidateTimeout.Start(func() {
		c.Clear()
		roleCallback(constants.Follower)
	})
}

func (c *Candidate) Clear() {
	if candidateTimeout != nil {
		candidateTimeout.Stop()
		candidateTimeout = nil
	}
	c.tickets = 0
}
