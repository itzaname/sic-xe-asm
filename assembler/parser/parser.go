package parser

import (
	"bufio"
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"fmt"
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

	g := graph.New()
	p := Parser{
		file:      file,
		scanner:   bufio.NewScanner(file),
		nodeGraph: &g,
	}
	return p, p.Parse()
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

		// Generate a string array of the line
		token, err := p.tokenizeLine(line)
		if err != nil {
			return fmt.Errorf("line %d: %s", lineNum, err.Error())
		}

		if p.isComment(token) {
			lineNum++
			continue
		}

		// Create graph node
		node, err := p.nodeFromToken(token)
		if err != nil {
			return fmt.Errorf("line %d: %s", lineNum, err.Error())
		}

		// Manually set needed debug info
		if item, ok := node.(*graph.InstructionNode); ok {
			item.Debug.Line = lineNum
		}
		if item, ok := node.(*graph.DirectiveNode); ok {
			item.Debug.Line = lineNum
		}

		p.nodeGraph.Append(node)

		lineNum++
	}

	_, err := p.nodeGraph.ResolveLiterals()
	if err != nil {
		return err
	}

	if _, err := p.nodeGraph.LinkNodes(); err != nil {
		return err
	}

	p.nodeGraph.UpdateAddr()

	return nil
}
