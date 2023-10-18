package main

import (
	"os"
	"time"
)

var logger, _ = NewLogger("")

type State int

const (
	Follower State = iota
	Leader
	Candidate
)

type Inform struct { //inform the leader, use singletons
	term             int
	state            State
	knownLeader      string
	knownNodes       []int
	totalNodes       int
	whoAmI           int
	myIP             string
	leaderTimeout    int
	candidateTimeout int
}

var inform *Inform

var kvStore *KVStore

var logStore *LogStore

var logIndex int

var data_to_submit *Action

func initInform(leaderTimeout int, canTimeout int, myIP string, totalNodes int) *Inform { //init the inform
	if inform == nil {
		inform = new(Inform)
		inform.term = 0
		inform.state = Follower
		inform.knownLeader = ""
		inform.knownNodes = make([]int, 0)
		inform.totalNodes = totalNodes
		inform.whoAmI = time.Now().Nanosecond()
		inform.myIP = myIP
		inform.leaderTimeout = leaderTimeout
		inform.candidateTimeout = canTimeout
	}
	return inform
}

func broadcastMe() {
	adapter.BroadOnline()
}

func transToLeader() {
	logger.Info("From %d transfrom to leader", inform.state)
	inform.state = Leader
	logIndex = logStore.PeekLastIndex() + 1
	go leaderLoop()
}

func transToCandidate() {
	logger.Info("From %d transfrom to candidate", inform.state)
	inform.state = Candidate
	voteMe()
	go candidateTimeout()
}

func transToFollower() {
	logger.Info("From %d transfrom to follower", inform.state)
	inform.state = Follower
	go followerLoop()
}

func main() {
	myIP := os.Args[1]
	logIndex = 0
	initConf()
	kvStore = new(KVStore)
	logStore = new(LogStore)
	inform = initInform(1000, 500, myIP, 3)
	adapter.Init()
	adapter.peers = getPeers()
	broadcastMe()
	go adapter.ListenLoop()
}
