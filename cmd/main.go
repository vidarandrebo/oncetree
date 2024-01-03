package main

import (
	"fmt"

	"github.com/vidarandrebo/once-tree/pkg/greeter"
)

func main() {
	test := greeter.NewGreeting()
	fmt.Println(test.Text)
}
