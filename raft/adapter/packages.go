package adapter

import (
	"encoding/json"

	"lxtend.com/m/constants"
)

type Packet struct {
	TypeOfMsg constants.PackegeType
	Term      int
	Data      json.RawMessage
}

type AppendEntries struct {
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

type AppendEntriesReply struct {
	Term    int  `json:"term"`
	Success bool `json:"success"`
}

type RequestVote struct {
	From         int `json:"from"`
	LastLogIndex int `json:"lastLogIndex"`
	LastLogTerm  int `json:"lastLogTerm"`
}
type RequestVoteReply struct {
	From   int  `json:"from"`
	Agree  bool `json:"agree"`
	MyTerm int  `json:"myTerm"`
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
