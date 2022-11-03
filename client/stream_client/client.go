package main

import (
	"context"
	"fmt"
	pb "go-grpc-example/proto/stream"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"time"
)

const PORT = "8888"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 1234}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 9999}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 1111}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}

/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	// 调用获取stream
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	// for循环获取服务端推送的消息
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// err==io.EOF则表示服务端关闭stream了
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	return nil
}

/*
1. 建立连接并获取client
2. 获取 stream 并通过 Send 方法不断推送数据到服务端
3. 发送完成后通过stream.CloseAndRecv() 关闭steam并接收服务端返回结果
*/
func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	// 获取 stream
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i <= 6; i++ {
		// 通过 Send 方法不断推送数据到服务端
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	// 发送完成后通过stream.CloseAndRecv() 关闭steam并接收服务端返回结果
	// (服务端则根据err==io.EOF来判断client是否关闭stream)
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	return nil
}

/*
1. 建立连接 获取client
2. 通过client获取stream
3. 开两个goroutine 分别用于Recv()和Send()
	3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
	3.2 Send()则由自己控制
4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
*/
func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	var wg sync.WaitGroup
	// 调用方法获取stream
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}

	// 开两个goroutine 分别用于Recv()和Send()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for n := 0; n <= 6; n++ {
			err := stream.Send(r)
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
			time.Sleep(time.Second)
		}

		// 发送完毕关闭stream
		err = stream.CloseSend()
		if err != nil {
			log.Printf("Send error:%v\n", err)
			return
		}
	}()

	wg.Wait()
	return nil
}
