gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative proto/storage/storage.proto
build:
	docker build -t tages-service . && docker run -it --rm -p 50051:50051 tages-service