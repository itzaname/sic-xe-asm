package parser

import (
	"bufio"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"log"
	"os"
)

type Parser struct {
	nodeGraph *graph.Graph
	file      *os.File
	scanner   *bufio.Scanner
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

func (p *Parser) Graph() *graph.Graph {
	return p.nodeGraph
}

func (p *Parser) Parse() error {
	return p.parseScanner()
}

// Private stuff no looking plz
func (p *Parser) readLine() (bool, string) {
	return p.scanner.Scan(), p.scanner.Text()
}

func (p *Parser) parseScanner() error {
	lineNum := 1
	for {
		read, line := p.readLine()
		if !read {
			break
		}

		token, err := p.tokenizeLine(line)
		if err != nil {
			return err
		}

		log.Println(token)
		log.Println(p.nodeFromToken(token))

		lineNum++
	}

	return nil
}
