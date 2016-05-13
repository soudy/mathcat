package main

import (
	"flag"
	"fmt"
	"github.com/soudy/mathcat"
	"gopkg.in/readline.v1"
)

var precision = flag.Int("precision", 2, "decimal precision used in results")

func repl() {
	p := mathcat.New()
	rl, err := readline.New("mathcat> ")
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

		if mathcat.IsWholeNumber(res) {
			fmt.Printf("%d\n", int64(res))
		} else {
			fmt.Printf("%.*f\n", *precision, res)
		}
	}
}

func main() {
	flag.Parse()
	repl()
}
