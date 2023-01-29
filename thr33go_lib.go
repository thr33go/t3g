package thr33go

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

func ErrExit(err error) {
	if err != nil {
		Red(err)
		os.Exit(1)
	}
}

func LenCheckRtn(s []string, idx int) string {
	if len(s) > idx {
		return s[idx]
	}
	return ""
}

func StringInSlice(str string, list []string) bool {
	for _, b := range list {
		if b == str {
			return true
		}
	}
	return false
}

func ReadFileLines(filename string) []string {
	file, _ := os.ReadFile(filename)
	file_str := string(file)
	return strings.Split(file_str, "\n")
}

var Red = color.New(color.Bold, color.FgHiRed).PrintFunc()
var Redln = color.New(color.Bold, color.FgHiRed).PrintlnFunc()
var Green = color.New(color.Bold, color.FgHiGreen).PrintFunc()
var Greenln = color.New(color.Bold, color.FgHiGreen).PrintlnFunc()
var Yellow = color.New(color.Bold, color.FgHiYellow).PrintFunc()
var Yellowln = color.New(color.Bold, color.FgHiYellow).PrintlnFunc()
var Blue = color.New(color.Bold, color.FgHiBlue).PrintFunc()
var Blueln = color.New(color.Bold, color.FgHiBlue).PrintlnFunc()
var Magenta = color.New(color.Bold, color.FgHiMagenta).PrintFunc()
var Magentaln = color.New(color.Bold, color.FgHiMagenta).PrintlnFunc()
var Cyan = color.New(color.Bold, color.FgHiCyan).PrintFunc()
var Cyanln = color.New(color.Bold, color.FgHiCyan).PrintlnFunc()
var White = color.New(color.Bold, color.FgHiWhite).PrintFunc()
var Whiteln = color.New(color.Bold, color.FgHiWhite).PrintlnFunc()
var Print = color.New(color.FgWhite).PrintFunc()
var Println = color.New(color.FgWhite).PrintlnFunc()
