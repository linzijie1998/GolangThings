# gRPC的流式传输

## 01 流式传输模式

根据请求和响应的是否是`流式（stream）传输`可以将服务分为以下几种传输模式：

| 传输模式                | 请求流 | 响应流 |
|---------------------|-----|-----|
| ServerStream        | ✅   |     |
| ClientStream        |     | ✅   |
| BidirectionalStream | ✅   | ✅   |

在`proto3`中可以使用`stream`关键字来定义RPC服务的传输模式：
```go
service MessageService {
    rpc ServerStream(stream Req) returns (Resp);
    rpc ClientStream(Req) returns (stream Resp);
    rpc BidirectionalStream(stream Req) returns (stream Resp);
}
```

## 02 代码实现
### 2.1 ServerStream
ServerStream模式的接口包含一个`Recv`方法用于接收流式传入的请求体，一个`SendAndClose`方法用于返回一个响应体并且关闭流：
```go
type MessageService_ServerStreamServer interface {
	SendAndClose(*Resp) error // 发送响应体并且关闭流
	Recv() (*Req, error)      // 接收请求流
	grpc.ServerStream
}
```
下面是一个ServerStream模式的实现Demo：
```go
// server
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

// client
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
```
### 2.2 ClientStream
ClientStream模式的接口仅需要一个`Send`方法用于返回响应体：
```go
type MessageService_ClientStreamServer interface {
	Send(*Resp) error // 发送响应流
	grpc.ServerStream
}
```
下面是一个ClientStream模式的实现Demo：
```go
// server
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
// client
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
```
### 2.3 BidirectionalStream
BidirectionalStream模式的接口包含一个`Send`方法发送响应体，以及一个`Recv`方法来接收请求流：
```go
type MessageService_BidirectionalStreamServer interface {
	Send(*Resp) error
	Recv() (*Req, error)
	grpc.ServerStream
}
```
下面是一个BidirectionalStream模式的实现Demo：
```go
// server
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

// client
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
```