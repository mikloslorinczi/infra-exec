package main

import (
	"fmt"
)

func main() {
	mydb := ConnectJSONDB("db.json")
	err := mydb.open()
	if err != nil {
		fmt.Println(err)
	}
	err = mydb.save()
	if err != nil {
		fmt.Println(err)
	}
}
