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

		if evaler.IsWholeNumber(res) {
			fmt.Printf("%d\n", int64(res))
		} else {
			fmt.Printf("%.2f\n", res)
		}
	}
}

func main() {
	repl()
}
