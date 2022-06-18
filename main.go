//usr/bin/env go run "$0" "$@"; exit

package main

import "fmt"
import "strings"
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

	var createExecOption = func(command string, isCompileMode bool) exec.Option {
		return exec.Option{
			IsCompileMode:              isCompileMode,
			Command:                    command,
			CompileOptions:             option.CompileArgs,
			Arguments:                  []string{option.Source.Path},
			ExecOptions:                option.ExecArgs,
			ShouldMeasureTime:          option.ShouldMeasureTime,
			ExitStatusWhenCompileError: exitStatusWhenCompileError,
			IsDebugMode:                isDebugMode,
		}
	}

	var s = option.Source

	switch s.Ext {

	case "py":
		{
			var o = createExecOption("python3", false)
			exec.Execute(o)
			os.Exit(0)
		}

	case "sh":
		{
			var o = createExecOption("bash", false)
			exec.Execute(o)
			os.Exit(0)
		}

	case "js":
		{
			var o = createExecOption("node", false)
			exec.Execute(o)
			os.Exit(0)
		}

	case "go":
		{
			if strings.HasSuffix(s.Base, "_test.go") { //test files
				var packageName = strings.Split(s.Base, "_")[0]
				var o = createExecOption("go", true)
				o.CompileOptions = append([]string{"test", "--count=1", "-v", fmt.Sprintf("./%v", packageName)}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = nil
				exec.Execute(o)
			} else { //normal files
				var output = s.PathWoExt + ".out"
				if !option.IsOnlyExecuteMode {
					var o = createExecOption("go", true)
					o.CompileOptions = append([]string{"build", "-o", output}, option.CompileArgs...)
					o.ExecOptions = nil
					exec.Execute(o)
				}
				if !option.IsOnlyCompileMode {
					var o = createExecOption(output, false)
					o.CompileOptions = nil
					o.Arguments = nil
					exec.Execute(o)
				}
			}
			os.Exit(0)
		}

	case "rs":
		{

			if s.Base == "main.rs" {

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

		}

	default:
		{
			util.Eprintf("Unsupported file type: %v\n", s.Ext)
			os.Exit(exitStatusWhenCompileError)
		}

	}

}
