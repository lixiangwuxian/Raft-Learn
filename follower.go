package raft

const (
	waitForData int = iota
	waitForCommit
)

// type FollowerPacket struct {
// 	typeOfMsg int
// 	term      int
// 	data      string
// }

// var tcp_to_leader *net.TCPConn

var data_to_submit Action

var leaderTimeout int

func cacheAction(data Action) {
	data_to_submit = data
}

func commitAction() {
	if data_to_submit.Type == Set {
		kvStore.Set(data_to_submit.Key, data_to_submit.Value)
	}
	if data_to_submit.Type == Delete {
		kvStore.Delete(data_to_submit.Key)
	}
	//no need to implement get
}

// func get_data_from_leader() (string, error) {
// 	var data []byte
// 	var err error
// 	for {
// 		_, err = tcp_to_leader.Read(data)
// 		if err != nil {
// 			return "", err
// 		}
// 		packet := parse_data(string(data))
// 		if packet.typeOfMsg == log_data {
// 			return string(data), nil
// 		}
// 	}
// }

// func parse_data(data string) FollowerPacket {
// 	packet := FollowerPacket{}
// 	json.Unmarshal([]byte(data), &packet)
// 	return packet
// }
