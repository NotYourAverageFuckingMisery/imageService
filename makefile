gen_proto:
	protoc -I ./proto --go_out ./ --go-grpc_out ./ ./proto/images.proto

local:
	go run cmd/main.go 
	
client:
	go run client/main.go