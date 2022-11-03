package main

import (
	"context"
	search "go-grpc-example/proto/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type service struct {
	search.UnimplementedSearchServiceServer
}

func (s *service) Search(ctx context.Context, req *search.SearchRequest) (res *search.SearchResponse, err error) {
	//fmt.Println(req.GetRequest())
	return &search.SearchResponse{Response: req.GetRequest() + " Server"}, nil
}

const PORT = "8888"

func main() {
	// 根据服务端输入的证书文件和密钥构造 TLS 凭证
	c, err := credentials.NewServerTLSFromFile("./conf/server_side_TLS/server.pem", "./conf/server_side_TLS/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	// 返回一个 ServerOption，用于设置服务器连接的凭据。
	// 用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项
	s := grpc.NewServer(grpc.Creds(c))
	lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	//将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。
	//这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
	search.RegisterSearchServiceServer(s, &service{})

	//gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
	s.Serve(lis)
}
