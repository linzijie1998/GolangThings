package main

import (
	"fmt"

	"github.com/linzijie1998/GolangThings/oop/model/animal"
)

func main() {
	// t, _ := person.NewTeacher("张三", 18)
	// fmt.Printf("name: %s, age: %d\n", t.GetName(), t.GetAge())

	var cat, dog animal.Animal

	cat = &animal.Cat{
		Name: "kitty",
		Age:  2,
	}
	fmt.Println(cat.Info())

	dog = &animal.Dog{
		Name: "sally",
		Age:  3,
	}
	fmt.Println(dog.Info())
}
