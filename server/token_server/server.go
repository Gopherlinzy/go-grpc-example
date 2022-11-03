package main

import (
	"fmt"
	"go-grpc-example/proto/token"
	"google.golang.org/grpc"
	"net"
)

type TokenService struct {
	token.UnimplementedTokenServiceServer
}

func main() {

	server := grpc.NewServer()
	token.RegisterTokenServiceServer(server, &TokenService{})

	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("start error:", err)
		return
	}
	fmt.Println("服务启动成功....")
	server.Serve(listen)
}
