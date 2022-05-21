package gob

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"
	"testing"
)

const mapSize = 2

func createMap() map[int]string {
	m := make(map[int]string, mapSize)
	for i := 0; i < mapSize; i++ {
		m[i] = strconv.Itoa(i)
	}
	return m
}

type Struct struct {
	ID   int64
	Name string
}

func BenchmarkGobEncodeMap(b *testing.B) {
	var buf bytes.Buffer
	m := createMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gob.NewEncoder(&buf).Encode(&m); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonEncodeMap(b *testing.B) {
	var buf bytes.Buffer
	m := createMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.NewEncoder(&buf).Encode(&m); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGobDecodeMap(b *testing.B) {
	m := createMap()

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&m); err != nil {
		b.Error(err)
	}
	bufBytes := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(bufBytes)
		if err := gob.NewDecoder(buf).Decode(&m); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonDecodeMap(b *testing.B) {
	m := createMap()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&m); err != nil {
		b.Error(err)
	}
	bufBytes := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(bufBytes)
		if err := json.NewDecoder(buf).Decode(&m); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGobEncodeStruct(b *testing.B) {
	var buf bytes.Buffer
	v := Struct{
		ID:   1,
		Name: "hoge",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := gob.NewEncoder(&buf).Encode(&v); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonEncodeStruct(b *testing.B) {
	var buf bytes.Buffer
	v := Struct{
		ID:   1,
		Name: "hoge",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.NewEncoder(&buf).Encode(&v); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGobDecodeStruct(b *testing.B) {
	v := Struct{
		ID:   1,
		Name: "hoge",
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&v); err != nil {
		b.Error(err)
	}
	bufBytes := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(bufBytes)
		if err := gob.NewDecoder(buf).Decode(&v); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsonDecodeStruct(b *testing.B) {
	v := Struct{
		ID:   1,
		Name: "hoge",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&v); err != nil {
		b.Error(err)
	}
	bufBytes := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(bufBytes)
		if err := json.NewDecoder(buf).Decode(&v); err != nil {
			b.Error(err)
		}
	}
}
