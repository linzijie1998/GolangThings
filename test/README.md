# Go语言的测试
回归测试：终端进行测试

集成测试：功能维度的测试，某个接口的测试

单元测试：单独的函数模块进行测试
## 1 Go语言中的测试规则
- 所有的测试文件以`_test.go`结尾；
- 测试函数为`func TestXxx(*testing.T)`的形式；
- 初始化逻辑放到`TestMain`中；
```go
func Func1() {
    // ...
}

func TestFunc1(t *testing.T) {
    // ...
}

func TestMain(m *testing.M) {
    // 测试前：数据装载、配置初始化等
    code := m.Run()
    // 测试后：释放资源等
    os.Exit(code)
}
```
## 2 单元测试

输入 -> 测试单元（函数、模块...）-> 输出

将输出和期望进行校对，保证质量、提升效率。
### 2.1 assert
获取第三方的assert包
```shell
go get -u "github.com/stretchr/testify/assert"
```
测试SayHello函数
```go
// hello.go
func SayHello() string {
	return "GoodBye!"
}
// hello_test.go
func TestSayHello(t *testing.T) {
	output := SayHello()
	expectOutput := "Hello!"
	assert.Equal(t, expectOutput, output)
}
```
### 2.2 覆盖率
测试判断是否及格的函数
```go
// judge.go
func JudgePassLine(score int16) bool {
	if score >= 60 {
		return true
	}
	return false
}
// judge_test.go
func TestJudgePassLine(t *testing.T) {
	isPass := JudgePassLine(70)
	assert.Equal(t, true, isPass)
}
```
使用命令行进行测试：
```shell
go test judge_test.go judge.go --cover
command-line-arguments  0.351s  coverage: 66.7% of statements
```
添加测试用例使其覆盖率到100%：
```go
// judge_test.go
func TestJudgePassLineTrue(t *testing.T) {
	isPass := JudgePassLine(70)
	assert.Equal(t, true, isPass)
}

func TestJudgePassLineFalse(t *testing.T) {
	isPass := JudgePassLine(50)
	assert.Equal(t, false, isPass)
}
```
### 2.3 依赖
一个复杂项目的依赖可能会有文件、数据库、缓存等强依赖。幂等：一个case重复运行很多次它的结果应该是一致的；稳定指单元测试相互隔离，能在任何时间、任何函数进行独立运行。

下面是一个文件的强依赖示例：
```go
// file.go
func ReadFirstLine() string {
	open, err := os.Open("test1.txt")
	defer open.Close()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func ProcessFirstLine() string {
	line := ReadFirstLine()
	dstLine := strings.ReplaceAll(line, "11", "00")
	return dstLine
}
// file_test.go
func TestProcessFirstLine(t *testing.T) {
	firstLine := ProcessFirstLine()
	assert.Equal(t, "line00", firstLine)
}
```
### 2.4 Mock
开源Mock测试包 monkey：https://github.com/bouk/monkey

快速Mock函数：为函数或者方法打桩从而将A函数替换为B函数

下面使用monkey的函数将ReadFirstLine函数替换为打桩的函数：
```go
// 对ReadFirstLine函数进行Mock，不再对本地文件依赖
func TestProcessFirstLineWithMock(t *testing.T) {
	monkey.Patch(ReadFirstLine, func() string {
		return "line110"
	})
	defer monkey.Unpatch(ReadFirstLine)
	line := ProcessFirstLine()
	assert.Equal(t, "line000", line)
}
```
### 2.5 单元测试Tips
- 一般覆盖率：50%~60%，较高覆盖率：80%+；
- 测试分支相互独立、全面覆盖；
- 测试单元粒度足够小，函数单一职责。

## 3. 基准测试
优化代码，需要对当前代码分析，Go语言内置了测试框架提供以基准测试。
```go
// benchmark.go
var ServerIndex [10]int

func InitServerIndex() {
	for i := 0; i < 10; i++ {
		ServerIndex[i] = i + 100
	}
}

func Select() int {
	return ServerIndex[rand.Intn(10)]
}
// benchmark_test.go
func BenchmarkSelect(b *testing.B) {
    // 串行测试
	InitServerIndex()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Select()
	}
}

func BenchmarkSelectParallel(b *testing.B) {
    // 并行测试
	InitServerIndex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Select()
		}
	})
}
```
`rand`函数为了保证全局随机性和并发安全，降低了并发的性能，因此结果上看`Select`函数的并行性能不如串行的性能。

优化方法：
```go
// go get -u "github.com/bytedance/gopkg"
func FastSelect() int {
	return ServerIndex[fastrand.Intn(10)]
}
```
