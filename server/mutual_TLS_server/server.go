package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	search "go-grpc-example/proto/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
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
	// 公钥中读取和解析公钥/私钥对
	cert, err := tls.LoadX509KeyPair("./conf/server/server.crt", "./conf/server/server.key")
	if err != nil {
		fmt.Println("LoadX509KeyPair error", err)
		return
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		fmt.Println("read ca pem error ", err)
		return
	}
	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("AppendCertsFromPEM error ")
		return
	}

	c := credentials.NewTLS(&tls.Config{
		//设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		//要求必须校验客户端的证书
		ClientAuth: tls.RequireAndVerifyClientCert,
		//设置根证书的集合，校验方式使用ClientAuth设定的模式
		ClientCAs: certPool,
	})
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
