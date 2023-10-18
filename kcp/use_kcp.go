package main

import (
	"encoding/json"
	"fmt"
	"time"

	kcp "github.com/xtaci/kcp-go"
)

func main() {
	server()
	time.Sleep(time.Second)
	// go client()
	time.Sleep(time.Second)
}

func server() {
	listen, _ := kcp.Listen("0.0.0.0:18230")
	for {
		conn, _ := listen.Accept()
		data := make([]byte, 1024)
		conn.Read(data)
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
	}
}

func client() {
	conn, err := kcp.Dial("10.22.35.255:18230")
	if err != nil {
		fmt.Print(err)
	}
	data, _ := json.Marshal("hello,world")
	conn.Write(data)
}
