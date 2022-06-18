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

	case "gp":
		{
			var o = createExecOption("gnuplot", false)
			o.CompileOptions = append([]string{"--persist"}, o.CompileOptions...)
			exec.Execute(o)
			os.Exit(0)
		}

	case "sql":
		{
			var o = createExecOption("sqlite3", false)
			o.CompileOptions = append(
				append(
					[]string{":memory:", "-init", "", "-batch"},
					o.CompileOptions...,
				),
				fmt.Sprintf(".read %v", s.Path),
			)
			o.Arguments = nil
			exec.Execute(o)
			os.Exit(0)
		}

	case "bats": //testing framework for Bash
		{
			var o = createExecOption("bats", false)
			o.CompileOptions = append([]string{"--print-output-on-failure", "--show-output-of-passing-tests"}, o.CompileOptions...)
			exec.Execute(o)
			os.Exit(0)
		}

	case "js":
		{
			var o = createExecOption("node", false)
			exec.Execute(o)
			os.Exit(0)
		}

	case "ts":
		{
			if !option.IsOnlyExecuteMode {
				var o = createExecOption("tsc", true)
				o.CompileOptions = append([]string{"--build"}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = nil
				exec.Execute(o)
			}
			if !option.IsOnlyCompileMode {
				var o = createExecOption("node", false)
				o.CompileOptions = nil
				o.Arguments = []string{fmt.Sprintf("%v/target/%v.js", s.Dir, s.Name)}
				exec.Execute(o)
			}
			os.Exit(0)
		}

	case "java":
		{
			if option.IsOnlyCompileMode {
				var o = createExecOption("gradle", true)
				o.CompileOptions = append([]string{"build"}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = nil
				exec.Execute(o)
			} else {
				var o = createExecOption("gradle", false)
				o.CompileOptions = append([]string{"run", "--quiet"}, option.CompileArgs...)
				o.Arguments = nil
				o.ExecOptions = func() []string { //To pass `a` and `b c`, we shall specify `['--args', '"a" "b c"']`.
					if len(o.ExecOptions) == 0 {
						return nil
					}
					for i := 0; i < len(o.ExecOptions); i++ {
						o.ExecOptions[i] = fmt.Sprintf(`"%v"`, o.ExecOptions[i])
					}
					return []string{"--args", strings.Join(o.ExecOptions, " ")}
				}()
				exec.Execute(o)
			}
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
