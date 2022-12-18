package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
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
	if bytes.Equal(buffer, []byte("tEXt")) {
		rawText := make([]byte, length)
		if _, err := chunk.Read(rawText); err != nil {
			return err
		}
		fmt.Println(string(rawText))
	}
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

func textChunk(text string) (io.Reader, error) {
	byteData := []byte(text)
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, int32(len(byteData))); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString("tEXt"); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(byteData); err != nil {
		return nil, err
	}
	crc := crc32.NewIEEE()
	if _, err := io.WriteString(crc, "tEXt"); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, crc.Sum32()); err != nil {
		return nil, err
	}
	return &buffer, nil
}

func main() {
	file, err := os.Open("./Lenna.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	newFile, err := os.Create("./Lenna2.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer newFile.Close()

	chunks, err := readChunks(file)
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(newFile, "\x89PNG\r\n\x1a\n"); err != nil {
		log.Fatalln(err)
	}
	if _, err := io.Copy(newFile, chunks[0]); err != nil {
		log.Fatalln(err)
	}
	txt, err := textChunk("hogehoge")
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := io.Copy(newFile, txt); err != nil {
		log.Fatalln(err)
	}
	for _, chunk := range chunks[1:] {
		if _, err := io.Copy(newFile, chunk); err != nil {
			log.Fatalln(err)
		}
	}

	newChunks, err := readChunks(newFile)
	if err != nil {
		panic(err)
	}
	for _,  chunk := range newChunks {
		dumpChunk(chunk)
	}

}
