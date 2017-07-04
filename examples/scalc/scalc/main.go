package main

import (
	"bufio"
	"fmt"
	"github.com/adamcolton/parlex/examples/scalc"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "scalc"
	app.Usage = "Reverse Polish Notation"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			return interactive()
		}
		r := scalc.Eval(strings.Join(c.Args(), " "))
		if r == nil {
			return fmt.Errorf("Failed to parse")
		}
		fmt.Println(format(r))
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "parse",
			Usage: "show the parse tree for an operation",
			Action: func(c *cli.Context) error {
				r := scalc.Parse(strings.Join(c.Args(), " "))
				if r == nil {
					return fmt.Errorf("Failed to parse")
				}
				fmt.Println(r)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
	}
}

func interactive() error {
	var stack string
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		if input == "exit\n" {
			return nil
		}
		r := scalc.Eval(stack + " " + input)
		if r == nil {
			fmt.Println("Bad input")
			continue
		}
		stack = format(r)
		fmt.Println(stack)
	}
	return nil
}

func format(stack []scalc.Pfloat) string {
	strs := make([]string, len(stack))
	for i, p := range stack {
		strs[i] = p.String()
	}
	return strings.Join(strs, " ")
}
