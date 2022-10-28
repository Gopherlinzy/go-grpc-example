package main

import (
	pb "go-grpc-example/proto/stream"
	"google.golang.org/grpc"
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
	server := grpc.NewServer()
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

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}
