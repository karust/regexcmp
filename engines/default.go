package engines

import (
	"fmt"
	"regexp"
)

type Default struct {
	engine *regexp.Regexp
}

func (re *Default) Compile(name, regexpr string) (err error) {
	re.engine, err = regexp.Compile(regexpr)
	if err != nil {
		return err
	}
	return nil
}

func (re *Default) CompileMany(expressions map[string]string) (err error) {
	finalRegexp := ""
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%v>%s)|", name, regexpr)
		finalRegexp += regexpr
	}
	finalRegexp = finalRegexp[:len(finalRegexp)-1]

	re.engine, err = regexp.Compile(finalRegexp)
	if err != nil {
		return err
	}
	return nil
}

func (re *Default) CountAll(data *[]byte, display bool) (int, error) {
	matches := re.engine.FindAllIndex(*data, -1)
	if display {
		for i, m := range matches {
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[m[0]:m[1]]))
		}
	}
	return len(matches), nil
}
