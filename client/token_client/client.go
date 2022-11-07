package main

import (
	"context"
	"fmt"
	tokenAuth "go-grpc-example/pkg/token"
	"go-grpc-example/proto/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const Address = "127.0.0.1:8888"

func main() {
	var opts []grpc.DialOption

	if tokenAuth.IsTLS {
		//打开tls 走tls认证
		// 根据客户端输入的证书文件和密钥构造 TLS 凭证。
		// 第二个参数 serverNameOverride 为服务名称。
		c, err := credentials.NewClientTLSFromFile("./conf/server_side_TLS/server.pem", "go-grpc-example")
		if err != nil {
			log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(c))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// token信息
	auth := tokenAuth.TokenAuth{
		token.TokenValidateParam{
			Token: "81dc9bdb52d04dc20036dbd8313ed055",
			Uid:   1234,
		},
	}
	opts = append(opts, grpc.WithPerRPCCredentials(&auth))
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		fmt.Println("grpc.Dial error:", err)
		return
	}
	defer conn.Close()
	// 实例化客户端
	client := token.NewTokenServiceClient(conn)

	// 调用具体方法
	token, err := client.Token(context.Background(), &token.Request{Name: "linzy"})
	if err != nil {
		fmt.Println("client.Token error:", err)
		return
	}
	fmt.Println("return result:", token)
}
