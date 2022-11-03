package main

import (
	"context"
	"fmt"
	"go-grpc-example/proto/token"
	"google.golang.org/grpc"
)

func main() {
	// token信息
	auth := token.TokenValidateParam{
		Token: "fa246d0262c3925617b0c72bb20eeb1d",
		Uid:   9999,
	}

	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
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
