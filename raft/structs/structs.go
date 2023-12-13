package structs

import (
	"lxtend.com/m/adapter"
	"lxtend.com/m/store"
)

type InformAndHandler struct {
	KnownLeader      string
	KnownNodes       []string
	FollowerTimeout  int //多长时间从follower转为leader，使用时需要加一个随机数
	CandidateTimeout int //多长时间选举超时
	Sender           adapter.Sender
	MyAddr           string //ip&port
	Volatile
	Persistent
}

type Volatile struct {
	CommitIndex int
	lastApplied int
}

type Persistent struct {
	CurrentTerm int
	VotedFor    string
	Store       store.InMemoryLogStore
}
