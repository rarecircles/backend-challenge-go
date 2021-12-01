package eth

import (
	"fmt"
	"strings"
)

type Log struct {
	Address []byte   `json:"address,omitempty"`
	Topics  [][]byte `json:"topics,omitempty"`
	Data    []byte   `json:"data,omitempty"`
	// supplement
	Index      uint32 `json:"index,omitempty"`
	BlockIndex uint32 `json:"blockIndex,omitempty"`
}

type LogEventDef struct {
	Name       string
	Parameters []*LogParameter
}

type LogParameter struct {
	Name     string
	TypeName string
	Indexed  bool
}

type LogEvent struct {
	Def  *LogEventDef
	Data interface{}
}

func (p *LogParameter) GetName(index int) string {
	name := p.Name
	if name == "" {
		name = fmt.Sprintf("unamed%d", index+1)
	}

	return name
}

func (l *LogEventDef) logID() []byte {
	return Keccak256([]byte(l.Signature()))
}

func (l *LogEventDef) Signature() string {
	var args []string
	for _, parameter := range l.Parameters {
		args = append(args, parameter.TypeName)
	}

	return fmt.Sprintf("%s(%s)", l.Name, strings.Join(args, ", "))
}

func (l *LogEventDef) String() string {
	var args []string
	for i, parameter := range l.Parameters {
		indexed := ""
		if parameter.Indexed {
			indexed = " indexed"
		}

		args = append(args, fmt.Sprintf("%s%s %s", parameter.TypeName, indexed, parameter.GetName(i)))
	}

	return fmt.Sprintf("%s(%s)", l.Name, strings.Join(args, ", "))
}
