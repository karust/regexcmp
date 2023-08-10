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
	allRegexps["email0"] = `[\w\.+-]+@[\w\.-]+\.[\w\.-]+`
	allRegexps["email1"] = `[a-z0-9_\.-]+\@[\da-z\.-]+\.[a-z\.]{2,6}`
	allRegexps["email2"] = `([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x22([^\x0d\x22\x5c\x80-\xff]|\x5c[\x00-\x7f])*\x22)(\x2e([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x22([^\x0d\x22\x5c\x80-\xff]|\x5c[\x00-\x7f])*\x22))*\x40([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x5b([^\x0d\x5b-\x5d\x80-\xff]|\x5c[\x00-\x7f])*\x5d)(\x2e([^\x00-\x20\x22\x28\x29\x2c\x2e\x3a-\x3c\x3e\x40\x5b-\x5d\x7f-\xff]+|\x5b([^\x0d\x5b-\x5d\x80-\xff]|\x5c[\x00-\x7f])*\x5d))*`
	//allRegexps["email3"] = `(?P<name>[-\w\d\.]+?)(?:\s+at\s+|\s*@\s*|\s*(?:[\[\]@]){3}\s*)(?P<host>[-\w\d\.]*?)\s*(?:dot|\.|(?:[\[\]dot\.]){3,5})\s*(?P<domain>\w+)`

	allRegexps["bitcoin"] = `\b([13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[ac-hj-np-zAC-HJ-NP-Z02-9]{11,71})`
	allRegexps["ip"] = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	allRegexps["ssn"] = `\d{3}-\d{2}-\d{4}`

	//allRegexps["uri"] = `(?P<protocol>(?:[^:]+)s?)?:\/\/(?:(?P<user>[^:\n\r]+):(?P<pass>[^@\n\r]+)@)?(?P<host>(?:www\.)?(?:[^:\/\n\r]+))(?::(?P<port>\d+))?\/?(?P<request>[^?#\n\r]+)?\??(?P<query>[^#\n\r]*)?\#?(?P<anchor>[^\n\r]*)?`
	//allRegexps["telephone"] = `\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`

	for i := 0; i < config.NumberNonMatching; i++ {
		allRegexps[fmt.Sprintf("non_matching%v", i)] = `[a-z0-9_\.-]+\@[\da-z\.-]+\.[a-z\.]{2,6}1`
	}

	for i := 0; i < config.NumberMatching; i++ {
		allRegexps[fmt.Sprintf("matching%v", i)] = `[\w\.+-]+@[\w\.-]+\.[\w\.-]+`
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
	var globM1, globM2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&globM1)

	err := engine.CompileMany(allRegexps)
	if err != nil {
		return err
	}

	start := time.Now()
	count, err := engine.CountAll(data, config.IsDisplayOutput)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	runtime.ReadMemStats(&globM2)

	fmt.Printf("  [%v regexps] count=%v, mem=%.2fMB, time=%v \n",
		len(allRegexps), count, float64(globM2.TotalAlloc-globM1.TotalAlloc)/1000/1000, time.Duration(elapsed).Round(time.Microsecond))
	return nil
}
