package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func dumpChunk(chunk io.Reader) error {
	var length int32
	if err := binary.Read(chunk, binary.BigEndian, &length); err != nil {
		return err
	}
	buffer := make([]byte, 4)
	if _, err := chunk.Read(buffer); err != nil {
		return err
	}
	fmt.Printf("chunk '%v' (%d bytes)\n", string(buffer), length)
	return nil
}

func readChunks(f *os.File) ([]io.Reader, error) {
	var chunks []io.Reader

	if _, err := f.Seek(8, 0); err != nil {
		return nil, err
	}
	offset := int64(8)

	for {
		var length int32
		if err := binary.Read(f, binary.BigEndian, &length); err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(f, offset, int64(length)+12))
		offset, _ = f.Seek(int64(length+8), 1)
	}

	return chunks, nil
}

func main() {
	file, err := os.Open("./Lenna.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	chunks, err := readChunks(file)
	if err != nil {
		panic(err)
	}
	for _, chunk := range chunks {
		if err := dumpChunk(chunk); err != nil {
			panic(err)
		}
	}
}
