package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/wolfag/go_grpc_hello/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:50069", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// sum(client)

	// hello(client, "tai")

	// primeStream(client, 120)

	// average(client)
	// max(client)
	// square(client, -98)

	sumWithDeadline(client, 4*time.Second) // ok
	sumWithDeadline(client, 1*time.Second) // timeout

}

func hello(client pb.GreeterClient, name string) {
	r, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r)
}

func sum(client pb.GreeterClient) {
	r, err := client.Sum(context.Background(), &pb.SumRequest{N1: 4, N2: 9})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r)
}
func sumWithDeadline(client pb.GreeterClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := client.SumWithDeadline(ctx, &pb.SumRequest{N1: 4, N2: 9})

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("DeadlineExceeded ")
			} else {
				log.Print("-----unknown")
			}
		} else {
			log.Fatal(err)
		}
	}
	log.Println(r)
}

func primeStream(client pb.GreeterClient, number int32) {
	stream, err := client.PrimeNumber(context.Background(), &pb.PrimeRequest{N: number})

	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			log.Printf("end stream %v \n", err)
			return
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(r.GetResult())
	}

}

func average(client pb.GreeterClient) {
	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	listReq := []pb.AverageRequest{
		{
			N: 3,
		},
		{
			N: 5,
		},
		{
			N: 6,
		},
	}

	for _, v := range listReq {
		err := stream.Send(&v)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.GetResult())
}

func max(client pb.GreeterClient) {
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan struct{})

	go func() {
		listReq := []pb.MaxRequest{
			{N: 4},
			{N: 7},
			{N: 1},
			{N: 9},
		}
		for _, v := range listReq {
			err := stream.Send(&v)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("---end stream")
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			log.Printf("rev max: %v\n", res.GetResult())
		}
		close(wait)
	}()

	<-wait
}

func square(client pb.GreeterClient, num int32) {
	res, err := client.Square(context.Background(), &pb.SquareRequest{N: num})
	if err != nil {
		log.Printf("call square root api err %v\n", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err msg: %v\n", errStatus.Message())
			log.Printf("err code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument num %v", num)
				return
			}
		}
	}

	log.Printf("square root response %v\n", res.GetResult())

}
