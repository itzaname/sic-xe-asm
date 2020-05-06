package parser

import (
	"fmt"
	"strconv"
)

func (p *Parser) readArrayStorageInput(input string) (interface{}, error) {
	if out, err := strconv.Atoi(input); err == nil {
		return out, nil
	}

	if len(input) < 4 {
		return nil, fmt.Errorf("expected data got: %s", input)
	}

	if input[0] == 'X' && input[1] == '\'' {
		return nil, fmt.Errorf("I can't handle this power")
	}

	if input[0] == 'C' && input[1] == '\'' {
		for i := 2; i < len(input); i++ {
			if input[i] == '\'' && input[i-1] != '\\' {
				fmt.Println(len(input), i)
				return input[2 : i-1], nil
			}
		}
	}

	return nil, fmt.Errorf("could not read input: %s", input)
}
