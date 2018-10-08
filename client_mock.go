package main

import (
	"net"
	"fmt"
	"os"
	"time"
	"bytes"
	"saber/proxy"
)

//模拟客户端

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:16379")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	var list = []string{"set a 1", "get a", "set a 2", "get a"}
	//var list = []string{"set a 1"}

	for i := 0; i < len(list); i++ {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println("connect success")
		fmt.Println(i)
		go sender(conn, list[i])
	}
	time.Sleep(time.Duration(10000 * time.Second))
}

func sender(conn net.Conn, content string) {
	r := bytes.NewReader([]byte(content))
	cmd, _ := proxy.ReadCommand(r)
	b := cmd.Format()
	fmt.Println(string(b))

	//resp := redisy.NewString([]byte(content))
	//r, err := redisy.EncodeToBytes(resp)
	//fmt.Println("format result  ", r)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "formet error: %s", err.Error())
	//}
	conn.Write(b)
	buffer := make([]byte, 2048)
	_, eer := conn.Read(buffer);
	bytes.Trim(buffer, " ")
	fmt.Println("send over ", string(buffer), content, eer)
}