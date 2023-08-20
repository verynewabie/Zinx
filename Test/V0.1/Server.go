package main

import "Zinx/znet"

func main() {
	s := znet.NewServer("[zinx V0.1]")

	s.Serve()
}
