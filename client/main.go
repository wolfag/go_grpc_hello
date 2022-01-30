package main

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "papa"})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(r)
}
