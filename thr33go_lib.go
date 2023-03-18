package lib

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/fatih/color"
)

func Set(slice interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}

	seen := make(map[interface{}]bool)
	result := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		v := sliceValue.Index(i).Interface()
		if _, ok := seen[v]; !ok {
			seen[v] = true
			result = reflect.Append(result, reflect.ValueOf(v))
		}
	}
	return result.Interface()
}

func Expr(expression string) any {
	program, err := expr.Compile(expression)
	if err != nil {
		fmt.Println("compile error:", err)
		return ""
	}
	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Println("runtime error:", err)
		return ""
	}
	return output
}

func Bin2byte(binstream string) []byte {
	var s string
	for i := 0; i < len(binstream); i = i + 8 {
		i64, _ := strconv.ParseInt(binstream[i:i+8], 2, 8)
		// r = append(r, []byte(string(i64)))
		s = s + string(i64)
	}
	return []byte(s)
}

func Hex2byte(hexstream string) []byte {
	hexstream = Re_space.ReplaceAllString(hexstream, "")
	src := []byte(hexstream)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	Errlog(err)
	return dst
}

type LogWriter struct{}

func (l LogWriter) Write(message []byte) (n int, err error) {
	c := color.New(color.FgHiMagenta, color.Bold)
	c.Printf("%s", message)
	return len(message), nil
}

func Errlog(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ErrExit(err error) {
	if err != nil {
		log.Println(err)
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

func ScanTimeout(duration string) (string, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	nameCh := make(chan string)
	errCh := make(chan error)

	go func() {
		var name string
		if _, err := fmt.Scan(&name); err != nil {
			errCh <- err
			return
		}
		nameCh <- name
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case name := <-nameCh:
		return name, nil
	}
}

func init() {
	lw := LogWriter{}
	// log.SetFlags(log.Ltime | log.Llongfile)
	log.SetFlags(log.Ltime)
	log.SetOutput(lw)
	Alphabet_slice = append(Lowercase, Uppercase...)
	Word_slice = append(Alphabet_slice, Numbers...)
	Any_slice = append(Word_slice, Symbols...)
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
var Re_space = *regexp.MustCompile(`\s+`)
var Uppercase []string = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var Lowercase []string = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var Numbers []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var Symbols []string = []string{"!", `"`, "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", ":", ";", "<", "=", ">", "?", "@", `[`, `\`, `]`, `^`, `_`, "`", `{`, `|`, `}`}
var Alphabet_slice []string
var Word_slice []string
var Any_slice []string
