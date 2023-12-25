package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
)

// CSV Upload Server sample.
// This server handle csv uploading (RFC4180: https://www.rfc-editor.org/rfc/rfc4180.html)
//
// curl -X POST -H "Content-Type: text/csv" -d "1,2,3,4,5" http://localhost:8080/

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// check Content-Type
		if v := r.Header.Get("Content-Type"); v != "text/csv" {
			http.Error(w, "invalid Content-Type: "+v, http.StatusBadRequest)
			return
		}
		// check Content-Length
		if r.ContentLength == 0 {
			http.Error(w, "invalid Content-Length: 0", http.StatusBadRequest)
			return
		}
		// check Method
		if r.Method != http.MethodPost {
			http.Error(w, "invalid Method: "+r.Method, http.StatusBadRequest)
			return
		}

		csvWriter := csv.NewReader(r.Body)
		// comment accept
		csvWriter.Comment = '#'
		// read all records
		records, err := csvWriter.ReadAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return json array
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		if err := e.Encode(records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
