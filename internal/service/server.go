package service

import (
	"fmt"
	"log"
	"net"

	"github.com/NotYourAverageFuckingMisery/imageService/internal/store"
	v1 "github.com/NotYourAverageFuckingMisery/imageService/proto/v1"

	"google.golang.org/grpc"
)

// Server runs on two different ports to provide concurrency limits with grpc.MaxConcurrentStreams
type Server struct {
	TServer *grpc.Server
	IServer *grpc.Server
}

func NewServer(store *store.DiskImageStore, tStreams uint32, iStreams uint32) *Server {

	s := &Server{
		TServer: grpc.NewServer(
			grpc.MaxConcurrentStreams(tStreams),
		),
		IServer: grpc.NewServer(
			grpc.MaxConcurrentStreams(iStreams),
		),
	}
	v1.RegisterImageInfoServiceServer(s.IServer, &ImageInfoServer{DiskImageStore: store})
	v1.RegisterTransferImageServiceServer(s.TServer, &TransferImageServer{DiskImageStore: store})

	return s
}

// This is bad, had no time to fix
func (s *Server) Run(tsAddr string, isAddr string) {
	tsLis, err := net.Listen("tcp", tsAddr)
	if err != nil {
		log.Fatalf("Failed to listen to: %v", err)
	}
	isLis, err := net.Listen("tcp", isAddr)
	if err != nil {
		log.Fatalf("Failed to listen to: %v", err)
	}

	fmt.Println("Image transfer service started at", tsAddr)

	go func() {
		if err := s.TServer.Serve(tsLis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	fmt.Println("Image info service started at", isAddr)

	if err := s.IServer.Serve(isLis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
