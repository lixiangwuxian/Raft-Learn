package leader

import (
	"time"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
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

func (l *Leader) OnMsg(packet adapter.Packet, inform *structs.InformAndHandler) {
	if packet.TypeOfMsg == constants.AppendEntriesReply {
		msgData := adapter.ParseAppendEntriesReply(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			l.Clear()
			roleCallback(constants.Follower)
		}
		if !msgData.Success {
			l.nextIndex[packet.SourceAddr]--
		} else {
			l.matchIndex[packet.SourceAddr] = msgData.CurrentIndex
			l.nextIndex[packet.SourceAddr] = msgData.CurrentIndex + 1

			// 检查是否可以增加 CommitIndex
			l.updateCommitIndex(inform)
		}
	} else if packet.TypeOfMsg == constants.RequestVoteReply {
		return //ignore, this should be handled by candidate
	} else if packet.TypeOfMsg == constants.AppendEntries {
		msgData := adapter.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			l.Clear()
			inform.KnownLeader = packet.SourceAddr
			roleCallback(constants.Follower)
		}
	} else if packet.TypeOfMsg == constants.RequestVote {
		msgData := adapter.ParseAppendEntries(packet.Data)
		if msgData.Term > inform.CurrentTerm {
			inform.Sender.RequestVoteReply(packet.SourceAddr, adapter.RequestVoteReply{Agree: true, MyTerm: inform.CurrentTerm}, inform.CurrentTerm)
			l.Clear()
			roleCallback(constants.Follower)
		}
	}
}

func (l *Leader) Init(inform *structs.InformAndHandler, changeCallback func(constants.State)) {
	roleCallback = changeCallback
	l.heartbeatTrigger = timeout.NewTimerControl(time.Millisecond * time.Duration(inform.FollowerTimeout) / 4)
	l.heartbeatTrigger.Start(l.heartBeat)
	l.nextIndex = make(map[string]int)
	l.matchIndex = make(map[string]int)
}

func (l *Leader) Clear() {
	l.nextIndex = nil
	l.matchIndex = nil
}

func (l *Leader) heartBeat() {
	for _, addr := range l.inform.KnownNodes {
		l.inform.Sender.AppendEntries(addr,
			adapter.AppendEntries{
				Term:         l.inform.CurrentTerm,
				PrevLogIndex: l.inform.Store.LastIndex(),
				PrevLogTerm:  l.inform.Store.LastTerm(),
				Entries:      l.inform.Store.GetSince(l.nextIndex[addr]),
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
