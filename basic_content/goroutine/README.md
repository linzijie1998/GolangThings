# Go语言并发编程
并发：多线程程序通过时间片的切换在一个核的CPU上运行

并行：多线程程序直接利用多核CPU实现同时运行
## 1 Goroutine
协程：用户态，轻量级线程，栈KB级别，可以由Go语言本身完成，比线程轻量

线程：内核态，线程跑多个协程，栈MB级别，切换、创建、停止都属于比较昂贵的系统操作

Go语言中使用`go`关键字创建协程：
```go
func hello(i int) {
    fmt.Println("hello, goroutine: " + fmt.Sprint(i))
}

func HelloGoRoutine() {
    for i := 0; i < 5; i++ {
        // 使用`go`关键字创建协程
        go func(j int) {
            hello(j)
        }(i)
    }
    // 使用Sleep阻塞，保证子协程完成前主线程不退出
    time.Sleep(time.Second)
}
```

## 2 CSP (Communicating Sequential Processes)
通过通信共享内存：需要一个`channel (通道)`，它遵循先入先出的原则，以保证收发数据的顺序；

通过共享内存通信：通过互斥量加锁获取临界区的权限，存在竞态条件和数据竞争的问题，影响程序的性能；

### 通过通信共享内存（Channel）
Go语言中使用`make`方法创建channel，`make(chan <元素类型>, <缓冲大小>)`
```go
make(chan int)       // 无缓冲通道
make(chan int, 2)    // 能存放2个元素的有缓冲通道
```
无缓冲通道也被称为**同步通道**

使用`channel`实现生产-消费模式的示例：
```go
func CalSquare() {
    src := make(chan int)
    dst := make(chan int, 3)
    // 子协程A 发送0~9的数字
    go func() {
        defer close(src)
        for i := 0; i < 10; i++ {
                src <- i
        }
    }()
    // 子协程B 获得协程A发送的数字 然后进行平方操作
    go func() {
        defer close(dst)
        for i := range src {
                dst <- i * i
        }
    }()
    // 主协程 获得B协程计算的结果 然后打印输出
    for i := range dst {
        fmt.Println(i)
    }
}
```
- 使用`<-`向`channel`发送消息，使用`for range`获取消息；
- 有创建就有关闭，使用`defer close(chennel)`关闭通道；
- 根据生产者的生产速度和消费者的消费速度，灵活的使用有缓冲和无缓冲通道，可以有效解决生产消费的不均衡的问题。
### 通过共享内存通信（Lock）
```go
var (
    x    int64      // 共享内存
    lock sync.Mutex // 互斥锁
)

func addWithLock() {
    for i := 0; i < 5000; i++ {
        // 临界区保护
        lock.Lock() // 加锁
        x += 1
        lock.Unlock() // 解锁
    }
}

func addWithoutLock() {
    for i := 0; i < 5000; i++ {
        // 临界区无保护
        x += 1
    }
}
```
- 不对临界区进行保护，输出结果未知；
- 使用互斥锁来保证并发安全；
## 3 WaitGroup
```go
const goRoutineNum = 5

func hello(i int) {
    fmt.Println("hello, goroutine: " + fmt.Sprint(i))
}

func HelloGoRoutine() {
    var wg sync.WaitGroup
    wg.Add(goRoutineNum) // 添加协程数
    for i := 0; i < goRoutineNum; i++ {
        go func(j int) {
            defer wg.Done() // 表示一个协程完成
            hello(j)
        }(i)
    }
    wg.Wait() // 阻塞
}
```
- `Add`添加协程数；
- `Done`表示一个协程的任务完成；
- `Wait`阻塞，保证子协程完成前主线程不退出；
