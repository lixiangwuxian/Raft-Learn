package main

import (
	"os"
	"strings"

	"lxtend.com/m/adapter"
	"lxtend.com/m/constants"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
	"lxtend.com/m/role"
	"lxtend.com/m/role/candidate"
	"lxtend.com/m/role/follower"
	"lxtend.com/m/role/leader"
	"lxtend.com/m/store"
	"lxtend.com/m/structs"
)

var inform *structs.InformAndHandler

var roleNow = constants.Follower

var roleMap = map[constants.State]role.Role{}

func initInform(followerTimeout int, candidateTimeout int, myAddr string, totalNodes int) *structs.InformAndHandler { //init the inform
	if inform == nil {
		inform = new(structs.InformAndHandler)
		inform.CurrentTerm = 0
		inform.FollowerTimeout = followerTimeout
		inform.CandidateTimeout = candidateTimeout
		inform.KnownLeader = ""
		inform.VotedFor = ""
		inform.Sender = &adapter.KcpSender{MyAddr: myAddr}
		inform.Store = store.InMemoryLogStore{}
		inform.MyAddr = myAddr
	}
	return inform
}

func onMsg(packet packages.Packet) {
	logger.Glogger.Info("receive msg from %s", packet.SourceAddr)
	roleMap[roleNow].OnMsg(packet, inform)
}

func onUserCommand(command string) {
	inform.Store.Append(packages.Entry{Term: inform.CurrentTerm, Command: command})
}

func changRole(state constants.State) {
	roleNow = state
	logger.Glogger.Info("change role to %s", state)
	roleMap[roleNow].Init(inform, changRole)
}

func main() {
	myAddr := os.Args[1]
	_, myPort := strings.Split(myAddr, ":")[0], strings.Split(myAddr, ":")[1]
	conf := initConf(os.Args[2])
	inform = initInform(20000, 10000, myAddr, 3)
	adapter.ListenHttp(myPort, onUserCommand, &inform.Store)
	netAdapter := adapter.InitAdapter(conf.Peers, myPort)
	inform.KnownNodes = conf.Peers
	roleMap[constants.Follower] = &follower.Follower{}
	roleMap[constants.Leader] = &leader.Leader{}
	roleMap[constants.Candidate] = &candidate.Candidate{}
	roleNow = constants.Follower
	roleMap[roleNow].Init(inform, changRole)
	netAdapter.ListenLoop(onMsg)
}
