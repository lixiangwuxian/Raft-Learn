package adapter

import "encoding/json"

type Packet struct {
	TypeOfMsg int
	Term      int
	Data      json.RawMessage
}

type AppendEntries struct { //type 1
	Term         int     `json:"term"`
	LeaderId     int     `json:"leaderId"`
	PrevLogIndex int     `json:"prevLogIndex"`
	PrevLogTerm  int     `json:"prevLogTerm"`
	Entries      []Entry `json:"entries"`
	LeaderCommit int     `json:"leaderCommit"`
}

type Entry struct {
	Term    int    `json:"term"`
	Command string `json:"command"`
}

type AppendEntriesReply struct { //type 2
	Term    int
	Success bool
}

type RequestVote struct { //type 3
	From         int
	LastLogIndex int
	LastLogTerm  int
}
type RequestVoteReply struct { //type 4
	From   int
	Agree  bool
	MyTerm int
}

func ParseAppendEntries(data json.RawMessage) AppendEntries {
	var appendEntries AppendEntries
	json.Unmarshal(data, &appendEntries)
	return appendEntries
}

func ParseAppendEntriesReply(data json.RawMessage) AppendEntriesReply {
	var appendEntriesReply AppendEntriesReply
	json.Unmarshal(data, &appendEntriesReply)
	return appendEntriesReply
}

func ParseRequestVote(data json.RawMessage) RequestVote {
	var requestVote RequestVote
	json.Unmarshal(data, &requestVote)
	return requestVote
}

func ParseRequestVoteReply(data json.RawMessage) RequestVoteReply {
	var requestVoteReply RequestVoteReply
	json.Unmarshal(data, &requestVoteReply)
	return requestVoteReply
}
