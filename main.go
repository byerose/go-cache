package main

import "fmt"

type people struct {
	name string
	age  int
}

func main() {
	p := people{
		"ii",
		12,
	}
	var i interface{} = p
	o := i.(people)
	fmt.Printf("%T,%T,%T,%p,%p,%p", p, i, o, &p, &i, &o)
	fmt.Print(len("ni"))
}
