gen_proto:
	protoc -I ./proto --go_out ./ --go-grpc_out ./ ./proto/images.proto

local:
	go run cmd/main.go 
	
client:
	go run client/main.go

prof_cmd:
	go build ./cmd/main.go
	./main -cpuprofile=main.prof

pprof_cmd:
	go tool pprof main main.prof