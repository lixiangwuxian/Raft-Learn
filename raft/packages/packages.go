package packages

import (
	"encoding/json"

	"lxtend.com/m/constants"
)

type Packet struct {
	TypeOfMsg  constants.PackegeType `json:"typeOfMsg"`
	SourceAddr string                `json:"sourceAddr"`
	Term       int                   `json:"term"`
	Data       json.RawMessage       `json:"data"`
}

type AppendEntries struct {
	Term         int     `json:"term"`
	PrevLogIndex int     `json:"prevLogIndex"`
	PrevLogTerm  int     `json:"prevLogTerm"`
	Entries      []Entry `json:"entries"`
	LeaderCommit int     `json:"leaderCommit"`
}

type AppendEntriesReply struct {
	Term         int  `json:"term"`
	Success      bool `json:"success"`
	CurrentIndex int  `json:"currentIndex"`
}

type Entry struct {
	Term    int    `json:"term"`
	Command string `json:"command"`
}

type RequestVote struct {
	LastLogIndex int `json:"lastLogIndex"`
	LastLogTerm  int `json:"lastLogTerm"`
}
type RequestVoteReply struct {
	Agree  bool `json:"agree"`
	MyTerm int  `json:"myTerm"`
}

type UserRequest struct {
	Command string `json:"command"`
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

func ParseUserRequest(data json.RawMessage) UserRequest {
	var userRequest UserRequest
	json.Unmarshal(data, &userRequest)
	return userRequest
}
