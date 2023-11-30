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

var inform *structs.Inform

var logStore *store.InMemoryLogStore

var logIndex int

var roleNow = constants.Follower

var roleMap = map[constants.State]role.Role{}

func initInform(leaderTimeout int, canTimeout int, myIP string, totalNodes int) *structs.Inform { //init the inform
	if inform == nil {
		inform = new(structs.Inform)
		inform.CurrentTerm = 0
		inform.FollowerTimeout = leaderTimeout
		inform.CandidateTimeout = canTimeout
	}
	return inform
}

func onMsg(packet adapter.Packet) {
	roleNow = roleMap[roleNow].OnMsg(packet, inform)
}

func main() {
	myIP := os.Args[1]
	logIndex = 0
	initConf()
	logStore = new(store.InMemoryLogStore)
	inform = initInform(1000, 500, myIP, 3)
	netAdapter := adapter.InitAdapter(getPeers())
	roleMap[constants.Follower] = follower.Follower{}
	roleMap[constants.Leader] = leader.Leader{}
	roleMap[constants.Candidate] = candidate.Candidate{}
	roleNow = follower.OnMsg(adapter.Packet{}, inform)
	netAdapter.ListenLoop(onMsg)
}
