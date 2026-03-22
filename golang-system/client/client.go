package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前菜单模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
		flag:       9999,
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(serverIp, fmt.Sprintf("%d", serverPort)))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}

	client.conn = conn

	//返回对象

	return client
}

// 菜单提示
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊天模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("..>>>>请输入合法范围内值<<<<<<")
		return false
	}
}

var serverIp string
var serverPort int

// client -ip 127.0.0.1

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器Ip（默认：127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置端口默认8888")

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		//根据不同模式处理不通业务
		switch client.flag {
		case 1:
			//fmt.Println("公聊模式")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式")
			break
		case 3:
			//fmt.Println("更新用户名称")
			client.UpdateName()
			break
		}
	}
}

// 更新用户名称
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write er:", err)
		return false
	}
	return true
}

// 公聊天
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>输入聊天内容，exit 退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn) // 这个就相当于下面这个for
	// for {
	// 	buf := make([]byte, 1024)
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(".......服务器链接失败...")
		return
	}

	fmt.Println(">>>>服务器链接成功....")
	go client.DealResponse()
	//启动客户端
	//select {}

	client.Run()
}
