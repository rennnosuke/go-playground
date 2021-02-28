package main

import (
	"encoding/csv"
	"os"
)

func main() {
	encoder := csv.NewWriter(os.Stdout)
	defer encoder.Flush()
	if err := encoder.Write([]string{"write line"}); err != nil {
		panic(err)
	}
}
