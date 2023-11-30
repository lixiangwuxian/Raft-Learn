package candidate

import (
	"time"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/structs"
	"lxtend.com/m/timeout"
)

var candidateTimeout *timeout.TimerTrigger
var msgCallback func(adapter.Packet)

type Candidate struct {
}

func (c Candidate) OnMsg(packet adapter.Packet, inform *structs.Inform) constants.State {
	return constants.Candidate
}

func (c Candidate) Init(inform *structs.Inform, changeCallback func(constants.State)) {
	candidateTimeout = timeout.NewTimerControl(time.Duration(inform.CandidateTimeout) * time.Millisecond)
}

func (c Candidate) Clear() {

}

// var canState CanState = CanState{0}

// func candidateTimeout() {
// 	time.Sleep(time.Duration(inform.candidateTimeout) * time.Millisecond)
// 	transToFollower()
// }

// func voteMe() {
// 	persist_inform.CurrentTerm++
// 	adapter.AskForVote(persist_inform.CurrentTerm)
// }

// func handleVoteTo(peerData RequestVoteReply) {
// 	if peerData.Agree {
// 		canState.vote++
// 	} else {
// 		if peerData.MyTerm > persist_inform.CurrentTerm {
// 			transToFollower()
// 		}
// 	}
// 	if canState.vote >= inform.totalNodes/2+1 {
// 		transToLeader()
// 	}
// }
