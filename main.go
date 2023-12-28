package main

import (
	"fmt"
	"github.com/Neal-C/interpreter-in-go/repl"
	"os"
	"os/user"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s ! This is the monkey programming language ! \n", currentUser.Username)
	fmt.Printf("Start typing commands \n")
	repl.Start(os.Stdin, os.Stdout)
}
