package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlieMap map[string]*User
	mapLock  sync.RWMutex
	//消息广播channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:       ip,
		Port:     port,
		OnlieMap: make(map[string]*User),
		Message:  make(chan string),
	}

	return server
}

//监听message 广播消息channel 的goroutine ,一旦有消息就发哦是那个给全部的User

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlieMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	fmt.Println("消息：" + sendMsg)
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务
	//fmt.Println("链接成功.....")
	user := NewUser(conn, this)
	user.Online()
	//接受用户发来的消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			//接受用户消息（去除'\n'）
			msg := string(buff[:n-1])
			//将消息广播
			//this.BroadCast(user, msg)
			user.DoMessage(msg)
		}
	}()

	//阻塞
	select {}
}

func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("litener accept err:", err)
			continue
		} else {
			fmt.Println("服务器启动 Ip:", this.Ip, "端口：", this.Port)
		}
		//do handler
		go this.Handler(conn)
	}

}
