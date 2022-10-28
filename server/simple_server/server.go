package main

import (
	"context"
	search "go-grpc-example/proto/search"
	"google.golang.org/grpc"
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
	s := grpc.NewServer() //创建 gRPC Server 对象
	//将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。
	//这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
	search.RegisterSearchServiceServer(s, &service{})

	lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	//gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
	s.Serve(lis)
}
