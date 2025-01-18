package main

import (
	"fmt"
	"net"
)

// 服务器需要给每个用户提供两个协程，一个用来读用户的信息，一个用来将msg发回给该用户，整个服务器还需要一个广播协程，用于向各个用户广播消息，所以需要存储
// 各个客户端的地址。
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("监听启动失败")
	} else {
		fmt.Println("开始监听")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接收连接失败")
		} else {
			fmt.Println("接收到了一个连接")
		}
		go chat(conn)
	}
}
func chat(conn net.Conn) {
	var rec []byte = make([]byte, 1024) //这里必须要分配空间
	_, err := conn.Read(rec)
	if err != nil {
		fmt.Println("读取数据失败:", err)
	} else {
		fmt.Println("读取到了：", string(rec))
	}
}
