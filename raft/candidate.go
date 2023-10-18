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
	inform.term++
	adapter.AskForVote(inform.term)
}

func handleVoteTo(peerData VoteTo) {
	if peerData.Agree {
		canState.vote++
	} else {
		if peerData.MyTerm > inform.term {
			transToFollower()
		}
	}
	if canState.vote >= inform.totalNodes/2+1 {
		transToLeader()
	}
}
