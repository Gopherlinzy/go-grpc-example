package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

const HelloServiceName = "server/tcp-server/server.HiLinzy"

type HelloServiceClient struct {
	*rpc.Client
}

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (h *HelloServiceClient) SayHi(request string, response *string) error {
	return h.Client.Call(HelloServiceName+".SayHi", request, &response)
}

func main() {
	//建立连接
	//dial, err := rpc.Dial("tcp", "127.0.0.1:8888")
	//client, err := DialHelloService("tcp", "127.0.0.1:8888")
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var result string
	for i := 0; i < 5; i++ {
		//发起请求
		//err = client.SayHi("linzy", &result)
		client.Call(HelloServiceName+".SayHi", "linzy", &result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("rpc service result:", result)
		time.Sleep(time.Second)
	}
}
