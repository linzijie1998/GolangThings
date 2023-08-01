package animal

import "fmt"

// Animal Animal接口，有一个Info方法返回动物的基本信息
type Animal interface {
	Info() string
}

type Cat struct {
	Name string
	Age  int
}

func (c *Cat) Info() string {
	return fmt.Sprintf("这是一只猫，它的名字叫%s, 今年%d岁了.", c.Name, c.Age)
}

type Dog struct {
	Name string
	Age  int
}

func (d *Dog) Info() string {
	return fmt.Sprintf("这是一只狗，它的名字叫%s, 今年%d岁了.", d.Name, d.Age)
}
