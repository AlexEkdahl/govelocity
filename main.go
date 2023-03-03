package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("hej")
		time.Sleep(time.Millisecond * 2000)
	}
}
