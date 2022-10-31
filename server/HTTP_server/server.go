package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

const HelloServiceName = "server/tcp-server/server.HiLinzy"

type HelloService struct{}

func (h *HelloService) SayHi(request string, response *string) error {
	format := time.Now().Format("2006-01-02 15:04:05")
	*response = "hi " + request + "---" + format
	return nil
}

func main() {
	//建立连接
	rpc.RegisterName(HelloServiceName, new(HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}

		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":8888", nil)
}
