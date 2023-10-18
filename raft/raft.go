package main

import (
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
	inform.state = Leader
	go leaderLoop()
}

func transToCandidate() {
	inform.state = Candidate
}

func transToFollower() {
	inform.state = Follower
	go followerLoop()
}

// func startListenFromBroadcast(inform *Inform, port int) {
// 	addr := fmt.Sprintf(":%d", port)
// 	conn, _ := net.ListenPacket("udp", addr)
// 	defer conn.Close()
// 	for {
// 		var inform Inform
// 		data := make([]byte, 1024)
// 		n, addr, _ := conn.ReadFrom(data)
// 		json.Unmarshal(data[:n], &inform)
// 		logger.Info("received from %s", addr.String())
// 	}
// }

func main() {
	myIP := "10.22.34.206:18230"
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
