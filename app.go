package main

import (
	"bytes"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func fixLineEndings(file string) (*bytes.Buffer, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] != '\n' {
			data[i] = '\n'
		}
	}

	return bytes.NewBuffer(data), nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s input output\n", filepath.Base(os.Args[0]))
		return
	}

	buffer, err := fixLineEndings(os.Args[1])
	if err != nil {
		fmt.Println("invalid file", err)
	}

	obj, err := assembler.GetObject(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ioutil.WriteFile(os.Args[2], []byte(obj), 0644); err != nil {
		fmt.Println(err)
		return
	}
}
