package main

import (
	"fmt"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

func main() {
	go server()
	time.Sleep(time.Second)
	go client()
	time.Sleep(time.Second * 5)
}

func server() {
	listen, _ := kcp.Listen("0.0.0.0:18230")
	conn, _ := listen.Accept()
	data := make([]byte, 1024)
	conn.Read(data)
	fmt.Println(string(data))
}

func client() {
	conn, err := kcp.Dial("10.22.35.255:18230")
	if err != nil {
		fmt.Print(err)
	}
	conn.Write([]byte("Hello World!"))
}
