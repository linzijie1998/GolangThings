# gRPC入门

## 01 环境配置

### 1.1 安装Go语言protocol编译插件
插件将会放到GOPATH的bin目录下：
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
如果无法使用`protoc`命令，则需要到GitHub的[releases](https://github.com/protocolbuffers/protobuf/releases/)页面下载插件，并且添加到PATH中
### 1.2 添加gRPC的Go依赖包
在go mod模式下添加grpc的依赖包：
```shell
go get -u google.golang.org/grpc
```

## 02 编写protobuf文件
```go
// 指定protobuf版本为v3
syntax = "proto3";

// option go_package = "<包路径(从mod开始写)>;<包名>";
option go_package = "github.com/linzijie1998/GolangThings/backend/gRPC/helloworld/pb/hello;hello";

package hello;

// 定义一个请求体
message Req {
  string message = 1;
}

// 定义一个响应体
message Resp{
  string message = 1;
}

// 服务
service HelloService {
  rpc SayHello(Req) returns (Resp);
}
```
## 03 代码生成
生成Go语言代码，在pb目录下执行以下命令：
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./hello/hello.proto
```
## 04 调用
### 4.1 服务端
```go
// 1. 取出Server
type HelloService struct {
    hello.UnimplementedHelloServiceServer
}

// 2. 实现接口
func (*HelloService) SayHello(ctx context.Context, req *hello.Req) (*hello.Resp, error) {
    msg := fmt.Sprintf("I got message from request: %s\n", req.GetMessage())
    log.Println(msg)
    return &hello.Resp{Message: msg}, nil
}

func main() {
	// 3. 注册服务
    listen, err := net.Listen("tcp", ":8888")
    if err != nil {
        log.Fatal(err)
    }
    server := grpc.NewServer()
    hello.RegisterHelloServiceServer(server, &HelloService{})
	// 4. 监听服务
    if err = server.Serve(listen); err != nil {
        log.Fatal(err)
    }
}
```
### 4.2 客户端
```go
var (
    ctx = context.Background()
)

func main() {
	// 1. 创建连接
    conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
	// 2. 创建Client
    client := hello.NewHelloServiceClient(conn)
	// 3. 调用方法
    resp, err := client.SayHello(ctx, &hello.Req{Message: "hello"})
    if err != nil {
        log.Fatal(err)
    }
    log.Println(resp.GetMessage())
}
```
## 参考
- https://grpc.io/docs/languages/go/quickstart/
- https://protobuf.dev/programming-guides/proto3/
- https://github.com/grpc/grpc-go/
