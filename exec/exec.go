package exec

// import "fmt"
import "errors"
import "os"
import "os/signal"
import "os/exec"

import "executer/util"

type Option struct {
	IsCompileMode              bool
	Command                    string
	CompileOptions             []string
	Arguments                  []string
	ExecOptions                []string
	ExitStatusWhenCompileError int
	IsDebugMode                bool
}

func Execute(o Option) {

	var exitStatusOnFailure = 1
	if o.IsCompileMode {
		exitStatusOnFailure = o.ExitStatusWhenCompileError
	}

	var args = make([]string, 0)
	args = append(args, o.CompileOptions...)
	args = append(args, o.Arguments...)
	args = append(args, o.ExecOptions...)

	if o.IsDebugMode {
		var l = []string{o.Command}
		l = append(l, args...)
		util.DebugPrint(util.ToStringPretty(l), true)
	}

	//     var cmd = exec.Command("bash", "-c", "for i in $(seq 5); do echo ${i}; sleep 1; done")
	//     var cmd = exec.Command("bash", "-c", "exit 150")
	var cmd = exec.Command(o.Command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		util.Eprintf("Failed to execute the command: %v\n", err)
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
	case err = <-done:
	case <-signalChannel:
		util.DebugPrint("\nSIGINT is caught.", o.IsDebugMode)
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			util.Eprintln("Failed to send SIGINT.")
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
