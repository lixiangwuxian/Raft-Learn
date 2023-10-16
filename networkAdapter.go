package raft

type adapter struct {
	peers map[int]string
}

type FollowerPacket struct {
	typeOfMsg int
	term      int
	data      string
}

func (*adapter) ReadNewestData() {
}

func (*adapter) WriteDataTo(peer int) {
}

func (*adapter) VoteTo(peer int) {
}

func (*adapter) SendHeartbeatTo(peer int) {
}

func (*adapter) SendCommitTo(peer int) {
}

func (*adapter) SendLogTo(peer int) {

}
