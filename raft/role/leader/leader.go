package leader

import (
	"time"

	"lxtend.com/m/constants"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var roleCallback func(constants.State)

type Leader struct {
	nextIndex        map[string]int
	matchIndex       map[string]int
	heartbeatTrigger *timeout.TimerTrigger
	inform           *structs.InformAndHandler
}

func (l *Leader) OnMsg(packet packages.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntriesReply {
		msgData := packages.ParseAppendEntriesReply(packet.Data)
		logger.Glogger.Info("receive append entries reply from %s, success: %t, current index: %d", packet.SourceAddr, msgData.Success, msgData.CurrentIndex)
		if msgData.Term > inform.CurrentTerm {
			l.Clear()
			roleCallback(constants.Follower)
		}
		if !msgData.Success {
			l.matchIndex[packet.SourceAddr] = msgData.CurrentIndex
		} else {
			l.matchIndex[packet.SourceAddr] = msgData.CurrentIndex
			// l.nextIndex[packet.SourceAddr] = msgData.CurrentIndex + 1
			l.updateCommitIndex(inform)
		}
	} else if packet.TypeOfMsg == constants.RequestVoteReply {
		return //ignore, this should be handled by candidate
	} else if packet.TypeOfMsg == constants.AppendEntries {
		msgData := packages.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			l.Clear()
			inform.KnownLeader = packet.SourceAddr
			roleCallback(constants.Follower)
		}
	} else if packet.TypeOfMsg == constants.RequestVote {
		msgData := packages.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			inform.Sender.RequestVoteReply(packet.SourceAddr, packages.RequestVoteReply{Agree: true, MyTerm: inform.CurrentTerm}, inform.CurrentTerm)
			l.Clear()
			roleCallback(constants.Follower)
		}
	}
}

func (l *Leader) Init(inform *structs.InformAndHandler, changeCallback func(constants.State)) {
	roleCallback = changeCallback
	l.inform = inform
	l.heartbeatTrigger = timeout.NewTimerControl(time.Millisecond * time.Duration(inform.FollowerTimeout) / 2)
	l.heartbeatTrigger.StartIntervalTask(l.heartBeat)
	// l.nextIndex = make(map[string]int)
	l.matchIndex = make(map[string]int)
}

func (l *Leader) Clear() {
	// l.nextIndex = nil
	l.matchIndex = nil
}

func (l *Leader) heartBeat() {
	for _, addr := range l.inform.KnownNodes {
		if l.inform.Store.GetSince(l.matchIndex[addr]) != nil {
			var commands string
			for _, command := range l.inform.Store.GetSince(l.matchIndex[addr]) {
				commands += command.Command + "\n"
			}
			logger.Glogger.Info("send heartbeat to %s, commands: %s", addr, commands)
		}
		logger.Glogger.Info("last index: %d, match index: %d", l.inform.Store.LastIndex(), l.matchIndex[addr])
		l.inform.Sender.AppendEntries(addr,
			packages.AppendEntries{
				Term:         l.inform.CurrentTerm,
				PrevLogIndex: l.matchIndex[addr],
				PrevLogTerm:  l.inform.Store.Get(l.matchIndex[addr] - 1).Term,
				Entries:      l.inform.Store.GetSince(l.matchIndex[addr]),
				LeaderCommit: l.inform.CommitIndex,
			})
	}
}

func (l *Leader) updateCommitIndex(inform *structs.InformAndHandler) {
	for i := len(inform.Store.GetSince(0)) - 1; i > inform.CommitIndex; i-- {
		count := 1 // 包括 Leader 自己
		for _, addr := range inform.KnownNodes {
			if l.matchIndex[addr] >= i {
				count++
			}
		}
		if count > len(inform.KnownNodes)/2 {
			inform.CommitIndex = i
			break
		}
	}
}
