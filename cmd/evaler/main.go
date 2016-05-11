package main

import (
	"fmt"
	"github.com/soudy/evaler"
	"gopkg.in/readline.v1"
)

func repl() {
	p := evaler.New()
	rl, err := readline.New("evaler> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		res, err := p.Run(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%.1f\n", res)
	}
}

func main() {
	repl()
}
