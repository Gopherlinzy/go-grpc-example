package main

import (
	"context"
	pb "go-grpc-example/proto/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const PORT = "8888"

func main() {
	// 根据客户端输入的证书文件和密钥构造 TLS 凭证。
	// 第二个参数 serverNameOverride 为服务名称。
	c, err := credentials.NewClientTLSFromFile("./conf/server_side_TLS/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}
	// 返回一个配置连接的 DialOption 选项。
	// 用于 grpc.Dial(target string, opts ...DialOption) 设置连接选项
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
