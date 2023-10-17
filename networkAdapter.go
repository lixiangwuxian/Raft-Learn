package raft

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

type Adapter struct {
	peers      map[int]string
	whoAmI     int
	listenConn net.Conn
	maskIP     string
}

type Packet struct {
	typeOfMsg int
	term      int
	data      json.RawMessage
}

type HeartBeat struct {
	from   int
	isEcho bool
}

type RequestVote struct { //request vote, this is the message content
	// term         int
	// candidateId  int
	lastLogIndex int
	lastLogTerm  int
}
type VoteTo struct {
	from  int
	agree bool
}

type Log struct {
	from int
	log  string
}

type Commit struct {
	from  int
	stage int
}

type CommitReply struct {
	from int
	okay bool
}

type AnyBodyThere struct {
	from int
	myIP string
}

type LeaderIs struct {
	from     int
	leaderIP string
}

const (
	heartbeat int = iota
	requestVote
	voteTo
	log_data
	commit
	commitReply
	anyBodyThere
	leaderIs
)

// type LeaderPacket struct {
// 	typeOfMsg int
// 	term      int
// 	data      string
// }

// type CandidatePacket struct {
// 	typeOfMsg int
// 	term      int
// 	data      string
// }

func (this *Adapter) Init(maskIP string) {
	this.whoAmI = time.Now().Nanosecond()
	this.peers = make(map[int]string)
	this.maskIP = maskIP
	listener, _ := kcp.Listen("0.0.0.0:18230")
	this.listenConn, _ = listener.Accept()
}

func (this *Adapter) ListenLoop() {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer := make([]byte, 1024)
	for {
		randLantency := rr.Intn(100)
		this.listenConn.SetReadDeadline(time.Now().Add(time.Millisecond*200 + time.Duration(randLantency)*time.Millisecond))
		_, err := this.listenConn.Read(buffer)
		if err != nil {
			//timeout
			if inform.state == Leader {
				continue
			} else if inform.state == Follower {
				inform.state = Candidate
				inform.term++
				this.AskForVote(inform.term)
			} else if inform.state == Candidate {
				continue
			}
		}
		packet := parseData(buffer)
		switch packet.typeOfMsg {
		case heartbeat:
			var data HeartBeat
			json.Unmarshal(packet.data, data)
			if !data.isEcho {
				//send heartbeat to leader
				this.SendHeartbeatTo(packet.term, false)
			}
		case requestVote:
			var data RequestVote
			json.Unmarshal(packet.data, data)
			//check if can vote
			if inform.term < packet.term {
				//vote
				this.VoteTo(inform.term, true)
			} else if inform.term == packet.term && logStore.Len() <= data.lastLogIndex {
				this.VoteTo(inform.term, true)
			} else {
				this.VoteTo(inform.term, false)
			}
		case voteTo:
			var data VoteTo
			json.Unmarshal(packet.data, data)

		case log_data:
			continue
		case commit:
			continue
		case commitReply:
			continue
		case anyBodyThere:
			continue
		case leaderIs:
			continue
		}
	}
}

func parseData(data []byte) Packet {
	var packet Packet
	json.Unmarshal(data, &packet)
	return packet
}

func (this *Adapter) BroadcastData(data string, broadCastIP string) {
	conn, _ := kcp.Dial(broadCastIP) //change later
	conn.Write([]byte(data))
}

func (a *Adapter) AskForVote(term int) { //use udp broadcast
	conn, _ := kcp.Dial(a.maskIP)
	pkg := RequestVote{term, a.whoAmI, logStore.Len(), logStore.PeekLastTerm()}
	conn.Write(pkg)

}

func (a *Adapter) WriteDataTo(peer int, data []byte) {
	conn, _ := kcp.Dial(a.peers[peer])
	conn.Write(data)
}

func (a *Adapter) VoteTo(peer int, agree bool) {
	pkg := VoteTo{a.whoAmI, agree}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data)
}

func (a *Adapter) SendHeartbeatTo(peer int, isEcho bool) {
	pkg := HeartBeat{a.whoAmI, isEcho}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data)
}

func (a *Adapter) SendCommitTo(peer int, logIndex int) {
	pkg := Commit{a.whoAmI, logIndex}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data)
}

func (a *Adapter) SendLogTo(peer int, log string) {
	pkg := Log{a.whoAmI, log}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data)
}
