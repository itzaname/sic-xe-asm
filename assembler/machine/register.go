package machine

import (
	"fmt"
	"strings"
)

type Register struct {
	Name   string
	Number uint8
}

var regList = []Register{
	{Name: "A", Number: 0},
	{Name: "X", Number: 1},
	{Name: "L", Number: 2},
	{Name: "B", Number: 3},
	{Name: "S", Number: 4},
	{Name: "T", Number: 5},
	{Name: "F", Number: 6},
}

func RegisterByNumber(num uint8) (*Register, error) {
	for i := 0; i < len(regList); i++ {
		if regList[i].Number == num {
			return &regList[i], nil
		}
	}

	return nil, fmt.Errorf("invalid register number: '%d'", num)
}

func RegisterByName(name string) (*Register, error) {
	name = strings.ToUpper(name)

	for i := 0; i < len(regList); i++ {
		if regList[i].Name == name {
			return &regList[i], nil
		}
	}

	return nil, fmt.Errorf("invalid register name: '%s'", name)
}
