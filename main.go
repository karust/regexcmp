package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/karust/regexcmp/engines"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

type Config struct {
	IsDisplayOutput    bool
	IsGroupRun         bool
	IsPrintGroupRegexp bool
	RegexFilePath      string
	ScanFilePath       string
	NumberMatching     int
	NumberNonMatching  int
	repeatScanTimes    int
	testCases          map[string]RegexpEngine
	execOrder          []string
}

var config Config

var RootCmd = &cobra.Command{
	Use:          "regexcmp",
	Short:        "Golang regex libraries comparison",
	Version:      version,
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.testCases = map[string]RegexpEngine{
			"rure":    &engines.Rure{},
			"pcre":    &engines.Pcre{},
			"default": &engines.Default{},
			"re2":     &engines.Re2{},
			"hyper":   &engines.Hyper{},
			"yara":    &engines.Yara{},
		}

		config.execOrder = []string{"rure", "pcre", "re2", "hyper", "yara", "default"}

		initRegexps()
	},

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			args = append(args, "1")
		}

		config.repeatScanTimes, err = strconv.Atoi(args[0])
		if err != nil {
			return errors.New("Provide numeric value: " + err.Error())
		}

		if len(args) >= 2 {
			testEngine := strings.ToLower(args[1])
			if _, ok := config.testCases[testEngine]; args[1] != "" && ok {
				config.execOrder = []string{testEngine}
			} else {
				fmt.Printf("No `%v` engine is available. Use one of %v\n", args[1], config.execOrder)
				return
			}
		}

		fmt.Println("Generate data...")
		var data = bytes.Repeat([]byte("mail@mail.co и nmber=+71112223334 URI:https://google.com 1.1.1.1 3FZbgi29cpjq2GjdwV8eyHuJJnkLtktZc5 "), config.repeatScanTimes)
		var locData *[]byte
		locData = &data

		fmt.Printf("Test data size: %.2fMB\n", float64(len(data))/1000/1000)

		if config.IsPrintGroupRegexp {
			groupRe := ""
			for name, regexpr := range allRegexps {
				regexpr = fmt.Sprintf("(?P<%v>%s)|", name, regexpr)
				groupRe += regexpr
			}
			fmt.Println("Group regexp: " + groupRe[:len(groupRe)-1])
		}

		for _, name := range config.execOrder {
			engine := config.testCases[name]
			fmt.Printf("Run %v:\n", strings.ToUpper(name))
			//fmt.Printf("Free memory: %dMB\n", memory.FreeMemory()/1024/1024)

			if config.IsGroupRun {
				err := runGroup(engine, locData)
				if err != nil {
					fmt.Printf("Error during test: %v", err)
				}
			} else {
				err := runSingle(engine, locData)
				if err != nil {
					fmt.Printf("Error during test: %v", err)
				}
			}
		}

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&config.IsDisplayOutput, "display", "d", false, "Display matched results")
	RootCmd.PersistentFlags().BoolVarP(&config.IsGroupRun, "group", "g", false, "Run grouped regexps (not individually)")
	RootCmd.PersistentFlags().BoolVarP(&config.IsPrintGroupRegexp, "print", "p", false, "Print constructed group regexp")
	RootCmd.PersistentFlags().StringVarP(&config.RegexFilePath, "rfile", "r", "./test/regexps.txt", "Host address to run server")
	RootCmd.PersistentFlags().StringVarP(&config.ScanFilePath, "sfile", "s", "./test/scanfile.txt", "Host address to run server")
	RootCmd.PersistentFlags().IntVarP(&config.NumberMatching, "matching", "m", 0, "Number of additional matching regexps")
	RootCmd.PersistentFlags().IntVarP(&config.NumberNonMatching, "nonmatching", "n", 0, "Number of additional non-matching regexps")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
