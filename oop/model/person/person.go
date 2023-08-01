package person

import "fmt"

type Person struct {
	name string // 姓名
	age  int    // 年龄，范围在[0, 150]
}

// NewPerson Person类的构造器，它对年龄字段进行了安全校验
func NewPerson(name string, age int) (*Person, error) {
	if age < 0 || age > 150 {
		return nil, fmt.Errorf("age range needs to be between 0 and 150, but found %d", age)
	}
	return &Person{name, age}, nil
}

// GetName Name字段的Get方法
func (p *Person) GetName() string {
	return p.name
}

// SetName Name字段的Set方法
func (p *Person) SetName(name string) {
	p.name = name
}

// GetAge Age字段的Get方法
func (p *Person) GetAge() int {
	return p.age
}

// SetAge Age字段的Set方法，对年龄字段进行了安全校验
func (p *Person) SetAge(age int) error {
	if age < 0 || age > 150 {
		return fmt.Errorf("age range needs to be between 0 and 150, but found %d", age)
	}
	p.age = age
	return nil
}
