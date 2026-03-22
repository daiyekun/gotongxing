package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	//监听当前用户是否活跃的channel
	isLive := make(chan bool)

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

			//用户的任意输入，代表当前用户活跃
			isLive <- true
		}
	}()

	//阻塞
	select {
	case <-isLive:
		//当前用户是活跃的，应该重置定时器
		//不做任何事情，为了激活select ,更新下面定时器
	case <-time.After(time.Minute * 10):
		//已经超时
		//当前的用户User强制关闭
		user.SendMsg("你被踢了")

		//销毁链接资源
		close(user.C)

		//关闭链接
		conn.Close()

		//推出当前handler
		return

	}
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
