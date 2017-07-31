package main

import (
	"fmt"
	"github.com/adamcolton/parlex/examples/parlexjson"
	"os"
	"strings"
)

func main() {
	s := strings.Join(os.Args[1:], " ")
	out, err := parlexjson.Format(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(out)
}
