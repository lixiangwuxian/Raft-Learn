package adapter

import (
	"encoding/binary"
	"encoding/json"

	kcp "github.com/xtaci/kcp-go"
	"lxtend.com/m/constants"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
)

type KcpSender struct {
	MyAddr string
}

var magicBytes = []byte{0x12, 0x34, 0x56, 0x78}

func (k *KcpSender) Send(remoteAddr string, data []byte) {
	conn, err := kcp.Dial(remoteAddr)
	if err != nil {
		return
	}
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))
	conn.Write(magicBytes)
	conn.Write(dataLen)
	conn.Write(data)
}

func (k *KcpSender) AppendEntries(peerAddr string, data packages.AppendEntries) {
	logger.Glogger.Info("send appendEntries to %s", peerAddr)
	data_byte, _ := json.Marshal(data)
	packet := packages.Packet{
		TypeOfMsg:  constants.AppendEntries,
		SourceAddr: k.MyAddr,
		Term:       data.Term,
		Data:       data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) AppendEntriesReply(peerAddr string, data packages.AppendEntriesReply) {
	logger.Glogger.Info("send appendEntriesReply to %s", peerAddr)
	data_byte, _ := json.Marshal(data)
	packet := packages.Packet{
		TypeOfMsg:  constants.AppendEntriesReply,
		SourceAddr: k.MyAddr,
		Term:       data.Term,
		Data:       data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) RequestVote(peerAddr string, data packages.RequestVote, currentTerm int) {
	logger.Glogger.Info("send requestVote to %s", peerAddr)
	data_byte, _ := json.Marshal(data)
	packet := packages.Packet{
		TypeOfMsg:  constants.RequestVote,
		SourceAddr: k.MyAddr,
		Term:       currentTerm,
		Data:       data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}

func (k *KcpSender) RequestVoteReply(peerAddr string, data packages.RequestVoteReply, currentTerm int) {
	logger.Glogger.Info("send requestVoteReply to %s", peerAddr)
	data_byte, _ := json.Marshal(data)
	packet := packages.Packet{
		TypeOfMsg:  constants.RequestVoteReply,
		SourceAddr: k.MyAddr,
		Term:       currentTerm,
		Data:       data_byte,
	}
	packet_byte, _ := json.Marshal(packet)
	k.Send(peerAddr, packet_byte)
}
