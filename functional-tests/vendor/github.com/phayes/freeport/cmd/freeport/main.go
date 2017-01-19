package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"strconv"
)

func main() {
	fmt.Println(strconv.Itoa(freeport.GetPort()))
}
