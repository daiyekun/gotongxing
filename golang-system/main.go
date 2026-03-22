package main

import "fmt"

func main() {
	fmt.Println("准备启动服务器Start....")
	server := NewServer("127.0.0.1", 8888)
	server.Start()
	fmt.Println("启动服务器成功....")
}
