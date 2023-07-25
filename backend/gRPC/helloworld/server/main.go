package main

import (
	"context"
	"fmt"
	"github.com/linzijie1998/GolangThings/backend/gRPC/helloworld/pb/hello"
	"google.golang.org/grpc"
	"log"
	"net"
)

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}

func (*HelloService) SayHello(ctx context.Context, req *hello.Req) (*hello.Resp, error) {
	msg := fmt.Sprintf("I got message from request: %s\n", req.GetMessage())
	log.Println(msg)
	return &hello.Resp{Message: msg}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	hello.RegisterHelloServiceServer(server, &HelloService{})
	if err = server.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
