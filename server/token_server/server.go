package main

import (
	"fmt"
	tokenAuth "go-grpc-example/pkg/token"
	"go-grpc-example/proto/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const Address = "127.0.0.1:8888"

type TokenService struct {
	token.UnimplementedTokenServiceServer
	tokenAuth.TokenAuth
}

func (u TokenService) Token(ctx context.Context, r *token.Request) (*token.Response, error) {
	// 验证token
	_, err := u.CheckToken(ctx)
	if err != nil {
		return nil, err
	}
	return &token.Response{Name: r.GetName()}, nil
}
func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("start error:", err)
		return
	}

	var opts []grpc.ServerOption

	if tokenAuth.IsTLS {
		// TLS认证
		// 根据服务端输入的证书文件和密钥构造 TLS 凭证
		c, err := credentials.NewServerTLSFromFile("./conf/server_side_TLS/server.pem", "./conf/server_side_TLS/server.key")
		if err != nil {
			log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
		}
		opts = append(opts, grpc.Creds(c))
	}
	server := grpc.NewServer(opts...)
	token.RegisterTokenServiceServer(server, &TokenService{})

	fmt.Println("服务启动成功....")
	server.Serve(listen)
}
