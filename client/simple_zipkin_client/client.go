package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go-grpc-example/proto/hello"
	"google.golang.org/grpc"
	"log"
)

const (
	PORT                      = "8888"
	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

func main() {
	reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(SERVICE_NAME, ZIPKIN_RECORDER_HOST_PORT)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	// 建立链接
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		))

	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	// 一定要记得关闭链接
	defer conn.Close()

	// 实例化客户端
	client := hello.NewUserServiceClient(conn)
	// 发起请求
	response, err := client.SayHi(context.Background(), &hello.Request{Name: "lin钟一"})
	if err != nil {
		log.Fatalf("client.SayHi err: %v", err)
	}
	fmt.Printf("resp: %s", response.GetResult())
}
