package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start ...")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("connect to server fail ...")
		return
	}
	for {
		_, err = conn.Write([]byte("Hello Zinx V0.1"))
		if err != nil {
			fmt.Println("write to server error,", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("receive from server error,", err)
			return
		}
		fmt.Printf("server call back %s, cnt=%d\n", buf, cnt)
		time.Sleep(time.Second)
	}
}
