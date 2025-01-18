package main

import (
	"fmt"
	"net"
)

type User struct {
	addr string
	name string
	msg  chan string
}

var allUser = make(map[string]User)   //存放所有登录的用户
var message = make(chan string, 1024) //用于存放任何人发送过来的信息，作为全局管道

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("监听启动失败")
	} else {
		fmt.Println("开始监听")
	}

	go savemessage() // 启动广播协程

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接收连接失败")
		}
		go chat(conn)
	}
}

func chat(conn net.Conn) {
	var rec []byte = make([]byte, 1024)
	user := User{
		addr: conn.RemoteAddr().String(),
		name: "xioaming",
		msg:  make(chan string, 10),
	}
	allUser[user.addr] = user
	msg := fmt.Sprintf("%s:%s上线了\n", user.addr, user.name)
	message <- msg
	fmt.Println(msg)

	for {
		_, err := conn.Read(rec)
		if err != nil {
			fmt.Println("读取数据失败:", err)
			break
		}
		msg := string(rec)
		message <- fmt.Sprintf("%s:%s\n", user.name, msg)
		go broadcast(&user, conn)
	}

	delete(allUser, user.addr)
	message <- fmt.Sprintf("%s下线了\n", user.name)

	conn.Close()
}

func savemessage() {
	for {
		msg := <-message
		for _, user := range allUser {
			user.msg <- msg
		}
	}
}
func broadcast(user *User, conn net.Conn) {
	for _msg := range user.msg {
		_, err := conn.Write([]byte(_msg))
		if err != nil {
			fmt.Println("广播失败")
		}
	}
}
