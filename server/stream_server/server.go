package main

import (
	"fmt"
	pb "go-grpc-example/proto/stream"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

const PORT = "8888"

func main() {
	server := grpc.NewServer() //创建 gRPC Server 对象
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

//服务端流式RPC，Server是Stream，Client为普通RPC请求
//客户端发送一次普通的RPC请求，服务端通过流式响应多次发送数据集
/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	// 具体返回多少个response根据业务逻辑调整
	for n := 0; n <= 6; n++ {
		// 通过 send 方法不断推送数据
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
	// 返回nil表示已经完成响应
	return nil
}

//客户端流式RPC，单向流
//客户端通过流式多次发送RPC请求给服务端，服务端发送一次响应给客户端
/*
1. for循环中通过stream.Recv()不断接收client传来的数据
2. err == io.EOF表示客户端已经发送完毕关闭连接了,此时在等待服务端处理完并返回消息
3. stream.SendAndClose() 发送消息并关闭连接(虽然在客户端流里服务器这边并不需要关闭 但是方法还是叫的这个名字，内部也只会调用Send())
*/
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	// for循环接收客户端发送的消息
	for {
		// 通过 Recv() 不断获取客户端 send()推送的消息
		r, err := stream.Recv()
		// err == io.EOF表示已经获取全部数据
		if err == io.EOF {
			// SendAndClose 返回并关闭连接
			// 在客户端发送完毕后服务端即可返回响应
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
/*
// 1. 建立连接 获取client
// 2. 通过client调用方法获取stream
// 3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
// 3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
// 3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
*/
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	var (
		wg    sync.WaitGroup //任务编排
		msgCh = make(chan *pb.StreamPoint)
	)
	wg.Add(1)
	go func() {
		n := 0
		defer wg.Done()
		for v := range msgCh {
			err := stream.Send(&pb.StreamResponse{
				Pt: &pb.StreamPoint{
					Name:  v.GetName(),
					Value: int32(n),
				},
			})
			if err != nil {
				fmt.Println("Send error :", err)
				continue
			}
			n++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			r, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv error :%v", err)
			}
			log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
			msgCh <- &pb.StreamPoint{
				Name: "gRPC Stream Server: Route",
			}
		}
		close(msgCh)
	}()

	wg.Wait() //等待任务结束

	return nil
}
