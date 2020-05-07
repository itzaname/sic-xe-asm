package machine

import "fmt"

type Directive struct {
	Name     string
	Storage  bool
	Resolved bool
}

var directiveList = []Directive{
	{Name: "BYTE", Storage: true},
	{Name: "WORD", Storage: true},
	{Name: "RESB", Storage: true},
	{Name: "RESW", Storage: true},
	{Name: "START", Storage: false},
	{Name: "BASE", Storage: false, Resolved: true},
	{Name: "NOBASE", Storage: false},
	{Name: "END", Storage: false, Resolved: true},
	{Name: "LTORG", Storage: false},
}

func DirectiveByName(name string) (*Directive, error) {
	for i := 0; i < len(directiveList); i++ {
		if directiveList[i].Name == name {
			return &directiveList[i], nil
		}
	}

	return nil, fmt.Errorf("invalid directive: '%s'", name)
}
