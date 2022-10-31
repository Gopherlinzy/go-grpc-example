package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

const HelloServiceName = "server/tcp-server/server.HiLinzy"

type HelloServiceInterface interface {
	SayHi(request string, response *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

type HelloService struct{}

func (h *HelloService) SayHi(request string, response *string) error {
	format := time.Now().Format("2006-01-02 15:04:05")
	*response = "hi " + request + "---" + format
	return nil
}

func main() {
	//注册服务
	//_ = rpc.RegisterName("HiLinzy", new(HelloService))
	RegisterHelloService(new(HelloService))
	//监听接口
	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		//监听请求
		accept, err := lis.Accept()
		if err != nil {
			log.Fatalf("Accept Error: %s", err)
		}
		//go rpc.ServeConn(accept)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(accept))
	}
}
