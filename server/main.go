package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/wolfag/go_grpc_hello/hello"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (*server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + req.Name,
	}, nil
}

func (*server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{
		Result: req.N1 + req.N2,
	}, nil
}

func (*server) PrimeNumber(req *pb.PrimeRequest, stream pb.Greeter_PrimeNumberServer) error {
	n := req.GetN()
	k := int32(2)
	log.Printf("go... %v\n", n)
	for n > 1 {
		if n%k == 0 {
			n = n / k

			stream.Send(&pb.PrimeResponse{Result: k})
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			log.Printf("k=%v\n", k)
		}
	}
	return nil
}

func (*server) Average(stream pb.Greeter_AverageServer) error {
	log.Printf("-----Average")
	total := float32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatal(err)
			return err
		}

		log.Printf("---recv: %v\n", req.GetN())

		total += req.GetN()
		count++
	}
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
