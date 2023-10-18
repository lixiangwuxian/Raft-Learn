package main

import "time"

type CanState struct {
	// term int
	vote int
}

var canState CanState = CanState{0}

func candidateTimeout() {
	time.Sleep(time.Duration(inform.candidateTimeout) * time.Millisecond)
	transToFollower()
}

func voteMe() {
	persist_inform.CurrentTerm++
	adapter.AskForVote(persist_inform.CurrentTerm)
}

func handleVoteTo(peerData RequestVoteReply) {
	if peerData.Agree {
		canState.vote++
	} else {
		if peerData.MyTerm > persist_inform.CurrentTerm {
			transToFollower()
		}
	}
	if canState.vote >= inform.totalNodes/2+1 {
		transToLeader()
	}
}
