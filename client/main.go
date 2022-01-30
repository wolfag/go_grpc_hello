package main

import (
	"context"
	"io"
	"log"

	pb "github.com/wolfag/go_grpc_hello/hello"
	"google.golang.org/grpc"
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

	primeStream(client, 120)

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
