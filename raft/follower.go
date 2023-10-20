package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

// const (
// 	waitForData int = iota
// 	waitForCommit
// )

// var leaderTimeout int

var lastVotedTerm int

var gotHeartBeat bool

func followerLoop() {
	for inform.state == Follower {
		time.Sleep(time.Duration(inform.leaderTimeout)*time.Millisecond + time.Duration(rand.Intn(100)+50)*time.Millisecond)
		if gotHeartBeat {
			gotHeartBeat = false
			continue
		} else {
			transToCandidate()
		}
	}
}

func voteToCan(peerData RequestVote) {
	//send vote to candidate
	//check if can vote
	if persist_inform.CurrentTerm == lastVotedTerm {
		adapter.VoteTo(persist_inform.CurrentTerm, false)
	} else if persist_inform.CurrentTerm < peerData.LastLogTerm {
		//vote
		adapter.VoteTo(persist_inform.CurrentTerm, true)
		lastVotedTerm = persist_inform.CurrentTerm
	} else if persist_inform.CurrentTerm == peerData.LastLogTerm && logStore.Len() <= peerData.LastLogIndex {
		adapter.VoteTo(persist_inform.CurrentTerm, true)
		lastVotedTerm = persist_inform.CurrentTerm
	} else {
		adapter.VoteTo(persist_inform.CurrentTerm, false)
	}
}

func handleAppendEntries(peerData AppendEntries) {
	gotHeartBeat = true
	inform.knownLeader = peerData.leaderId
	var leaderEntries = make([]Action, 0)
	json.Unmarshal(peerData.entries, &leaderEntries)
	clearDiffEntries(leaderEntries, peerData.prevLogIndex)
	if peerData.leaderCommit > inform.commitIndex {
		logStore.CommitIndex = min(logStore.CommitIndex, logStore.PeekLastIndex())
	}
	if peerData.Term < persist_inform.CurrentTerm {
		go adapter.SendAppendEntriesReply(inform.whoAmI, false)
	} else if checkEntries(peerData.prevLogIndex, peerData.prevLogTerm) {
		go adapter.SendAppendEntriesReply(inform.whoAmI, true)
	} else {
		go adapter.SendAppendEntriesReply(inform.whoAmI, false)
	}
}

func checkEntries(prevLogIndex int, prevLogTerm int) bool {
	return logStore.Get(prevLogIndex).Term == prevLogTerm
}

func clearDiffEntries(entries []Action, prevLogIndex int) {
	for index, entry := range entries {
		if index+prevLogIndex+1 > logStore.Len() {
			logStore.Append(entry)
		} else if logStore.Get(index+prevLogIndex+1).Term != entry.Term {
			logStore.Set(index, entry)
		}
	}
}
