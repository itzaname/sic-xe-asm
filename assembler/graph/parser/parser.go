package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

type Parser struct {
	file    *os.File
	scanner *bufio.Scanner
}

func New(path string) (Parser, error) {
	file, err := os.Open(path)
	if err != nil {
		return Parser{}, err
	}

	return Parser{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (p *Parser) readLine() (bool, string) {
	return p.scanner.Scan(), p.scanner.Text()
}

func (p *Parser) cleanLine(items []string) []string {
	output := []string{}
	for i := 0; i < len(items); i++ {
		oLine := items[i]
		oLine = strings.TrimSpace(oLine)
		if len(oLine) > 0 {
			output = append(output, oLine)
		}
	}

	return output
}

// strings.Fields would be the ideal thing to use
// but since it will remove the spaces it isn't suitable
// for data arguments where the space is intentional
func (p *Parser) splitLine(input string) []string {
	output := []string{}
	buffer := ""
	for i := 0; i < len(input); i++ {
		buffer += string(input[i])
		if unicode.IsSpace(rune(input[i])) {
			output = append(output, buffer)
			buffer = ""
		}
	}
	output = append(output, buffer)

	return output
}

var dataDelimiter = []uint8{0x43, 0x58} // C, X
func (p *Parser) isDataDelimiter(input uint8) bool {
	for i := 0; i < len(dataDelimiter); i++ {
		if dataDelimiter[i] == input {
			return true
		}
	}
	return false
}

func (p *Parser) tokenizeLine(input string) ([]string, error) {
	line := p.splitLine(input)

	tokens := []string{}
	dataMode := false
	dataBuffer := ""

	for i := 0; i < len(line); i++ {
		start := -1
		end := -1
		for x := 0; x < len(line[i]); x++ {
			if !dataMode && unicode.IsSpace(rune(line[i][x])) {
				continue
			}

			if len(line[i]) > x+1 {
				if p.isDataDelimiter(line[i][x]) && line[i][x+1] == 0x27 {
					dataMode = true
					if start < 0 {
						start = x
					}
					x += 2
				}
			}

			if start < 0 {
				start = x
			}
			end = x

			if dataMode && line[i][x] == 0x27 {
				final := dataBuffer + line[i][start:x+1]
				tokens = append(tokens, final)
				dataBuffer = ""
				dataMode = false
				start = -1
				break
			}
		}
		end += 1
		if start >= 0 && end >= 0 {
			if dataMode {
				dataBuffer = dataBuffer + line[i][start:end]
			} else {
				tokens = append(tokens, line[i][start:end])
			}
		}

	}

	return tokens, nil
}

func (p *Parser) Parse() error {

	for {
		read, line := p.readLine()
		if !read {
			break
		}

		oof, _ := p.tokenizeLine(line)
		log.Println(len(oof), oof)
	}

	return nil
}
