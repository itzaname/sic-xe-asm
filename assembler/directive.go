package assembler

import (
	"ci.itzana.me/itzaname/sic-xe-asm/assembler/parser/graph"
	"encoding/binary"
	"fmt"
	"strconv"
)

type directiveManager struct {
	*Assembler
}

var dirFuncMap map[string]func(node *graph.DirectiveNode) ([]byte, bool, error)

func (dir *directiveManager) Start(node *graph.DirectiveNode) ([]byte, bool, error) {
	if len(node.Name) > 6 {
		return nil, false, fmt.Errorf("program name too long. %d characters when maximum is 6", len(node.Name))
	}
	dir.Assembler.obj.Header.Name = node.Name

	addr, err := strconv.Atoi(node.Data.(string))
	if err != nil {
		return nil, false, err
	}
	dir.Assembler.obj.Header.BaseAddress = addr

	return nil, false, nil
}

func (dir *directiveManager) Base(node *graph.DirectiveNode) ([]byte, bool, error) {
	target := node.Data.(graph.Node)
	dir.flags.base = true
	dir.flags.baseAddr = target.Address()
	return nil, false, nil
}

func (dir *directiveManager) NoBase(node *graph.DirectiveNode) ([]byte, bool, error) {
	dir.flags.base = false
	dir.flags.baseAddr = 0
	return nil, false, nil
}

func (dir *directiveManager) Byte(node *graph.DirectiveNode) ([]byte, bool, error) {
	data := node.Data.(*graph.Storage)
	return data.Data.([]byte), true, nil
}

func (dir *directiveManager) Word(node *graph.DirectiveNode) ([]byte, bool, error) {
	data := node.Data.(*graph.Storage)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(data.Data.(int)))
	return bs[1:], false, nil
}

func (dir *directiveManager) ResB(node *graph.DirectiveNode) ([]byte, bool, error) {
	dir.flags.endModule = true
	return nil, false, nil
}

func (dir *directiveManager) ResW(node *graph.DirectiveNode) ([]byte, bool, error) {
	dir.flags.endModule = true
	return nil, false, nil
}

func (dir *directiveManager) End(node *graph.DirectiveNode) ([]byte, bool, error) {
	if node.Data != nil {
		target := node.Data.(graph.Node)
		dir.obj.End.Start = target.Address()
	}
	return nil, false, nil
}

func (dir *directiveManager) Null(node *graph.DirectiveNode) ([]byte, bool, error) {
	return nil, false, nil
}

func (asm *Assembler) handleDirective(node *graph.DirectiveNode) ([]byte, bool, error) {
	if dirFuncMap == nil {
		dirManager := directiveManager{asm}
		dirFuncMap = map[string]func(node *graph.DirectiveNode) ([]byte, bool, error){}

		// Directive functions
		dirFuncMap["START"] = dirManager.Start
		dirFuncMap["BASE"] = dirManager.Base
		dirFuncMap["NOBASE"] = dirManager.Base
		dirFuncMap["BYTE"] = dirManager.Byte
		dirFuncMap["WORD"] = dirManager.Word
		dirFuncMap["RESB"] = dirManager.ResB
		dirFuncMap["RESW"] = dirManager.ResW
		dirFuncMap["END"] = dirManager.End
		dirFuncMap["LTORG"] = dirManager.Null
		dirFuncMap["EQU"] = dirManager.Null
	}

	if f, ok := dirFuncMap[node.Directive.Name]; ok {
		return f(node)
	}

	return nil, false, fmt.Errorf("unhandled directive '%s'", node.Directive.Name)
}
