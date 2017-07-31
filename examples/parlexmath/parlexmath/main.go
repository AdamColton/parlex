package main

import (
	"fmt"
	"github.com/adamcolton/parlex/examples/parlexmath"
	"os"
	"strings"
)

func main() {
	v, err := parlexmath.Eval(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Println(v)
}
