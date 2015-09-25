package main

import (
	"fmt"
	"time"
	"usercmd"
	"common"
	"github.com/gogo/protobuf/proto"
	"github.com/cpp2go/gonet"
)

type Client struct {
	gonet.TcpTask
	mclient *gonet.TcpClient
}

func NewClient() *Client {
	s := &Client{
		TcpTask: *gonet.NewTcpTask(nil),
	}
	s.Derived = s
	return s
}

func (this *Client) Connect(addr string) bool {

	conn, err := this.mclient.Connect(addr)
	if err != nil {
		fmt.Println("连接失败 ", addr)
		return false
	}

	this.Conn = conn

	this.Start()

	fmt.Println("连接成功 ", addr)
	return true
}

func (this *Client) ParseMsg(data []byte, flag byte) bool {

	this.Verify()

	this.AsyncSend(data, flag)

	return true
}

func (this *Client) SendCmd(cmd usercmd.UserCmd, msg common.Message) bool {
	data, flag, err := common.EncodeCmd(uint16(cmd), msg)
	if err != nil {
		fmt.Println("[服务] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return false
	}
	return this.AsyncSend(data, flag)
}

func (this *Client) OnClose() {

}

func main() {

	for {

		client := NewClient()

		if !client.Connect("127.0.0.1:80") {
			return
		}

		retCmd := &usercmd.ReqUserLogin{
			Account:  proto.String("abcd"),
			Password: proto.String("123456"),
			Key:      proto.Uint32(100),
		}
		client.SendCmd(usercmd.UserCmd_Login, retCmd)

		time.Sleep(time.Millisecond * 1)
	}
}
