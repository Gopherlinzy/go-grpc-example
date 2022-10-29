package main

import (
	pb "go-grpc-example/proto/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	"io"
	"log"
	"net"
	"time"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

const PORT = "8888"

func main() {
	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
	server := grpc.NewServer(grpc.Creds(altsTC)) //创建 gRPC Server 对象
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

//服务端流式RPC，Server是Stream，Client为普通RPC请求
//客户端发送一次普通的RPC请求，服务端通过流式响应多次发送数据集
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

//客户端流式RPC，单向流
//客户端通过流式多次发送RPC请求给服务端，服务端发送一次普通的RPC请求给客户端
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
		time.Sleep(time.Second)
	}
	return nil
}

//双向流，由客户端发起流式的RPC方法请求，服务端以同样的流式RPC方法响应请求
//首个请求一定是client发起，具体交互方法（谁先谁后，一次发多少，响应多少，什么时候关闭）根据程序编写方式来确定（可以结合协程）
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
	return nil
}
