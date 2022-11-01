gen:
	protoc --go_out=. --go-grpc_out=. ./proto/search/*.proto
	protoc --go_out=. --go-grpc_out=. ./proto/stream/*.proto
	protoc --go_out=. --go-grpc_out=. ./proto/hello/*.proto