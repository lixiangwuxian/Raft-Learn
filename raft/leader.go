package main

import (
	"strings"
	"sync"
	"time"
)

var leader_mutex sync.Locker = &sync.Mutex{}
var good_client_count int = 0

func leaderLoop() {
	for inform.state == Leader {
		adapter.SendHeartbeat()
		time.Sleep(time.Duration(inform.leaderTimeout/2) * time.Millisecond)
	}
}

func dataFromClient(data string) {
	logIndex++
	action := parseAction(data)
	data_to_submit = &action
	adapter.SendLog(*data_to_submit)
}

func addCacheReply() {
	leader_mutex.Lock()
	good_client_count++
	if good_client_count >= inform.totalNodes/2+1 {
		commitAction()
	}
	leader_mutex.Unlock()
}

func needKeepLeader(data HeartBeat) {
	if data.MyTerm > inform.term {
		inform.term = data.MyTerm
		inform.knownLeader = data.MyIP
		transToFollower()
	}
}

func parseAction(data string) Action {
	tokens := strings.Split(data, " ")
	action := Action{}
	if tokens[0] == "set" {
		action.Type = Set
		action.Key = tokens[1]
		action.Value = tokens[2]
	}
	if tokens[0] == "get" {
		action.Type = Get
		action.Key = tokens[1]
	}
	if tokens[0] == "delete" {
		action.Type = Delete
		action.Key = tokens[1]
	}
	action.Term = inform.term
	action.Index = logIndex
	return action
}
