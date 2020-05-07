package main

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s input output\n", filepath.Base(os.Args[0]))
		return
	}

	obj, err := assembler.GetObject(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ioutil.WriteFile(os.Args[2], []byte(obj), 0644); err != nil {
		fmt.Println(err)
		return
	}
}
