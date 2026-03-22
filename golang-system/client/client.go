package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIP:   serverIp,
		ServerPort: serverPort,
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

var serverIp string
var serverPort int

// client -ip 127.0.0.1

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器Ip（默认：127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置端口默认8888")

}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(".......服务器链接失败...")
		return
	}

	fmt.Println(">>>>服务器链接成功....")

	// 启动一个协程来处理从服务器接收到的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := client.conn.Read(buf)
			if err != nil {
				fmt.Println("连接断开:", err)
				return
			}
			fmt.Println("收到服务器消息:", string(buf[:n]))
		}
	}()

	// 启动一个协程来处理用户输入并发送给服务器
	go func() {
		for {
			var msg string
			fmt.Print("请输入消息: ")
			// 注意：在实际生产中，混合使用 fmt.Scan 和 net.Read 可能会有缓冲问题，
			// 这里仅做简单演示。更推荐使用 bufio.NewReader(os.Stdin)
			if _, err := fmt.Scanln(&msg); err != nil {
				continue
			}
			if msg == "quit" {
				client.conn.Close()
				return
			}
			_, err := client.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("发送失败:", err)
				return
			}
		}
	}()

	//启动客户端
	select {}
}
