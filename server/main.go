package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"
	"time"

	pb "github.com/wolfag/go_grpc_hello/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) Max(stream pb.Greeter_MaxServer) error {
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("end Max")
			return nil
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
		num := req.GetN()
		if num > max {
			max = num
		}
		err = stream.Send(&pb.MaxResponse{Result: max})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (*server) Square(ctx context.Context, req *pb.SquareRequest) (*pb.SquareResponse, error) {
	num := req.GetN()
	log.Printf("num: %v", num)
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "expect num > 0, but got %v", num)
	}

	return &pb.SquareResponse{
		Result: math.Sqrt(float64(num)),
	}, nil
}

func (*server) SumWithDeadline(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	// simulator expensive task
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("client cancel")
			return nil, status.Error(codes.DeadlineExceeded, "client cancel request")
		}
		time.Sleep(1 * (time.Second))
	}
	return &pb.SumResponse{
		Result: req.N1 + req.N2,
	}, nil
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
