package main

import (
	"fmt"
	"os"
	"os/user"
	"mua/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("MUA-lang v0.0.1 | Welcome `%s` on linux!\n",
		user.Username)
	fmt.Printf("Type \"help\" for more information.\n")
	repl.Start(os.Stdin, os.Stdout)
}