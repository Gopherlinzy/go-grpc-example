package main

import (
	"context"
	search2 "go-grpc-example/proto/search"
	"google.golang.org/grpc"
	"log"
)

const PORT = "8888"

func main() {
	//创建与给定目标（服务端）的连接交互
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	//创建 SearchService 的客户端对象
	client := search2.NewSearchServiceClient(conn)
	//发送 RPC 请求，等待同步响应，得到回调后返回响应结果
	resp, err := client.Search(context.Background(), &search2.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	//输出响应结果
	log.Printf("resp: %s", resp.GetResponse())
}
