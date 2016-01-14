package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(`{ "Error": "something broke" }`)
	os.Exit(17)
}
