package main

type CanState struct {
	// term int
	vote int
}

var canState CanState

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
