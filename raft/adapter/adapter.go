package adapter

import (
	"encoding/binary"
	"encoding/json"
	"net"

	kcp "github.com/xtaci/kcp-go"
)

type Adapter struct {
	peers    map[int]string
	listener net.Listener
}

const (
	data_from_client = iota
	request_vote
	request_vote_reply
	append_entries
	append_entries_reply
)

func InitAdapter(peers map[int]string) *Adapter {
	a := new(Adapter)
	a.peers = make(map[int]string)
	a.listener, _ = kcp.Listen("0.0.0.0:18230")
	return a
}

func (a *Adapter) ListenLoop(onMsg func(packet Packet)) {
	for {
		listenConn, _ := a.listener.Accept()
		//用来对齐消息
		magicBytes := make([]byte, 4)
		listenConn.Read(magicBytes)
		magicNum := binary.BigEndian.Uint32(magicBytes)
		if magicNum != 0x12345678 {
			continue
		}
		lengthByte := make([]byte, 4)
		listenConn.Read(lengthByte)
		length := binary.BigEndian.Uint32(lengthByte)
		data := make([]byte, length)
		listenConn.Read(data)
		packet := parseData(data, listenConn.RemoteAddr().String())
		onMsg(packet)
	}
}

func parseData(data []byte, remote string) Packet {
	var packet Packet
	json.Unmarshal(data, &packet)
	packet.SourceAddr = remote
	return packet
}
