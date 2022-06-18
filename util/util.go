package util

import "fmt"
import "os"
import "io"
import "errors"
import "strings"

func Eprintln[T any](t T) {
	fmt.Fprintln(os.Stderr, t)
}

func Eprintf[T any](formatString string, t T) {
	fmt.Fprintf(os.Stderr, formatString, t)
}

func DebugPrint[T any](t T, isDebugMode bool) {
	if isDebugMode {
		Eprintln(t)
	}
}

func ToStringPretty[T any](l []T) string {
	if l == nil {
		return "nil"
	}
	var ret = "[" + fmt.Sprintf("'%v'", l[0])
	for i := 1; i < len(l); i++ {
		ret += ", " + fmt.Sprintf("'%v'", l[i])
	}
	ret += "]"
	return ret
}

func IsFile(path string) bool {
	if info, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) && !info.IsDir() {
		return true
	}
	return false
}

func ReadFileUnchecked(file string) []string {
	var f, _ = os.Open(file)
	defer f.Close()
	var b, _ = io.ReadAll(f)
	return strings.Split(string(b), "\n")
}
