package main

import (
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

func cacheAction(data *Action) bool {
	if data_to_submit != nil {
		return false
	}
	data_to_submit = data
	logger.Info("Cached data from leader, waiting for commit")
	return true
}

func commitAction() {
	if data_to_submit.Type == Set {
		kvStore.Set(data_to_submit.Key, data_to_submit.Value)
	}
	if data_to_submit.Type == Delete {
		kvStore.Delete(data_to_submit.Key)
	}
	data_to_submit = nil
	logger.Info("Data committed")
	//no need to implement get
}

func voteToCan(peerData RequestVote) {
	//send vote to candidate
	//check if can vote
	if inform.term == lastVotedTerm {
		adapter.VoteTo(inform.term, false)
	} else if inform.term < peerData.LastLogTerm {
		//vote
		adapter.VoteTo(inform.term, true)
		lastVotedTerm = inform.term
	} else if inform.term == peerData.LastLogTerm && logStore.Len() <= peerData.LastLogIndex {
		adapter.VoteTo(inform.term, true)
		lastVotedTerm = inform.term
	} else {
		adapter.VoteTo(inform.term, false)
	}
}

func handleHeartbeat(peerData HeartBeat) {
	gotHeartBeat = true
	inform.knownLeader = peerData.MyIP
}
