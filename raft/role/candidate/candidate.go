package candidate

import (
	"time"

	"lxtend.com/m/constants"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var roleCallback func(constants.State)

type Candidate struct {
	tickets          int
	candidateTimeout *timeout.TimerTrigger
}

func (c *Candidate) OnMsg(packet packages.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntriesReply {
		return //ignore, this should be handled by leader
	} else if packet.TypeOfMsg == constants.RequestVoteReply {
		msgData := packages.ParseRequestVoteReply(packet.Data)
		if msgData.Agree {
			logger.Glogger.Info("receive vote, now tickets is %d", c.tickets)
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
		msgData := packages.ParseAppendEntries(packet.Data)
		if msgData.Term >= inform.CurrentTerm {
			inform.CurrentTerm = msgData.Term
			inform.KnownLeader = packet.SourceAddr
			c.Clear()
			roleCallback(constants.Follower)
		} else {
			inform.Sender.AppendEntriesReply(packet.SourceAddr, packages.AppendEntriesReply{
				Term:         inform.CurrentTerm,
				Success:      false,
				CurrentIndex: inform.Store.LastIndex(),
			})
		}
	} else if packet.TypeOfMsg == constants.RequestVote {
		msgData := packages.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			inform.Sender.RequestVoteReply(packet.SourceAddr, packages.RequestVoteReply{Agree: true, MyTerm: inform.CurrentTerm}, inform.CurrentTerm)
			c.Clear()
			roleCallback(constants.Follower)
		}
	}
}

func (c *Candidate) Init(inform *structs.InformAndHandler, changeCallback func(constants.State)) {
	c.candidateTimeout = timeout.NewTimerControl(time.Duration(inform.CandidateTimeout) * time.Millisecond)
	inform.CurrentTerm++
	c.tickets = 1
	roleCallback = changeCallback
	for _, peer := range inform.KnownNodes {
		if peer != inform.MyAddr {
			inform.Sender.RequestVote(peer, packages.RequestVote{LastLogIndex: inform.Store.LastIndex(), LastLogTerm: inform.Store.LastTerm()}, inform.CurrentTerm)
			logger.Glogger.Info("sent request vote to %s", peer)
		}
	}
	c.candidateTimeout.Start(func() {
		c.Clear()
		roleCallback(constants.Follower)
	})
}

func (c *Candidate) Clear() {
	if c.candidateTimeout != nil {
		c.candidateTimeout.Stop()
		c.candidateTimeout = nil
	}
	c.tickets = 0
}
