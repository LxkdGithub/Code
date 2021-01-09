package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	server, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return
	}
	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Accept failed: %v\n", err)
			continue
		}
		go process(client)
	}

}


// 普通的TcpServer
//func process(client net.Conn) {
//	remoteAddr := client.RemoteAddr().String()
//	fmt.Printf("Connection from %s\n", remoteAddr)
//	client.Write([]byte("Hello world!"))
//	// []byte() == char * a = "dsv";
//	// []byte{} == char a[] = {'1', '3'};
//	client.Close()
//}


// Sccks5版本

func process(client net.Conn) {
	if err := Socks5Auth(client); err != nil {
		fmt.Println("auth error: ", err)
		client.Close()
		return
	}

	target, err := Socks5Connect(client)
	if err != nil {
		fmt.Println("connetc error:", err)
		client.Close()
		return
	}

	Socks5Forward(client, target)
}

func Socks5Auth(client net.Conn) error {
	buf := make([]byte, 256)

	// 读取VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	// 读取METHODS  这里nMETHODS就是大小
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	// 认证各种方法并挑选一种
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 {
		return errors.New("write rsp err: " + err.Error())
	}

	return nil
}

func Socks5Connect(client net.Conn) (net.Conn, error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		return nil, errors.New("invalid ver/cmd/type")
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		return nil, errors.New("invalid ver/cmd")
	}
	// cmd的2代表BINd和3代表UDP Associate
	// ......

	// 下面是地址
	// 0x01：4个字节，对应 IPv4 地址
	//
	//0x02：先来一个字节 n 表示域名长度，然后跟着 n 个字节。注意这里不是 NUL 结尾的。
	//
	//0x03：16个字节，对应 IPv6 地址

	addr := ""
	switch atyp {
	case 1: // ipv4
		n, err := io.ReadFull(client, buf[:4])
		if n != 4 {
			return nil, errors.New("invalid IPv4 address: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3: // 域名地址 n+n字节
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err := io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			return nil, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])

	case 4:  // IPv6 先不管
		return nil, errors.New("IPv6: no supported yet")

	default:
		return nil, errors.New("invalid atyp")
	}

	// 开始端口
	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		return nil, errors.New("read port: " + err.Error())
	}
	port := binary.BigEndian.Uint16(buf[:2])

	//创建连接
	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		return nil, errors.New("dial error: " + err.Error())
	}
	// 得到的LocalAddr也没什么用 就不返回了
	// 返回格式是 VER|REP(o成功错误)|RSV|ATYP|BND.ADDR|BND.PORT
	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		dest.Close()
		return nil, errors.New("write rsp: " + err.Error())
	}

	return dest, nil
}

func Socks5Forward(client, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}












