package main

import (
	"github.com/linzijie1998/GolangThings/backend/gRPC/stream/pb/message"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	message.RegisterMessageServiceServer(server, &MessageService{})
	err = server.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
