package engines

import (
	"errors"
	"fmt"

	"github.com/hillu/go-yara/v4"
)

type Yara struct {
	compiledRules *yara.Rules
}

func (y *Yara) Compile(name, regexp string) error {
	compiler, err := yara.NewCompiler()
	if err != nil {
		return err
	}
	r := y.ruleFromRegexp(name, regexp)
	err = compiler.AddString(r, "test_namespace")
	if err != nil {
		return err
	}

	y.compiledRules, err = y.compileRules(compiler)
	if err != nil {
		return err
	}
	return nil
}

func (y *Yara) CompileMany(expressions map[string]string) error {
	compiler, err := yara.NewCompiler()
	if err != nil {
		return err
	}

	for name, regexpr := range expressions {
		r := y.ruleFromRegexp(name, regexpr)
		err = compiler.AddString(r, "test_namespace")
		if err != nil {
			return err
		}
	}

	y.compiledRules, err = y.compileRules(compiler)
	if err != nil {
		return err
	}
	return nil
}

func (y *Yara) CountAll(data *[]byte, display bool) (int, error) {
	matchRules, err := y.yaraScan(*data, y.compiledRules)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, m := range matchRules {
		count += len(m.Strings)
		if display {
			for i, match := range m.Strings {
				fmt.Printf("\t%v)[%v] %v\n", i+1, match.Name, string(match.Data))
			}
		}
	}
	return count, nil
}

func (y *Yara) ruleFromRegexp(name, expr string) string {
	return fmt.Sprintf(`
	rule %s 
	{ 
		strings: 
			$re1 = /%s/ fullword
		condition: 
			$re1 
	}`, name, expr,
	)
}

func (y *Yara) compileRules(compiler *yara.Compiler) (rules *yara.Rules, err error) {
	rules, err = compiler.GetRules()
	if err != nil {
		return nil, errors.New("Failed to compile rules")
	}
	return rules, err
}

func (y *Yara) yaraScan(content []byte, rules *yara.Rules) (match yara.MatchRules, err error) {
	sc, err := yara.NewScanner(rules)
	if err != nil {
		return nil, err
	}

	var m yara.MatchRules
	err = sc.SetCallback(&m).ScanMem(content) //.ScanFile("./test_data/test_1")
	return m, err
}
