package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go-grpc-example/pkg/Interceptor"
	"go-grpc-example/proto/hello"
	"google.golang.org/grpc"
	"log"
)

const PORT = "8888"

func main() {
	// 建立链接
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			//Interceptor.UnaryClientInterceptor()
			// 按照顺序依次执行截取器
			grpc_middleware.ChainUnaryClient(Interceptor.UnaryClientInterceptor(),
				Interceptor.UnaryClientInterceptorTwo()),
		))
	//conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(),
	//	grpc.WithUnaryInterceptor(Interceptor.UnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	// 一定要记得关闭链接
	defer conn.Close()

	// 实例化客户端
	client := hello.NewUserServiceClient(conn)
	// 发起请求
	response, err := client.SayHi(context.Background(), &hello.Request{Name: "lin钟一"})
	if err != nil {
		log.Fatalf("client.SayHi err: %v", err)
	}
	fmt.Printf("resp: %s", response.GetResult())
}
