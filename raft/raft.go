package main

import (
	"os"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/role"
	"lxtend.com/m/role/candidate"
	"lxtend.com/m/role/follower"
	"lxtend.com/m/role/leader"
	"lxtend.com/m/store"
	"lxtend.com/m/structs"
)

var logger, _ = NewLogger("")

var inform *structs.InformAndHandler

var logStore *store.InMemoryLogStore

var logIndex int

var roleNow = constants.Follower

var roleMap = map[constants.State]role.Role{}

func initInform(followerTimeout int, candidateTimeout int, myIP string, totalNodes int) *structs.InformAndHandler { //init the inform
	if inform == nil {
		inform = new(structs.InformAndHandler)
		inform.CurrentTerm = 0
		inform.FollowerTimeout = followerTimeout
		inform.CandidateTimeout = candidateTimeout
		inform.KnownLeader = ""
		inform.VotedFor = ""
		inform.Sender = &adapter.KcpSender{}
		inform.Store = store.InMemoryLogStore{}
	}
	return inform
}

func onMsg(packet adapter.Packet) {
	roleMap[roleNow].OnMsg(packet, inform)
}

func changRole(state constants.State) {
	roleNow = state
	roleMap[roleNow].Init(inform, changRole)
}

func main() {
	myIP := os.Args[1]
	logIndex = 0
	initConf()
	logStore = new(store.InMemoryLogStore)
	inform = initInform(1000, 500, myIP, 3)
	netAdapter := adapter.InitAdapter(getPeers())
	roleMap[constants.Follower] = &follower.Follower{}
	roleMap[constants.Leader] = &leader.Leader{}
	roleMap[constants.Candidate] = &candidate.Candidate{}
	roleNow = constants.Follower
	roleMap[roleNow].Init(inform, changRole)
	netAdapter.ListenLoop(onMsg)
}
