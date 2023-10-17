package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	go server()
	time.Sleep(time.Second)
	go client()
	time.Sleep(time.Second * 3)
}

func server() {
	conn, _ := net.ListenPacket("udp4", ":18230")
	data := make([]byte, 1024)
	conn.ReadFrom(data)
	fmt.Println(string(data))
}

func client() {
	//mask is 255.255.252.0
	//ip is 10.22.34.71
	//34=32+2
	//broadcast address is 10.22.35.255
	pc, _ := net.ListenPacket("udp4", ":18231")
	addr, _ := net.ResolveUDPAddr("udp4", "10.22.35.255:18230")
	pc.WriteTo([]byte("data to transmit"), addr)
}
