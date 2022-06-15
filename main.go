//usr/bin/env go run "$0" "$@"; exit

package main

import "fmt"
import "errors"
import "os"
import "os/signal"
import "os/exec"

const (
	isDebugModeDefault         = 1
	exitStatusWhenCompileError = 100
)

var isDebugMode = false

func toStringPretty[T any](l []T) string {
    if (l == nil) {
        return "nil"
    }
    var ret = "[" + fmt.Sprintf("'%v'", l[0])
    for i := 1; i < len(l); i++ {
        ret += ", " + fmt.Sprintf("'%v'", l[i])
    }
    ret += "]"
    return ret
}

func eprintln[T any](t T) {
    fmt.Fprintln(os.Stderr, t)
}

func eprintf[T any](formatString string, t T) {
    fmt.Fprintf(os.Stderr, formatString, t)
}

func debugPrint[T any](t T) {
    if (isDebugMode) {
        eprintln(t)
    }
}

type ExecOption struct {
    isCompileMode bool
    command string
    compileOptions []string
    arguments []string
    execOptions []string
}

func execute(execOption ExecOption) {

    var exitStatusOnFailure = 1
    if (execOption.isCompileMode) {
        exitStatusOnFailure = exitStatusWhenCompileError
    }

    var args = make([]string, 0)
    args = append(args, execOption.compileOptions...)
    args = append(args, execOption.arguments...)
    args = append(args, execOption.execOptions...)

    if (isDebugMode) {
        var l = []string{execOption.command}
        l = append(l, args...)
        debugPrint(toStringPretty(l))
    }

//     var cmd = exec.Command("bash", "-c", "for i in $(seq 5); do echo ${i}; sleep 1; done")
//     var cmd = exec.Command("bash", "-c", "exit 150")
    var cmd = exec.Command(execOption.command, args...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Start(); err != nil {
        eprintf("Failed to execute the command: %v\n", err)
        os.Exit(exitStatusOnFailure)
    }

	var done = make(chan error)
    go func() {
        done <- cmd.Wait()
    }()

	var signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

    var err error
    select {
        case err = <- done:
        case <- signalChannel:
            debugPrint("\nSIGINT is caught.")
            if err := cmd.Process.Signal(os.Interrupt); err != nil {
                eprintln("Failed to send SIGINT.")
                os.Exit(exitStatusOnFailure)
            }
            os.Exit(exitStatusOnFailure)
    }

    if err != nil {
        var e *exec.ExitError
        if errors.As(err, &e) {
            os.Exit(e.ProcessState.ExitCode())
        } else {
            os.Exit(exitStatusOnFailure)
        }
    }

}

func main() {

	if (isDebugModeDefault != 0) || (os.Getenv("EXECUTER_DEBUG") == "1") {
		isDebugMode = true
	}

    var execOption = ExecOption {
        isCompileMode: false,
        command: "bash",
        compileOptions: []string{},
        arguments: []string{},
        execOptions: nil,
    }

    execute(execOption)

}
