package raft

import (
	"encoding/json"
	"fmt"
	"net"
)

var logger, _ = NewLogger("")

type State int

const (
	Follower State = iota
	Leader
	Candidate
)

type Inform struct { //inform the leader, use singletons
	term        int
	state       State
	knownLeader int
	knownNodes  []int
	totalNodes  int
	// leaderTimeout    int
	// candidateTimeout int
}

var inform *Inform

var kvStore *KVStore

var logStore *LogStore

func initInform(leaderTimeout int, canTimeout int) *Inform { //init the inform
	if inform == nil {
		inform = new(Inform)
		inform.term = 0
		inform.state = Follower
		inform.knownLeader = -1
		inform.knownNodes = make([]int, 0)
		// inform.leaderTimeout = leaderTimeout
		// inform.candidateTimeout = canTimeout
	}
	return inform
}

type RequestVote struct { //request vote, this is the message content
	term         int
	candidateId  int
	lastLogIndex int
	lastLogTerm  int
}

func broadcastMe(myIP string, maskIP string, port int) { //broadcast node itself to other nodes by udp broadcast
	var maskAddr *net.IPAddr
	maskAddr, _ = net.ResolveIPAddr("ip", maskIP)
	maskUDPAddr := &net.UDPAddr{
		IP:   maskAddr.IP,
		Port: port,
	}
	conn, _ := net.DialUDP("udp", nil, maskUDPAddr)
	defer conn.Close()
	conn.Write([]byte(myIP))
	logger.Info("broadcasted")
}

func startListenFromBroadcast(inform *Inform, port int) {
	addr := fmt.Sprintf(":%d", port)
	conn, _ := net.ListenPacket("udp", addr)
	defer conn.Close()
	for {
		var inform Inform
		data := make([]byte, 1024)
		n, addr, _ := conn.ReadFrom(data)
		json.Unmarshal(data[:n], &inform)
		logger.Info("received from %s", addr.String())
	}
}

func main() {
	inform := initInform(1000, 500)
	myIP := "192.168.0.1"
	maskIP := "192.168.0.255"
	go startListenFromBroadcast(inform, 18082)
	for inform.knownLeader == -1 {
		broadcastMe(myIP, maskIP, 18082)
	}
}
