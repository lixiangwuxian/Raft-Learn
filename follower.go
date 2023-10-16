package raft

import (
	"encoding/json"
	"net"
	"strings"
)

const (
	waitForData int = iota
	waitForCommit
)

const (
	heartbeat int = iota
	commit
	log_data
)

// type FollowerPacket struct {
// 	typeOfMsg int
// 	term      int
// 	data      string
// }

var store *KVStore

var tcp_to_leader *net.TCPConn

var data_to_submit string

var leaderTimeout int

func follower_loop() {
	state := waitForData
	for {
		data, err := get_data_from_leader() // block, return err if time out
		if err != nil {
			return
		}
		packet := parse_data(data)
		switch packet.typeOfMsg {
		case heartbeat:
			continue
		case commit:
			if state == waitForCommit {
				commitAction(data_to_submit)
				state = waitForData
			}
			continue
		case log_data:
			if state == waitForData {
				data_to_submit = packet.data
				state = waitForCommit
			}
		}
	}
}

func commitAction(data string) {
	//get set delete
	tokens := strings.Split(data, " ")
	if tokens[0] == "get" {
		store.Get(tokens[1])
	} else if tokens[0] == "set" {
		store.Set(tokens[1], tokens[2])
	} else if tokens[0] == "delete" {
		store.Delete(tokens[1])
	}
}

func get_data_from_leader() (string, error) {
	var data []byte
	var err error
	// tcp_to_leader.SetReadDeadline()
	for {
		_, err = tcp_to_leader.Read(data)
		if err != nil {
			return "", err
		}
		packet := parse_data(string(data))
		if packet.typeOfMsg == log_data {
			return string(data), nil
		}
	}
}

func parse_data(data string) FollowerPacket {
	packet := FollowerPacket{}
	json.Unmarshal([]byte(data), &packet)
	return packet
}
