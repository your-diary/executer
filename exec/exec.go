package exec

// import "fmt"
import "errors"
import "os"
import "os/signal"
import "os/exec"
import "time"

import "executer/util"

type Option struct {
	IsCompileMode              bool
	Command                    string
	CompileOptions             []string
	Arguments                  []string
	ExecOptions                []string
	ShouldMeasureTime          bool
	ExitStatusWhenCompileError int
	IsDebugMode                bool
}

func Execute(o Option) {

	var start = time.Now()
	var exit = func(exitStatus int) {
		var elapsedSeconds float64 = float64(time.Now().Sub(start).Milliseconds()) / 1000
		if !o.IsCompileMode && (o.ShouldMeasureTime || o.IsDebugMode) {
			util.Eprintf("\nElapsed: %.2f(s)\n", elapsedSeconds)
		}
		if exitStatus != 0 {
			os.Exit(exitStatus)
		}
	}

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

	var cmd = exec.Command(o.Command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		util.Eprintf("Failed to execute the command: %v\n", err)
		exit(exitStatusOnFailure)
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
			exit(exitStatusOnFailure)
		}
		<-done
		exit(exitStatusOnFailure)
	}

	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			exit(e.ProcessState.ExitCode())
		} else {
			exit(exitStatusOnFailure)
		}
	}

	exit(0)

}
