package follower

import (
	"math/rand"
	"time"

	"lxtend.com/m/constants"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var followerTimeout *timeout.TimerTrigger
var roleChangeCallback func(constants.State)

type Follower struct {
	// LogIndex
}

func (f *Follower) OnMsg(packet packages.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntries {
		followerTimeout.Reset()
		data := packages.ParseAppendEntries(packet.Data)
		if data.Term >= inform.CurrentTerm {
			inform.CurrentTerm = data.Term
			inform.KnownLeader = packet.SourceAddr
		} else {
			return
		}
		logger.Glogger.Info("receive append entries from %s", packet.SourceAddr)
		if data.Entries != nil {
			var commands string
			for _, command := range data.Entries {
				commands += command.Command + "\n"
			}
			logger.Glogger.Info("receive append entries from %s, commands: %s", packet.SourceAddr, commands)
		}
		if inform.Store.LastTerm() == data.PrevLogTerm && inform.Store.LastIndex() == data.PrevLogIndex {
			logger.Glogger.Info("start append, current index %d, current term %d, leader PrevLogIndex %d, leader PrevLogTerm %d", inform.Store.LastIndex(), inform.Store.LastTerm(), data.PrevLogIndex, data.PrevLogTerm)
			inform.Store.Appends(data.Entries)
			inform.Sender.AppendEntriesReply(packet.SourceAddr,
				packages.AppendEntriesReply{Term: inform.CurrentTerm,
					Success:      true,
					CurrentIndex: inform.Store.LastIndex(),
				})
			logger.Glogger.Info("append done, current index %d, current term %d", inform.Store.LastIndex(), inform.Store.LastTerm())
		} else {
			logger.Glogger.Info("append entries failed, current index %d, current term %d, leader PrevLogIndex %d, leader PrevLogTerm %d", inform.Store.LastIndex(), inform.Store.LastTerm(), data.PrevLogIndex, data.PrevLogTerm)
			inform.Sender.AppendEntriesReply(packet.SourceAddr,
				packages.AppendEntriesReply{Term: inform.CurrentTerm,
					Success:      false,
					CurrentIndex: inform.Store.LastIndex(),
				})
			return
		}
		inform.Sender.AppendEntriesReply(packet.SourceAddr,
			packages.AppendEntriesReply{Term: inform.CurrentTerm,
				Success:      true,
				CurrentIndex: inform.Store.LastIndex(),
			})
	} else if packet.TypeOfMsg == constants.AppendEntriesReply {
		return
	} else if packet.TypeOfMsg == constants.RequestVote {
		logger.Glogger.Info("receive vote request from %s", packet.SourceAddr)
		if packet.Term > inform.CurrentTerm {
			followerTimeout.Reset()
			inform.Sender.RequestVoteReply(packet.SourceAddr,
				packages.RequestVoteReply{
					Agree:  true,
					MyTerm: inform.CurrentTerm,
				},
				inform.CurrentTerm)
			inform.CurrentTerm = packet.Term
			logger.Glogger.Info("voted for %s", packet.SourceAddr)
		} else {
			followerTimeout.Reset()
			inform.Sender.RequestVoteReply(packet.SourceAddr,
				packages.RequestVoteReply{
					Agree:  false,
					MyTerm: inform.CurrentTerm,
				},
				inform.CurrentTerm)
			logger.Glogger.Info("reject vote for %s, peer term is %d, my term is %d", packet.SourceAddr, packet.Term, inform.CurrentTerm)
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
