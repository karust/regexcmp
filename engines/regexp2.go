package engines

import (
	"bytes"
	"fmt"

	regexp2 "github.com/dlclark/regexp2"
)

type Regexp2 struct {
	engine *regexp2.Regexp
}

func (r2 *Regexp2) Compile(name, regexpr string) (err error) {
	r2.engine, err = regexp2.Compile(regexpr, 0x0200)
	if err != nil {
		return err
	}
	return nil
}

func (r2 *Regexp2) CompileMany(expressions map[string]string) (err error) {
	finalRegexp := ""
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%v>%s)|", name, regexpr)
		finalRegexp += regexpr
	}
	finalRegexp = finalRegexp[:len(finalRegexp)-1]
	r2.engine, err = regexp2.Compile(finalRegexp, 0x0200)
	if err != nil {
		return err
	}
	return nil
}

func (r2 *Regexp2) CountAll(data *[]byte, display bool) (int, error) {
	var matches [][]rune
	m, _ := r2.engine.FindRunesMatch(bytes.Runes(*data))
	for m != nil {
		matches = append(matches, m.Runes())
		m, _ = r2.engine.FindNextMatch(m)
	}

	if display {
		for i, m := range matches {
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[m[0]:m[1]]))
		}
	}
	return len(matches), nil
}
