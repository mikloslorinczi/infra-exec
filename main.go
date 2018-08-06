package main

import (
	"flag"
	"fmt"
)

var (
	command string
	outfile string
)

func main() {
	flag.Parse()
	fmt.Println("\nCommand Executor")
	fmt.Println("Command:", command)
	fmt.Println("File:", outfile)
}

func init() {
	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&outfile, "f", "outfile", "Otput file")
}
