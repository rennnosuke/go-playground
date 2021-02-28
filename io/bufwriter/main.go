package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	buf := bufio.NewWriter(os.Stdout)
	for i := range make([]int, 10) {
		if _, err := buf.WriteString(fmt.Sprintf("loop %d\n", i)); err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
		if err := buf.Flush(); err != nil {
			panic(err)
		}
	}
}
