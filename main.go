package main

import (
	"executer/exec"
	"executer/option"
	"executer/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mattn/go-isatty"
)

const (
	isDebugModeDefault         = 0
	exitStatusWhenCompileError = 180
)

var isDebugMode = false

func main() {

	if isatty.IsTerminal(os.Stderr.Fd()) {
		util.Eprintln("\u001B[1;034m==================================================================\u001B[0m")
	}

	if (isDebugModeDefault != 0) || (os.Getenv("EXECUTER_DEBUG") == "1") {
		isDebugMode = true
	}

	var option, err = option.Parse(os.Args)
	if err != nil {
		util.Eprintf("Failed to parse command-line options: %v\n", err)
		os.Exit(exitStatusWhenCompileError)
	}
	util.DebugPrint(option, isDebugMode)

	if option.IsOnlyCompileMode {
		util.Eprintln("\u001B[094mOnly-compile mode.\u001B[0m")
	}

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

	//yrun.sh
	//We respect `yrun.sh` iff the following three conditions are met.
	//1. It exists.
	//2. It isn't empty.
	//3. It doesn't consist only of comments.
	{
		var file = "./yrun.sh"
		if util.IsFile(file) {

			//checks if `./yrun.sh` is empty
			var lines = util.ReadFileUnchecked(file)
			var isYrunShEmpty = true
			for _, line := range lines {
				var l = strings.TrimSpace(line)
				if !((l == "") || strings.HasPrefix(l, "#")) {
					isYrunShEmpty = false
					break
				}
			}

			if !isYrunShEmpty {
				var o = createExecOption("bash", false)
				o.CompileOptions = nil
				o.Arguments = append([]string{file}, o.Arguments...)
				exec.Execute(o)
				os.Exit(0)
			}

		}
	}

	var s = option.Source

	switch s.Ext {

	case "py":
		{
			var o exec.Option
			if runtime.GOOS == "darwin" {
				o = createExecOption("python3.11", false)
			} else {
				o = createExecOption("python3", false)
			}
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

	case "awk":
		{
			var prog = strings.Join(util.ReadFileUnchecked(s.Path), "\n")
			//We require there is a `BEGIN` block to avoid stdin's begin read.
			if !strings.Contains(prog, "BEGIN {") {
				util.Eprintln("The input doesn't include `BEGIN { ... }` block.")
				os.Exit(exitStatusWhenCompileError)
			}
			var o = createExecOption("awk", false)
			o.Arguments = []string{prog}
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
			if strings.HasSuffix(s.Original, "test.ts") {
				var o = createExecOption("npm", false)
				o.CompileOptions = append([]string{"test"}, o.CompileOptions...)
				o.Arguments = nil
				exec.Execute(o)
			} else {
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
			}
			os.Exit(0)
		}

	case "c", "cpp":
		{
			var output = s.PathWoExt + ".out"
			if !option.IsOnlyExecuteMode {
				var o = func() exec.Option {

					if runtime.GOOS == "darwin" {
						if s.Ext == "c" {
							return createExecOption("gcc-12", true)
						}
						return createExecOption("g++-12", true)
					}

					if s.Ext == "c" {
						return createExecOption("gcc", true)
					}
					return createExecOption("g++", true)

				}()
				o.CompileOptions = append([]string{"-fdiagnostics-color=always", "-Wno-unused-result", "-Wfatal-errors", "-o", output}, option.CompileArgs...)
				if s.Ext == "c" {
					o.CompileOptions = append(o.CompileOptions, "-l", "m")
				}
				o.ExecOptions = nil
				exec.Execute(o)
			}
			if !option.IsOnlyCompileMode {
				var o = createExecOption(output, false)
				o.CompileOptions = nil
				o.Arguments = nil
				exec.Execute(o)
			}
			os.Exit(0)
		}

	case "java":
		{
			if util.IsFile("./settings.gradle") { //project
				if !option.IsOnlyExecuteMode {
					var o = createExecOption("gradle", true)
					o.CompileOptions = append([]string{"build", "--quiet", "--console", "plain"}, option.CompileArgs...)
					o.Arguments = nil
					o.ExecOptions = nil
					exec.Execute(o)
				}
				if !option.IsOnlyCompileMode {
					var fqcn = func() string {
						var packageName = regexp.MustCompile(`package ([^;]+);`).FindStringSubmatch(
							strings.Join(util.ReadFileUnchecked(s.Path), "\n"),
						)[1]
						return fmt.Sprintf("%v.%v", packageName, s.Name)
					}()
					var o = createExecOption("java", false)
					o.CompileOptions = []string{
						"-enableassertions",
						"--class-path",
						"./app/build/classes/java/main:./app/build/classes/java/test/",
					}
					o.Arguments = []string{fqcn}
					exec.Execute(o)
				}
			} else { //non-project (unit file)
				if !option.IsOnlyExecuteMode {
					var o = createExecOption("javac", true)
					o.ExecOptions = nil
					exec.Execute(o)
				}
				if !option.IsOnlyCompileMode {
					var o = createExecOption("java", false)
					o.CompileOptions = []string{"-enableassertions"}
					o.Arguments = []string{s.Name}
					exec.Execute(o)
				}
			}
			os.Exit(0)
		}

	case "hs":
		{
			var cabalFiles, _ = filepath.Glob("*.cabal")
			if cabalFiles != nil { //project
				if strings.Contains(s.Path, "/test/") { //test files
					var o = createExecOption("cabal", true)
					o.CompileOptions = append([]string{"test", "-v0", "--test-show-details=streaming", "--test-option=--color", "--ghc-options=-Wall"}, option.CompileArgs...)
					o.Arguments = nil
					o.ExecOptions = nil
					exec.Execute(o)
				} else {
					var packageName = regexp.MustCompile(`\.cabal$`).ReplaceAllString(cabalFiles[0], "")
					if !option.IsOnlyExecuteMode {
						var o = createExecOption("cabal", true)
						o.CompileOptions = append([]string{"build", "-v0", "--ghc-options=-Wall"}, option.CompileArgs...)
						o.Arguments = nil
						o.ExecOptions = nil
						exec.Execute(o)
					}
					if !option.IsOnlyCompileMode && (s.Base == "Main.hs") {
						var o = createExecOption("cabal", false)
						o.CompileOptions = []string{"exec", packageName}
						o.Arguments = nil
						exec.Execute(o)
					}
				}
			} else {
				var output = s.PathWoExt + ".out"
				if !option.IsOnlyExecuteMode {
					var o = createExecOption("ghc", true)
					o.CompileOptions = append([]string{"-v0", "-Wall", "-Wno-type-defaults", "-o", output}, option.CompileArgs...)
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
				if util.IsFile("./go.mod") { //project
					if !option.IsOnlyExecuteMode {
						var o = createExecOption("go", true)
						o.CompileOptions = append([]string{"build"}, option.CompileArgs...)
						o.Arguments = nil
						o.ExecOptions = nil
						exec.Execute(o)
					}
					if !option.IsOnlyCompileMode && (s.Base == "main.go") {
						var output = func() string {
							var moduleName = regexp.MustCompile(`^module (.*)$`).FindStringSubmatch(
								util.ReadFileUnchecked("./go.mod")[0],
							)[1]
							return fmt.Sprintf("./%v", moduleName)
						}()
						var o = createExecOption(output, false)
						o.CompileOptions = nil
						o.Arguments = nil
						exec.Execute(o)
					}
				} else { //non-project (unit file)
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
			}
			os.Exit(0)
		}

	case "rs":
		{
			if !util.IsFile("./Cargo.toml") {
				util.Eprintln("`Cargo.toml` not found. `cd ..` in Vim may help.")
				os.Exit(exitStatusWhenCompileError)
			}
			if s.Base == "main.rs" {
				if strings.Contains(s.Path, "/atcoder/") {
					os.Setenv("RUST_BACKTRACE", "0")
					var o = createExecOption("cargo", true)
					o.CompileOptions = append([]string{"test"}, option.CompileArgs...)
					o.Arguments = nil
					o.ExecOptions = nil
					exec.Execute(o)
				} else {
					if !option.IsOnlyExecuteMode {
						var o = createExecOption("cargo", true)
						o.CompileOptions = append([]string{"build"}, option.CompileArgs...)
						o.Arguments = nil
						o.ExecOptions = nil
						exec.Execute(o)
					}
					if !option.IsOnlyCompileMode {
						var output = func() string {
							var packageName = regexp.MustCompile(`^name = "(.*)"$`).FindStringSubmatch(
								util.ReadFileUnchecked("./Cargo.toml")[1],
							)[1]
							return fmt.Sprintf("./target/debug/%v", packageName)
						}()
						var o = createExecOption(output, false)
						o.CompileOptions = nil
						o.Arguments = nil
						exec.Execute(o)
					}
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
