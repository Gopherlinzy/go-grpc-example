package Interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
	"runtime"
	"time"
)

//  UnaryClientInterceptor 一个简单的 unary interceptor 示例。
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Println("我是第一个拦截器")
		// 预处理(pre-processing)
		start := time.Now()
		// 获取正在运行程序的操作系统
		cos := runtime.GOOS
		// 将操作系统信息附加到传出请求
		ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)

		// 可以看做是当前 RPC 方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
		// 调用RPC方法(invoking RPC method)
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 后处理(post-processing)
		end := time.Now()
		log.Printf("RPC: %s,,client-OS: '%v' req:%v start time: %s, end time: %s, err: %v", method, cos, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
		return err
	}
}

func UnaryClientInterceptorTwo() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Println("我是第二个拦截器")
		// 可以看做是当前 RPC 方法，一般在拦截器中调用 invoker 能达到调用 RPC 方法的效果，当然底层也是 gRPC 在处理。
		// 调用RPC方法(invoking RPC method)
		_ = invoker(ctx, method, req, reply, cc, opts...)
		return nil
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 预处理(pre-processing)
		start := time.Now()
		// 从传入上下文获取元数据
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("couldn't parse incoming context metadata")
		}

		// 检索客户端操作系统，如果它不存在，则此值为空
		os := md.Get("client-os")
		// 获取客户端IP地址
		ip, err := getClientIP(ctx)
		if err != nil {
			return nil, err
		}

		// RPC 方法真正执行的逻辑
		// 调用RPC方法(invoking RPC method)
		m, err := handler(ctx, req)
		end := time.Now()
		// 记录请求参数 耗时 错误信息等数据
		// 后处理(post-processing)
		log.Printf("RPC: %s,client-OS: '%v' and IP: '%v' req:%v start time: %s, end time: %s, err: %v", info.FullMethod, os, ip, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
		return m, err
	}
}

// GetClientIP检查上下文以检索客户机的ip地址
func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Printf("opening client streaming to the server method: %v", method)
		// 调用Streamer函数，获得ClientStream
		stream, err := streamer(ctx, desc, cc, method)
		return newStreamClient(stream), err
	}
}

// 嵌入式 streamClient 允许我们访问SendMsg和RecvMsg函数
type streamClient struct {
	grpc.ClientStream
}

func newStreamClient(c grpc.ClientStream) grpc.ClientStream {
	return &streamClient{c}
}

// RecvMsg从流中接收消息
func (e *streamClient) RecvMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ClientStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

// RecvMsg从流中接收消息
func (e *streamClient) SendMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ClientStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := newStreamServer(ss)
		return handler(srv, wrapper)
	}
}

// 嵌入式EdgeServerStream允许我们访问RecvMsg函数
type streamServer struct {
	grpc.ServerStream
}

func newStreamServer(s grpc.ServerStream) grpc.ServerStream {
	return &streamServer{s}
}

// RecvMsg从流中接收消息
func (e *streamServer) RecvMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

// RecvMsg从流中接收消息
func (e *streamServer) SendMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	if err := e.ServerStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}
