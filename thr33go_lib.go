package lib

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

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
	}
	os.Exit(1)
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
