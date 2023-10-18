package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func main() {
	go server()
	// go server()
	time.Sleep(time.Second)
	go client()
	time.Sleep(time.Second * 10)
}

func server() {
	conn, _ := net.ListenPacket("udp4", ":18230")
	data := make([]byte, 1024)
	conn.ReadFrom(data)
	for i := 0; i < 1024; i++ {
		// fmt.Println(data[i])
		if data[i] == 0 {
			data = data[:i]
			break
		}
	}
	// fmt.Println(string(data))
	var str string
	json.Unmarshal(data, &str)
	fmt.Println(str)
	fmt.Println(string(data))
}

func client() {
	//mask is 255.255.252.0
	//ip is 10.22.34.71
	//34=32+2
	//broadcast address is 10.22.35.255
	pc, _ := net.ListenPacket("udp4", ":18231")
	addr, _ := net.ResolveUDPAddr("udp4", "10.22.35.255:18230")
	for {
		pc.WriteTo([]byte("data to transmit"), addr)
	}
}
