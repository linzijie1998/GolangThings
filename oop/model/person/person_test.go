package person_test

import (
	"testing"

	"github.com/linzijie1998/GolangThings/oop/model/person"
)

const (
	name       = "张三"
	validAge   = 25
	invalidAge = -8
)

func TestNewPerson(t *testing.T) {
	var err error

	_, err = person.NewPerson(name, validAge)
	if err != nil {
		t.Errorf("creating valid Person returned an error: %s", err.Error())
	}

	_, err = person.NewPerson(name, invalidAge)
	if err != nil {
		t.Errorf("creating invalid Person returned an error: %s", err.Error())
	}
}
