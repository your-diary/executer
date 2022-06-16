//usr/bin/env go run "$0" "$@"; exit

package main

// import "fmt"
import "os"

import "executer/option"
import "executer/util"
import "executer/exec"

const (
	isDebugModeDefault         = 1
	exitStatusWhenCompileError = 100
)

var isDebugMode = false

func main() {

	if (isDebugModeDefault != 0) || (os.Getenv("EXECUTER_DEBUG") == "1") {
		isDebugMode = true
	}

	var option, err = option.Parse(os.Args)
	if err != nil {
		util.Eprintf("Failed to parse command-line options: %v\n", err)
		os.Exit(exitStatusWhenCompileError)
	}
	util.DebugPrint(option, isDebugMode)

	var execOption = exec.Option{
		IsCompileMode:              false,
		Command:                    "python3",
		CompileOptions:             []string{"--version"},
		Arguments:                  []string{},
		ExecOptions:                nil,
		ExitStatusWhenCompileError: exitStatusWhenCompileError,
		IsDebugMode:                isDebugMode,
	}

	exec.Execute(execOption)

}
