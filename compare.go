package main

import (
	"fmt"
	"runtime"
	"time"
)

type RegexpEngine interface {
	Compile(string, string) error
	CompileMany(map[string]string) error
	CountAll(*[]byte, bool) (int, error)
}

var allRegexps = make(map[string]string)

func initRegexps() {
	allRegexps["email"] = `(?P<name>[-\w\d\.]+?)(?:\s+at\s+|\s*@\s*|\s*(?:[\[\]@]){3}\s*)(?P<host>[-\w\d\.]*?)\s*(?:dot|\.|(?:[\[\]dot\.]){3,5})\s*(?P<domain>\w+)`
	allRegexps["bitcoin"] = `\b([13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[ac-hj-np-zAC-HJ-NP-Z02-9]{11,71})`
	allRegexps["ssn"] = `\d{3}-\d{2}-\d{4}`
	allRegexps["uri"] = `[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?`
	allRegexps["tel"] = `\+\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`

	if config.GenerateNonMatching {
		keys := []string{}
		for name := range allRegexps {
			keys = append(keys, name)
		}
		for _, k := range keys {
			allRegexps["non_matching_"+k] = allRegexps[k] + k
		}
	}

	for i := 0; i < config.NumberNonMatching; i++ {
		allRegexps[fmt.Sprintf("non_matching%d", i)] = fmt.Sprintf(`[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?%d`, i)
	}

	for i := 0; i < config.NumberMatching; i++ {
		allRegexps[fmt.Sprintf("matching%d", i)] = `[\w]+://[^/\s?#]+[^\s?#]+(?:\?[^\s#]*)?(?:#[^\s]*)?`
	}
}

func runSingle(engine RegexpEngine, data *[]byte) error {
	var globM1, globM2 runtime.MemStats
	totalCount, totalElapsed := 0, 0

	runtime.GC()
	runtime.ReadMemStats(&globM1)

	for name, expr := range allRegexps {
		var m1, m2 runtime.MemStats
		runtime.ReadMemStats(&m1)

		err := engine.Compile(name, expr)
		if err != nil {
			return err
		}

		start := time.Now()

		count, err := engine.CountAll(data, config.IsDisplayOutput)
		if err != nil {
			return err
		}
		totalCount += count

		elapsed := time.Since(start)
		totalElapsed += int(elapsed)
		runtime.ReadMemStats(&m2)

		fmt.Printf("  [%v] count=%v, mem=%.2fKB, time=%v \n",
			name, count, float64(m2.TotalAlloc-m1.TotalAlloc)/1000, time.Duration(elapsed).Round(time.Microsecond))
	}

	runtime.ReadMemStats(&globM2)

	fmt.Printf("Total. Counted: %v, Memory: %.2fMB, Duration: %v \n",
		totalCount, float64(globM2.TotalAlloc-globM1.TotalAlloc)/1000/1000, time.Duration(totalElapsed).Round(time.Microsecond))
	return nil
}

func runGroup(engine RegexpEngine, data *[]byte) error {
	var globM1, globM2, globM3 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&globM1)

	err := engine.CompileMany(allRegexps)
	if err != nil {
		return err
	}
	runtime.ReadMemStats(&globM2)

	start := time.Now()
	count, err := engine.CountAll(data, config.IsDisplayOutput)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	runtime.ReadMemStats(&globM3)

	fmt.Printf("  [%v regexps] count=%v, compiled_mem=%.2fMB, mem=%.2fMB, time=%v \n",
		len(allRegexps), count,
		float64(globM2.TotalAlloc-globM1.TotalAlloc)/1000/1000,
		float64(globM3.TotalAlloc-globM1.TotalAlloc)/1000/1000,
		time.Duration(elapsed).Round(time.Microsecond),
	)
	return nil
}
