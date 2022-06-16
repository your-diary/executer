package util

import "fmt"
import "os"

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
