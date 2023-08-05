package engines

import (
	"fmt"

	"github.com/GRbit/go-pcre"
)

type Pcre struct {
	engine pcre.Regexp
}

func (re *Pcre) Compile(name, regexpr string) (err error) {
	re.engine, _ = pcre.CompileJIT(regexpr, 0, pcre.STUDY_JIT_COMPILE)
	return nil
}

func (re *Pcre) CompileMany(expressions map[string]string) (err error) {
	finalRegexp := ""
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%v>%s)|", name, regexpr)
		finalRegexp += regexpr
	}
	finalRegexp = finalRegexp[:len(finalRegexp)-1]

	re.engine, _ = pcre.CompileJIT(finalRegexp, 0, pcre.STUDY_JIT_COMPILE)
	return nil
}

func (re *Pcre) CountAll(data *[]byte, display bool) (int, error) {
	matches := re.engine.FindAllIndex(*data, pcre.NOTEMPTY)
	if display {
		for i, m := range matches {
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[m[0]:m[1]]))
		}
	}
	return len(matches), nil
}
