package main

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler"
	"fmt"
)

func main() {
	obj, err := assembler.GetObject("fuck")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(obj)
}
