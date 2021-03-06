package machine

import "fmt"

type HeaderRecord struct {
	Name        string
	BaseAddress int
	Length      int
}

func (hr *HeaderRecord) String() string {
	name := hr.Name
	nameLen := len(name)
	for i := nameLen; i < 6; i++ {
		name = name + " "
	}
	return "H" + name + fmt.Sprintf("%.6X", hr.BaseAddress) + fmt.Sprintf("%.6X", hr.Length)
}

type TextRecord struct {
	Start  int
	Length int
	Object string
}

func (tx *TextRecord) String() string {
	return "T" + fmt.Sprintf("%.6X", tx.Start) + fmt.Sprintf("%.2X", tx.Length) + tx.Object
}

type EndRecord struct {
	Start int
}

func (er *EndRecord) String() string {
	return "E" + fmt.Sprintf("%.6X", er.Start)
}

type ModificationRecord struct {
	Address int
	Length  int
}

func (mr *ModificationRecord) String() string {
	return "M" + fmt.Sprintf("%.6X", mr.Address) + fmt.Sprintf("%.2X", mr.Length)
}
