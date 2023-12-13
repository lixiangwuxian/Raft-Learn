package adapter

import (
	"encoding/binary"
	"encoding/json"
	"net"

	kcp "github.com/xtaci/kcp-go"
	"lxtend.com/m/logger"
	"lxtend.com/m/packages"
)

type Adapter struct {
	peers    []string
	listener net.Listener
}

func InitAdapter(peers []string, port string) *Adapter {
	a := new(Adapter)
	a.peers = peers
	a.listener, _ = kcp.Listen("0.0.0.0:" + port)
	return a
}

func (a *Adapter) ListenLoop(onMsg func(packet packages.Packet)) {
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
		packet := parseData(data)
		logger.Glogger.Info("receive msg from %s", packet.SourceAddr)
		logger.Glogger.Info("msg term is %d, msg type is %d", packet.Term, packet.TypeOfMsg)
		onMsg(packet)
	}
}

func parseData(data []byte) packages.Packet {
	var packet packages.Packet
	json.Unmarshal(data, &packet)
	return packet
}
