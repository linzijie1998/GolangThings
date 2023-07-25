package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/linzijie1998/GolangThings/backend/gRPC/gateway/pb/message"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
)

var (
	wg  = sync.WaitGroup{}
	ctx = context.Background()
)

type MessageService struct {
	message.UnimplementedMessageServiceServer
}

func (*MessageService) Chat(ctx context.Context, req *message.Req) (*message.Resp, error) {
	return &message.Resp{Content: fmt.Sprintf("I got your message: %s", req.GetContent())}, nil
}

func (*MessageService) Echo(ctx context.Context, req *message.Req) (*message.Resp, error) {
	return &message.Resp{Content: fmt.Sprintf("I got your message: %s", req.GetContent())}, nil
}

func registerGRPC() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	message.RegisterMessageServiceServer(server, &MessageService{})
	if err = server.Serve(listen); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

func registerGateway() {
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8888", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	mux := runtime.NewServeMux()
	gwServer := http.Server{
		Handler: mux,
		Addr:    ":8090",
	}
	if err = message.RegisterMessageServiceHandler(ctx, mux, conn); err != nil {
		log.Fatal(err)
	}
	if err = gwServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

func main() {
	wg.Add(2)
	go registerGRPC()
	go registerGateway()
	wg.Wait()
}
