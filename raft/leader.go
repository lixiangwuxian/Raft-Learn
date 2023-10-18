package main

import (
	"encoding/json"
	"strings"
	"time"
)

// var leader_mutex sync.Locker = &sync.Mutex{}
// var good_client_count int = 0

var data_to_submit []*Action

func leaderLoop() {
	for inform.state == Leader {
		// adapter.SendHeartbeat()
		for _, data := range data_to_submit {
			entries, _ := json.Marshal(data)
			adapter.SendAppendEntries(entries)
		}
		time.Sleep(time.Duration(inform.leaderTimeout/2) * time.Millisecond)
	}
}

func dataFromClient(data string) {
	logIndex++
	action := parseAction(data)
	data_to_submit = append(data_to_submit, &action)
}

func handleAppendEntriesReply() {

}

// func addCacheReply() {
// 	leader_mutex.Lock()
// 	good_client_count++
// 	if good_client_count >= inform.totalNodes/2+1 {
// 		commitAction()
// 	}
// 	leader_mutex.Unlock()
// }

// func needKeepLeader(data HeartBeat) {
// 	if data.MyTerm > persist_inform.CurrentTerm {
// 		persist_inform.CurrentTerm = data.MyTerm
// 		inform.knownLeader = data.MyIP
// 		transToFollower()
// 	}
// }

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
	action.Term = persist_inform.CurrentTerm
	action.Index = logIndex
	return action
}
