package main

import (
	"encoding/json"
	"os"
)

func main() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	body := map[string]interface{}{
		"id":   1,
		"name": "test",
	}
	if err := encoder.Encode(body); err != nil {
		panic(err)
	}
}
