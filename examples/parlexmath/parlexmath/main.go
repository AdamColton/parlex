package main

import (
	"fmt"
	"github.com/adamcolton/parlex/examples/parlexmath"
	"os"
	"strings"
)

func main() {
	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		fmt.Fprint(os.Stderr, "Please provide input to parse\n")
		return
	}
	v, err := parlexmath.Eval(input)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Println(v)
}
