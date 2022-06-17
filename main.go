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
	isDebugModeDefault         = 0
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

	var createExecOption = func(command string, isCompileMode bool) exec.Option {
		return exec.Option{
			IsCompileMode:              isCompileMode,
			Command:                    command,
			CompileOptions:             option.CompileArgs,
			Arguments:                  []string{option.Source},
			ExecOptions:                option.ExecArgs,
			ShouldMeasureTime:          option.ShouldMeasureTime,
			ExitStatusWhenCompileError: exitStatusWhenCompileError,
			IsDebugMode:                isDebugMode,
		}
	}

	var ext = path.Ext(option.Source)   //`.py`
	var base = path.Base(option.Source) //`main.py`

	if ext == ".py" {

		var o = createExecOption("python3", false)

		exec.Execute(o)

		os.Exit(0)

	} else if ext == ".rs" {

		if base == "main.rs" {

			if option.IsOnlyCompileMode {

				var o = createExecOption("cargo", true)
				o.CompileOptions = append([]string{"check", "--quiet"}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = nil

				exec.Execute(o)

			} else {

				var o = createExecOption("cargo", false)
				o.CompileOptions = append([]string{"run", "--quiet"}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = option.ExecArgs

				exec.Execute(o)

			}

		} else {

			var o = createExecOption("cargo", true)
			o.CompileOptions = append([]string{"check", "--quiet"}, option.CompileArgs...)
			o.Arguments = nil
			o.ExecOptions = nil

			exec.Execute(o)

		}

		os.Exit(0)

	} else {

		util.Eprintf("Unsupported file type: %v\n", ext)
		os.Exit(exitStatusWhenCompileError)

	}

}
