package common

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

const (
	MaxCompressSize = 1024
	MaxCmdSize      = 2
)

type Message interface {
	Marshal() (data []byte, err error)
	MarshalTo(data []byte) (n int, err error)
	Size() (n int)
	Unmarshal(data []byte) error
}

func zlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, err := w.Write(src)
	if err != nil {
		return nil
	}
	w.Close()
	return in.Bytes()
}

func zlibUnCompress(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil
	}
	return out.Bytes()
}

// 生成二进制数据,返回数据和是否压缩标识
func EncodeCmd(cmd uint16, msg Message) ([]byte, byte, error) {
	msglen := msg.Size()
	if msglen >= MaxCompressSize {
		data, err := msg.Marshal()
		if err != nil {
			fmt.Println("[协议] 编码错误 ", err)
			return nil, 0, err
		}
		mbuff := zlibCompress(data)
		p := make([]byte, MaxCmdSize+len(mbuff))
		p[0] = byte(cmd)
		p[1] = byte(cmd >> 8)
		copy(p[MaxCmdSize:], mbuff)
		return p, 1, nil
	}
	p := make([]byte, MaxCmdSize+msglen)
	_, err := msg.MarshalTo(p[MaxCmdSize:])
	if err != nil {
		fmt.Println("[协议] 编码错误 ", err)
		return nil, 0, err
	}
	p[0] = byte(cmd)
	p[1] = byte(cmd >> 8)
	return p, 0, nil
}

// 获取指令号
func GetCmd(buf []byte) uint16 {
	if len(buf) < MaxCmdSize {
		return 0
	}
	return uint16(buf[0]) | uint16(buf[1])<<8
}

// 生成protobuf数据
func DecodeCmd(buf []byte, flag byte, pb Message) Message {
	if len(buf) < MaxCmdSize {
		fmt.Println("[协议] 数据错误 ", buf)
		return nil
	}
	var mbuff []byte
	if flag == 1 {
		mbuff = zlibUnCompress(buf[MaxCmdSize:])
	} else {
		mbuff = buf[MaxCmdSize:]
	}
	err := pb.Unmarshal(mbuff)
	if err != nil {
		fmt.Println("[协议] 解码错误 ", err, ",", mbuff)
		return nil
	}
	return pb
}
