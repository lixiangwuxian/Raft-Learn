package adapter

import (
	"encoding/json"

	kcp "github.com/xtaci/kcp-go"
)

type KcpSender struct {
}

func (k *KcpSender) Send(remoteAddr string, data []byte) {
	conn, err := kcp.Dial(remoteAddr)
	if err != nil {
		return
	}
	conn.Write(data)

}

func (k *KcpSender) AppendEntries(peerAddr string, data AppendEntries) {
	data_byte, _ := json.Marshal(data)
	packet := Packet{
		TypeOfMsg: append_entries,
		Term:      data.Term,
		Data:      data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) AppendEntriesReply(peerAddr string, data AppendEntriesReply) {
	data_byte, _ := json.Marshal(data)
	packet := Packet{
		TypeOfMsg: append_entries_reply,
		Term:      data.Term,
		Data:      data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) RequestVote(peerAddr string, data RequestVote, currentTerm int) {
	data_byte, _ := json.Marshal(data)
	packet := Packet{
		TypeOfMsg: request_vote,
		Term:      currentTerm,
		Data:      data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) RequestVoteReply(peerAddr string, data RequestVoteReply, currentTerm int) {
	data_byte, _ := json.Marshal(data)
	packet := Packet{
		TypeOfMsg: request_vote_reply,
		Term:      currentTerm,
		Data:      data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}
