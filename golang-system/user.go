package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// 用户上线业务
func (this *User) Online() {
	//用户上线了,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlieMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播用户上线了
	this.server.BroadCast(this, "已上线了\n")
}

// 用户下线业务
func (this *User) Offline() {
	//用户上线了,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	delete(this.server.OnlieMap, this.Name)
	this.server.mapLock.Unlock()

	//广播用户下线了
	this.server.BroadCast(this, "下线了\n")
}

// 处理用户消息
func (this *User) DoMessage(msg string) {
	fmt.Println("输入消息...>>>：", msg)
	if strings.Contains(msg, "who") {
		//if  msg == "who" msg {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlieMap {
			sendmsg := "[" + user.Addr + "]" + user.Name + ":" + "online \r\n"
			this.SendMsg(sendmsg)

		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && strings.Contains(msg, "rename|") {
		newName := strings.Split(msg, "|")[1]
		//判断是否存在
		_, ok := this.server.OnlieMap[newName]
		if ok {
			this.SendMsg("curr Name Exist")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlieMap, this.Name)
			this.server.OnlieMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("你的Name 已经修改:" + this.Name + "\n")
		}
	} else if len(msg) > 4 && strings.Contains(msg, "to|") {
		//消息格式： to|小王|消息内容
		//1 获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用\"to|张三|你好\"格式。\n")
			return
		}

		//2 获取用户名称，得到对方User对象
		remoteUser, ok := this.server.OnlieMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		//3.获取消息内容，通对方的User对象发送消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("五消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name + "私信：" + content)

	} else {
		this.server.BroadCast(this, msg)
	}

}

// 监听当前用户User channerl 消息的gorountine
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

// 当前对应客户发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
