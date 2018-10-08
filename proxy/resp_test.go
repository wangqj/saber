package proxy

import (
	"testing"
	"rocket/redisy"
	"fmt"
	"bytes"
)

func TestCommand_Format(t *testing.T) {
	//ex: parse command
	body := []byte("setex name 10 walu\r\n") //其实inline command对最后的\r\n没有要求
	r := bytes.NewReader(body)

	cmd, _ := redisy.ReadCommand(r)
	fmt.Println(string(cmd.Format()))
	//so:
	//cmd.Name() == "setex"
	//cmd.Value(1) == "name"
	//cmd.Integer(2) == 10 (int64)
	//cmd.Value(3)  == "walu"

	//----------------------------
	//encode command
	cmd, _ = redisy.NewCommand("setex", "name", "10", "walu")
	body = cmd.Format()
	fmt.Println("dd=", body)

	//ex: parse command
	body2 := []byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n")
	r2 := bytes.NewReader(body2)

	data, _ := redisy.ReadCommand(r2)

	fmt.Println(data.Args[0], data.Args[1], data.Args[2])
}
