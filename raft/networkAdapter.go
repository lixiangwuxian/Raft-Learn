package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

type Adapter struct {
	peers       map[int]string
	onlinePeers map[int]string
	listenConn  net.Conn
}

type Packet struct {
	TypeOfMsg int
	Term      int
	Data      json.RawMessage
}

// type HeartBeat struct {
// 	From   int
// 	MyIP   string
// 	MyTerm int
// }

type RequestVote struct {
	From         int
	LastLogIndex int
	LastLogTerm  int
}
type RequestVoteReply struct {
	From   int
	Agree  bool
	MyTerm int
}

type AppendEntries struct {
	Term         int
	leaderId     int
	prevLogIndex int
	prevLogTerm  int
	entries      []byte
	leaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

// type Cache struct {
// 	From       int
// 	ActionToDo Action
// }

// type CacheReply struct {
// 	From int
// 	Okay bool
// }

// type Commit struct {
// 	From  int
// 	Stage int
// }

// type CommitReply struct {
// 	From int
// 	Okay bool
// }

type IMOnline struct {
	From int
	MyIP string
}

type IMHere struct {
	From int
	MyIP string
}

const (
	// heartbeat int = iota
	data_from_client = iota
	request_vote
	// vote_to
	request_vote_reply
	// cache_data
	// cache_reply
	// commit
	// commit_reply
	append_entries
	append_entries_reply
	im_online
	im_here
)

func (a *Adapter) Init() {
	a.peers = make(map[int]string)
	// a.maskIP = maskIP
	listener, _ := kcp.Listen("0.0.0.0:18230")
	a.listenConn, _ = listener.Accept()
	go a.ListenLoop()
}

func (a *Adapter) ListenLoop() {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	buffer := make([]byte, 1024)
	for {
		randLantency := rr.Intn(100)
		a.listenConn.SetReadDeadline(time.Now().Add(time.Millisecond*200 + time.Duration(randLantency)*time.Millisecond))
		_, err := a.listenConn.Read(buffer)
		for i := 0; i < 1024; i++ {
			fmt.Println(buffer[i])
			if buffer[i] == 0 {
				buffer = buffer[:i]
				break
			}
		}
		if err != nil {
			//timeout
			if inform.state == Leader {
				continue
			} else if inform.state == Follower {
				transToCandidate()
			} else if inform.state == Candidate {
				continue
			}
		}
		packet := parseData(buffer)
		switch packet.TypeOfMsg {
		// case heartbeat:
		// 	var data HeartBeat
		// 	json.Unmarshal(packet.Data, &data)
		// 	if inform.state == Follower {
		// 		handleHeartbeat(data)
		// 	}
		// 	if inform.state == Candidate {
		// 		transToFollower()
		// 	}
		// 	if inform.state == Leader {
		// 		needKeepLeader(data)
		// 	}
		case data_from_client:
			var data string
			json.Unmarshal(packet.Data, &data)
			if inform.state == Leader {
				dataFromClient(data)
			}
		case request_vote:
			var data RequestVote
			json.Unmarshal(packet.Data, &data)
			voteToCan(data)
		case request_vote_reply:
			var data RequestVoteReply
			json.Unmarshal(packet.Data, &data)
			if inform.state == Candidate {
				handleVoteTo(data)
			}
		case append_entries:
			var data AppendEntries
			json.Unmarshal(packet.Data, &data)
			if inform.state == Follower {
				handleAppendEntries(data)
			}
		case append_entries_reply:
			var data AppendEntriesReply
			json.Unmarshal(packet.Data, &data)
			if inform.state == Leader {
				handleAppendEntriesReply()
			}
		// case cache_data:
		// 	var data Cache
		// 	json.Unmarshal(packet.Data, &data)
		// 	if inform.state == Follower {
		// 		if data.ActionToDo.Index >= logStore.PeekLastIndex() {
		// 			cacheAction(&data.ActionToDo)
		// 		}
		// 	}
		// case cache_reply:
		// 	var data CacheReply
		// 	json.Unmarshal(packet.Data, &data)
		// 	if inform.state == Leader {
		// 		if data.Okay {
		// 			addCacheReply()
		// 		}
		// 	}
		// case commit:
		// 	var data Commit
		// 	json.Unmarshal(packet.Data, &data)
		// 	if inform.state == Follower {
		// 		commitAction()
		// 	}
		// case commit_reply:
		// 	continue
		case im_online:
			var data IMOnline
			json.Unmarshal(packet.Data, &data)
			a.onlinePeers[data.From] = data.MyIP
			a.SendImHere(data.From)
		case im_here:
			var data IMHere
			json.Unmarshal(packet.Data, &data)
			a.onlinePeers[data.From] = data.MyIP
		}
	}
}

func parseData(data []byte) Packet {
	var packet Packet
	json.Unmarshal(data, &packet)
	return packet
}

func (a *Adapter) BroadcastData(data []byte, typeOfMsg int) {
	fullData, _ := json.Marshal(Packet{typeOfMsg, persist_inform.CurrentTerm, data})
	for _, peer := range a.onlinePeers {
		conn, _ := kcp.Dial(peer)
		conn.Write(fullData)
	}
}

func (a *Adapter) FirstBroadcastData(data []byte, typeOfMsg int) {
	fullData, _ := json.Marshal(Packet{typeOfMsg, persist_inform.CurrentTerm, data})
	for _, peer := range a.peers {
		conn, _ := kcp.Dial(peer)
		conn.Write(fullData)
	}
}

func (a *Adapter) WriteDataTo(peer int, data []byte, typeOfMsg int) {
	conn, _ := kcp.Dial(a.onlinePeers[peer])
	fullData, _ := json.Marshal(Packet{typeOfMsg, persist_inform.CurrentTerm, data})
	conn.Write(fullData)
}
func (a *Adapter) AskForVote(term int) { //use kcp broadcast
	pkg := RequestVote{inform.whoAmI, logStore.Len(), logStore.PeekLastTerm()}
	data, _ := json.Marshal(pkg)
	a.BroadcastData(data, request_vote)
}

func (a *Adapter) VoteTo(peer int, agree bool) {
	pkg := RequestVoteReply{inform.whoAmI, agree, persist_inform.CurrentTerm}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data, request_vote_reply)
}

func (a *Adapter) SendAppendEntries(entries []byte) {
	pkg := AppendEntries{persist_inform.CurrentTerm, inform.whoAmI, logStore.PeekLastIndex(), logStore.PeekLastTerm(), entries, logStore.PeekLastIndex()}
	data, _ := json.Marshal(pkg)
	a.BroadcastData(data, append_entries)
}

func (a *Adapter) SendAppendEntriesReply(peer int, success bool) {
	pkg := AppendEntriesReply{persist_inform.CurrentTerm, success}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data, append_entries_reply)
}

// func (a *Adapter) SendHeartbeat() {
// 	pkg := HeartBeat{inform.whoAmI, inform.myIP, persist_inform.CurrentTerm}
// 	data, _ := json.Marshal(pkg)
// 	a.BroadcastData(data)
// }

// func (a *Adapter) SendLog(action Action) {
// 	pkg := Cache{inform.whoAmI, action}
// 	data, _ := json.Marshal(pkg)
// 	a.BroadcastData(data)
// }

// func (a *Adapter) SendCommitTo(peer int, logIndex int) {
// 	pkg := Commit{inform.whoAmI, logIndex}
// 	data, _ := json.Marshal(pkg)
// 	a.WriteDataTo(peer, data)
// }

// func (a *Adapter) SendLogTo(peer int, action Action, index int) {
// 	pkg := Cache{inform.whoAmI, action}
// 	data, _ := json.Marshal(pkg)
// 	a.WriteDataTo(peer, data)
// }

func (a *Adapter) BroadOnline() {
	pkg := IMOnline{inform.whoAmI, inform.myIP}
	data, _ := json.Marshal(pkg)
	a.FirstBroadcastData(data, im_online)
}

func (a *Adapter) SendImHere(peer int) {
	pkg := IMHere{inform.whoAmI, inform.myIP}
	data, _ := json.Marshal(pkg)
	a.WriteDataTo(peer, data, im_here)
}

var adapter Adapter
