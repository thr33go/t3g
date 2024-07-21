package t3g

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
	"github.com/kylelemons/godebug/diff"
)

type FileList struct {
	Filemap map[string]*os.File
}

func (f *FileList) FH(filename string) *os.File {
	if f.Filemap == nil {
		f.Filemap = make(map[string]*os.File)
	}
	fileHandle, ok := f.Filemap[filename]
	if !ok {
		fileHandle, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		f.Filemap[filename] = fileHandle
		return fileHandle
	}
	return fileHandle
}

func (f *FileList) Clear() {
	for k, v := range f.Filemap {
		err := v.Close()
		if err != nil {
			continue
		}
		info, _ := os.Stat(k)
		if info.Size() == 0 {
			os.Remove(k)
		}
	}
}

func File(filename string) *os.File {
	h := FLt3g.FH(filename)
	return h
}

func FileClear() {
	FLt3g.Clear()
}

func Set(slice interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return slice
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
		fmt.Println(expression, "compile error:", err)
		return ""
	}
	output, err := expr.Run(program, nil)
	if err != nil {
		fmt.Println(expression, "runtime error:", err)
		return ""
	}
	return output
}

func Bin2byte(binstream string) []byte {
	var s string
	for i := 0; i < len(binstream); i = i + 8 {
		i64, _ := strconv.ParseInt(binstream[i:i+8], 2, 8)
		// r = append(r, []byte(string(i64)))
		s = s + fmt.Sprint(i64)
	}
	return []byte(s)
}

func Bin2str(binstream string, digits int) string {
	var s string
	for i := 0; i < len(binstream); i = i + digits {
		i64, _ := strconv.ParseInt(binstream[i:i+digits], 2, 8)
		s = s + fmt.Sprint(i64) + " "
	}
	return s
}

func Hex2byte(hexstream string) []byte {
	hexstream = Re_space.ReplaceAllString(hexstream, "")
	src := []byte(hexstream)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	ErrLog(err)
	return dst
}

func XorHex(hexstream1, hexstream2 string) (string, error) {
	bytes1, err := hex.DecodeString(hexstream1)
	if err != nil {
		return "", err
	}
	bytes2, err := hex.DecodeString(hexstream2)
	if err != nil {
		return "", err
	}

	var long, short []byte
	if len(bytes1) > len(bytes2) {
		long, short = bytes1, bytes2
	} else {
		long, short = bytes2, bytes1
	}

	extendedShort := make([]byte, len(long))
	for i := range long {
		extendedShort[i] = short[i%len(short)]
	}

	result := make([]byte, len(long))
	for i := range long {
		result[i] = long[i] ^ extendedShort[i]
	}

	return hex.EncodeToString(result), nil
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type LogWriter struct{}

func (l LogWriter) Write(message []byte) (n int, err error) {
	c := color.New(color.FgHiMagenta, color.Bold)
	c.Printf("%s", message)
	return len(message), nil
}

func ErrLog(err error) {
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

func DiffAdded(a, b string) string {
	a_slice := strings.Split(a, "\n")
	b_slice := strings.Split(b, "\n")
	chunk := diff.DiffChunks(a_slice, b_slice)
	added := ""
	for _, c := range chunk {
		if c.Added == nil {
			continue
		}
		added += strings.Join(c.Added, "") + "\n"
	}
	return added
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

func GrayS(contents string) string {
	return "\033[0;90m" + contents + "\033[0m"
}

var Red = color.New(color.Bold, color.FgHiRed).PrintFunc()
var Redln = color.New(color.Bold, color.FgHiRed).PrintlnFunc()
var Redf = color.New(color.Bold, color.FgHiRed).PrintfFunc()
var RedS = color.New(color.Bold, color.FgHiRed).SprintFunc()
var Green = color.New(color.Bold, color.FgHiGreen).PrintFunc()
var Greenln = color.New(color.Bold, color.FgHiGreen).PrintlnFunc()
var Greenf = color.New(color.Bold, color.FgHiGreen).PrintfFunc()
var GreenS = color.New(color.Bold, color.FgHiGreen).SprintFunc()
var Yellow = color.New(color.Bold, color.FgHiYellow).PrintFunc()
var Yellowln = color.New(color.Bold, color.FgHiYellow).PrintlnFunc()
var Yellowf = color.New(color.Bold, color.FgHiYellow).PrintfFunc()
var YellowS = color.New(color.Bold, color.FgHiYellow).SprintFunc()
var Blue = color.New(color.Bold, color.FgHiBlue).PrintFunc()
var Blueln = color.New(color.Bold, color.FgHiBlue).PrintlnFunc()
var Bluef = color.New(color.Bold, color.FgHiBlue).PrintfFunc()
var BlueS = color.New(color.Bold, color.FgHiBlue).SprintFunc()
var Magenta = color.New(color.Bold, color.FgHiMagenta).PrintFunc()
var Magentaln = color.New(color.Bold, color.FgHiMagenta).PrintlnFunc()
var Magentaf = color.New(color.Bold, color.FgHiMagenta).PrintfFunc()
var MagentaS = color.New(color.Bold, color.FgHiMagenta).SprintFunc()
var Cyan = color.New(color.Bold, color.FgHiCyan).PrintFunc()
var Cyanln = color.New(color.Bold, color.FgHiCyan).PrintlnFunc()
var Cyanf = color.New(color.Bold, color.FgHiCyan).PrintfFunc()
var CyanS = color.New(color.Bold, color.FgHiCyan).SprintFunc()
var White = color.New(color.Bold, color.FgHiWhite).PrintFunc()
var Whiteln = color.New(color.Bold, color.FgHiWhite).PrintlnFunc()
var Whitef = color.New(color.Bold, color.FgHiWhite).PrintfFunc()
var WhiteS = color.New(color.Bold, color.FgHiWhite).SprintFunc()
var Re_space = *regexp.MustCompile(`\s+`)
var Uppercase []string = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var Lowercase []string = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var Numbers []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var Symbols []string = []string{"!", `"`, "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", ":", ";", "<", "=", ">", "?", "@", `[`, `\`, `]`, `^`, `_`, "`", `{`, `|`, `}`}
var Alphabet_slice []string
var Word_slice []string
var Any_slice []string
var FLt3g FileList
