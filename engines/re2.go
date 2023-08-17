package engines

import (
	"fmt"

	re2 "github.com/wasilibs/go-re2"
)

type Re2 struct {
	engine *re2.Regexp
}

func (re *Re2) Compile(name, regexpr string) (err error) {
	re.engine, err = re2.Compile(regexpr)
	if err != nil {
		return err
	}
	return nil
}

func (re *Re2) CompileMany(expressions map[string]string) (err error) {
	finalRegexp := ""
	for name, regexpr := range expressions {
		regexpr = fmt.Sprintf("(?P<%s>%s)|", name, regexpr)
		finalRegexp += regexpr
	}
	finalRegexp = finalRegexp[:len(finalRegexp)-1]
	re.engine, err = re2.Compile(finalRegexp)
	if err != nil {
		return err
	}
	return nil
}

func (re *Re2) CountAll(data *[]byte, display bool) (int, error) {
	matches := re.engine.FindAllIndex(*data, -1)

	if display {
		for i, m := range matches {
			// for j, name := range re.engine.SubexpNames() {
			// 	if j != 0 && name != "" && string(matches[i][j]) != "" {
			// 		fmt.Println(name, string(matches[i][j]))
			// 	}
			// }
			fmt.Printf("\t%v. %v\n", i+1, string((*data)[m[0]:m[1]]))
		}
	}
	return len(matches), nil
}
