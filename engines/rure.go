package engines

import (
	"fmt"

	"github.com/BurntSushi/rure-go"
)

type Rure struct {
	engine *rure.Regex
}

func (re *Rure) Compile(name, regexpr string) (err error) {
	re.engine, err = rure.Compile(regexpr)
	if err != nil {
		return err
	}
	return nil
}

func (re *Rure) CompileMany(expressions map[string]string) (err error) {
	finalRegexp := ""
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%v>%s)|", name, regexpr)
		finalRegexp += regexpr
	}
	finalRegexp = finalRegexp[:len(finalRegexp)-1]

	re.engine, err = rure.Compile(finalRegexp)
	if err != nil {
		return err
	}
	return nil
}

func (re *Rure) CountAll(data *[]byte, display bool) (int, error) {
	matches := re.engine.FindAllBytes(*data)
	if display {
		for i := 0; i < len(matches)/2; i++ {
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[matches[i*2]:matches[i*2+1]]))
		}
	}
	return len(matches) / 2, nil
}
