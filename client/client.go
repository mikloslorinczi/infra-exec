package main

import (
	"flag"
	"fmt"
)

var nodeName string

func main() {

	flag.Parse()

	fmt.Printf("Hello, my name is %v\n", nodeName)

}

func init() {
	flag.StringVar(&nodeName, "name", "", "Name of node")
}
