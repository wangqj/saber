package main

import (
	"os"
	"net"
	"fmt"
	"time"
)

func main() {
	netListen, err := net.Listen("tcp", "127.0.0.1:16379")
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer netListen.Close()
	for i := 0; ; i++ {
		fmt.Println(i)
		conn, _ := netListen.Accept()
		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Millisecond))

		//fmt.Println(conn.RemoteAddr().String(), " tcp connect success")
		go handle(conn)

	}
}

func handle(proxyConn net.Conn) {
	//defer proxyConn.Close()
	//fmt.Println("handle request")
	proxyBuffer := make([]byte, 2048)
	proxyConn.Read(proxyBuffer)
	//proxyBuffer = bytes.TrimRight(proxyBuffer, "\x00")
	//fmt.Println("request content is %s", string(proxyBuffer))
	proxyConn.Write([]byte("+OK\r\n"))
	//time.Sleep(3*time.Second)

	//fmt.Println("over")
}
