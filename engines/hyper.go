package engines

import (
	"fmt"

	"github.com/flier/gohs/hyperscan"
)

type Hyper struct {
	engine hyperscan.BlockDatabase
}

func (h *Hyper) Compile(name, regexp string) error {
	flags, err := hyperscan.ParseCompileFlag("L")
	if err != nil {
		return err
	}

	h.engine, err = hyperscan.NewBlockDatabase(&hyperscan.Pattern{Expression: regexp, Flags: flags})
	if err != nil {
		return err
	}
	return nil
}

func (h *Hyper) CompileMany(expressions map[string]string) error {
	flags, err := hyperscan.ParseCompileFlag("L")
	if err != nil {
		return err
	}

	patterns := []*hyperscan.Pattern{}
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%v>%s)", name, regexpr)
		p := &hyperscan.Pattern{Expression: regexpr, Flags: flags}
		patterns = append(patterns, p)
	}

	h.engine, err = hyperscan.NewBlockDatabase(patterns...)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hyper) CountAll(data *[]byte, display bool) (int, error) {
	matches := h.engine.FindAllIndex(*data, -1)
	if display {
		for i, m := range matches {
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[m[0]:m[1]]))
		}
	}
	return len(matches), nil
}
