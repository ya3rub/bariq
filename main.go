package main

import (
	"fmt"
	"os"
	"os/user"

	"bariq/repl"
)

var pl = fmt.Println

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, Starting now...\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
