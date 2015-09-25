package main

import (
	"common"
	"fmt"
	"github.com/cpp2go/gonet"
	"net"
	"usercmd"
)

type EchoTask struct {
	gonet.TcpTask
}

func NewEchoTask(conn net.Conn) *EchoTask {
	s := &EchoTask{
		TcpTask: *gonet.NewTcpTask(conn),
	}
	s.Derived = s
	return s
}

func (this *EchoTask) ParseMsg(data []byte, flag byte) bool {

	cmd := usercmd.UserCmd(common.GetCmd(data))

	switch cmd {
	case usercmd.UserCmd_Login:
		{
			revCmd, ok := common.DecodeCmd(data, flag, &usercmd.ReqUserLogin{}).(*usercmd.ReqUserLogin)
			if !ok {
				return false
			}

			fmt.Println("> ", cmd, ",", *revCmd.Account, ",", *revCmd.Password, ",", *revCmd.Key)

			this.Verify()

			retCmd := &usercmd.ReqUserLogin{
				Account:  revCmd.Account,
				Password: revCmd.Password,
				Key:      revCmd.Key,
			}
			this.SendCmd(usercmd.UserCmd_Login, retCmd)
		}
	}

	return true
}

func (this *EchoTask) SendCmd(cmd usercmd.UserCmd, msg common.Message) bool {
	data, flag, err := common.EncodeCmd(uint16(cmd), msg)
	if err != nil {
		fmt.Println("[服务] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return false
	}
	return this.AsyncSend(data, flag)
}

func (this *EchoTask) OnClose() {

}

type EchoServer struct {
	gonet.Service
	tcpser *gonet.TcpServer
}

var serverm *EchoServer

func EchoServer_GetMe() *EchoServer {
	if serverm == nil {
		serverm = &EchoServer{
			tcpser: &gonet.TcpServer{},
		}
		serverm.Derived = serverm
	}
	return serverm
}

func (this *EchoServer) Init() bool {
	err := this.tcpser.Bind(":80")
	if err != nil {
		fmt.Println("绑定端口失败")
		return false
	}
	return true
}

func (this *EchoServer) Reload() {

}

func (this *EchoServer) MainLoop() {
	conn, err := this.tcpser.Accept()
	if err != nil {
		return
	}
	NewEchoTask(conn).Start()
}

func (this *EchoServer) Final() bool {
	this.tcpser.Close()
	return true
}

func main() {

	EchoServer_GetMe().Main()

}
