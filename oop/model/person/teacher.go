package person

type Teacher struct {
	*Person
}

func NewTeacher(name string, age int) (*Teacher, error) {
	p, err := NewPerson(name, age)
	if err != nil {
		return nil, err
	}
	return &Teacher{p}, nil
}

func (t *Teacher) GetAge() int {
	return t.age + 100000
}
