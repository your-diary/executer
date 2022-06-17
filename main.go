//usr/bin/env go run "$0" "$@"; exit

package main

// import "fmt"
import "os"

// import "strings"
import "path"

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

	var ext = path.Ext(option.Source)   //`.py`
	var base = path.Base(option.Source) //`main.py`

	if ext == ".py" {

		var execOption = exec.Option{
			IsCompileMode:              false,
			Command:                    "python3",
			CompileOptions:             option.CompileArgs,
			Arguments:                  []string{option.Source},
			ExecOptions:                option.ExecArgs,
			ExitStatusWhenCompileError: exitStatusWhenCompileError,
			IsDebugMode:                isDebugMode,
		}

		exec.Execute(execOption)

		os.Exit(0)

	} else if ext == ".rs" {

		if base == "main.rs" {

			if option.IsOnlyCompileMode {

				var execOption = exec.Option{
					IsCompileMode:              true,
					Command:                    "cargo",
					CompileOptions:             append([]string{"check", "--quiet"}, option.CompileArgs...),
					Arguments:                  nil,
					ExecOptions:                nil,
					ExitStatusWhenCompileError: exitStatusWhenCompileError,
					IsDebugMode:                isDebugMode,
				}

				exec.Execute(execOption)

			} else {

				var execOption = exec.Option{
					IsCompileMode:              false,
					Command:                    "cargo",
					CompileOptions:             append([]string{"run", "--quiet"}, option.CompileArgs...),
					Arguments:                  nil,
					ExecOptions:                option.ExecArgs,
					ExitStatusWhenCompileError: exitStatusWhenCompileError,
					IsDebugMode:                isDebugMode,
				}

				exec.Execute(execOption)

			}

		} else {

			var execOption = exec.Option{
				IsCompileMode:              true,
				Command:                    "cargo",
				CompileOptions:             append([]string{"check", "--quiet"}, option.CompileArgs...),
				Arguments:                  nil,
				ExecOptions:                nil,
				ExitStatusWhenCompileError: exitStatusWhenCompileError,
				IsDebugMode:                isDebugMode,
			}

			exec.Execute(execOption)

		}

		os.Exit(0)

	} else {

		util.Eprintf("Unsupported file type: %v\n", ext)
		os.Exit(exitStatusWhenCompileError)

	}

}
