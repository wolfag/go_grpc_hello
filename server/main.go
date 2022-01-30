package main

import (
	"log"
	"net"

	pb "github.com/wolfag/helloworld/hello"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
