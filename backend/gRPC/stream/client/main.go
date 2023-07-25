package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/linzijie1998/GolangThings/backend/gRPC/stream/pb/message"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"time"
)

var (
	ctx = context.Background()
)

func serverStreamDemo(client message.MessageServiceClient) {
	serverStream, err := client.ServerStream(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err := serverStream.Send(&message.Req{Content: fmt.Sprintf("This is %dth req message!", i+1)})
		if err != nil {
			log.Fatal(err)
		}
	}
	resp, err := serverStream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetContent())
}

func clientStreamDemo(client message.MessageServiceClient) {
	clientStream, err := client.ClientStream(ctx, &message.Req{Content: "I need resp message!"})
	if err != nil {
		log.Fatal(err)
	}
	for {
		resp, err := clientStream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Receive over!")
				break
			}
			log.Fatal(err)
		}
		log.Println(resp.GetContent())
	}
}

func bidirectionalStreamDemo(client message.MessageServiceClient) {
	bidirectionalStream, err := client.BidirectionalStream(ctx)
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 10; i++ {
			err := bidirectionalStream.Send(&message.Req{Content: fmt.Sprintf("This is %d req message!", i+1)})
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
		if err = bidirectionalStream.CloseSend(); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		for {
			resp, err := bidirectionalStream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("Receive over!")
					break
				}
				log.Fatal(err)
			}
			log.Println(resp.GetContent())
		}
		wg.Done()
	}()

	wg.Wait()
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := message.NewMessageServiceClient(conn)

	//serverStreamDemo(client)
	//clientStreamDemo(client)
	bidirectionalStreamDemo(client)
}
