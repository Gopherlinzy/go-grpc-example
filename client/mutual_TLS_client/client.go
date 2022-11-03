package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "go-grpc-example/proto/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

const PORT = "8888"

func main() {
	// 公钥中读取和解析公钥/私钥对
	cert, err := tls.LoadX509KeyPair("./conf/client/client.crt", "./conf/client/client.key")
	if err != nil {
		fmt.Println("LoadX509KeyPair error ", err)
		return
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./conf/ca.crt")
	if err != nil {
		fmt.Println("ReadFile ca.crt error ", err)
		return
	}
	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("certPool.AppendCertsFromPEM error ")
		return
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "go-grpc-example",
		//ServerName: "linzy",
		RootCAs: certPool,
	})

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
