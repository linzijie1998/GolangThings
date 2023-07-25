package main

import (
	"context"
	"github.com/linzijie1998/GolangThings/backend/gRPC/helloworld/pb/hello"
	"google.golang.org/grpc"
	"log"
)

var (
	ctx = context.Background()
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := hello.NewHelloServiceClient(conn)
	resp, err := client.SayHello(ctx, &hello.Req{Message: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetMessage())
}
