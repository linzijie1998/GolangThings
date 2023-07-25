# gRPC的Gateway实现

## 01 准备工作
### 1.1 安装Go插件
```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
### 1.2 添加proto3文件
在`pb`文件夹下创建`google/api`文件夹，并且将[http.proto](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)和[annotations.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto)文件放到该目录下。

## 02 代码实现
### 2.1 修改proto3文件
```go
// 导入googleapi
import "google/api/annotations.proto";

// 添加option
service MessageService {
  rpc Chat(Req) returns (Resp) {
    option(google.api.http)={
      post: "/api/message/chat",
      body: "*",
    };
  };

  rpc Echo(Req) returns (Resp) {
    option(google.api.http)={
      get: "/api/message/chat/{content}",
    };
  };
}
```
### 2.2 生成代码
添加`--grpc-gateway_out`和`--grpc-gateway_opt`两个参数：
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out . --grpc-gateway_opt=paths=source_relative ./message/message.proto
```
### 2.3 启动服务
我们创建两个goroutine，分别用于注册gRPC服务和注册gateway：
```go
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
```

## 参考
https://grpc-ecosystem.github.io/grpc-gateway/