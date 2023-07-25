package main

import (
	"errors"
	"fmt"
	"github.com/linzijie1998/GolangThings/backend/gRPC/stream/pb/message"
	"io"
	"log"
	"time"
)

type MessageService struct {
	message.UnimplementedMessageServiceServer
}

func (*MessageService) ServerStream(server message.MessageService_ServerStreamServer) error {
	for {
		req, err := server.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Received over!")
				if err := server.SendAndClose(&message.Resp{Content: "Received over!"}); err != nil {
					return err
				}
				break
			}
			return err
		}
		log.Println(req.GetContent())
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*MessageService) ClientStream(req *message.Req, server message.MessageService_ClientStreamServer) error {
	log.Println(req.GetContent())
	for i := 0; i < 10; i++ {
		err := server.Send(&message.Resp{Content: fmt.Sprintf("This is %dth resp message!", i+1)})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*MessageService) BidirectionalStream(server message.MessageService_BidirectionalStreamServer) error {
	type meta struct {
		msg string
		err error
	}
	metaChan := make(chan meta)

	go func() {
		for {
			req, err := server.Recv()
			if err != nil {
				metaChan <- meta{err: err}
				break
			}
			metaChan <- meta{msg: req.GetContent()}
		}
	}()

	for {
		m := <-metaChan
		if m.err != nil {
			if errors.Is(m.err, io.EOF) {
				log.Println("Receive over!")
				return nil
			}
			return m.err
		}
		err := server.Send(&message.Resp{Content: fmt.Sprintf("I got your message: %s", m.msg)})
		if err != nil {
			return err
		}
	}
}
