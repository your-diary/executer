package option

import "os"
import "fmt"

import "golang.org/x/exp/slices"

type Options struct {
	Source            string
	CompileArgs       []string
	ExecArgs          []string
	IsOnlyCompileMode bool
	IsOnlyExecuteMode bool
	ShouldMeasureTime bool
}

var optionList = []string{
	"--compile-args",
	"--args",
	"--only-compile",
	"--only-execute",
	"--time",
	"-h",
	"--help",
}

func printUsage() {
	fmt.Println(`Usage
  executer <file> [<option(s)>]

Options
  --compile-args <arg(s)>    #Passes <arg(s)> when compilation.
  --args <arg(s)>            #Passes <arg(s)> when execution.
  --only-compile             #Just compiles and skips execution.
  --only-execute             #Just executes and skips compilation.
  --time                     #Measures the execution time.
  -h/--help                  #Shows this help.`)
}

var exit func(int) = os.Exit //for mock

func Parse(args []string) (Options, error) {

	var ret = Options{}

	var i = 0
	for i < len(args)-1 {
		i++
		var arg = args[i]
		switch arg {
		case "-h", "--help":
			printUsage()
			exit(0)
		case "--only-compile":
			ret.IsOnlyCompileMode = true
		case "--only-execute":
			ret.IsOnlyExecuteMode = true
		case "--time":
			ret.ShouldMeasureTime = true
		case "--args":
			if i == len(args)-1 {
				return ret, fmt.Errorf("`--args` with no argument")
			}
			ret.ExecArgs = []string{}
			for i < len(args)-1 {
				i++
				var arg = args[i]
				if slices.Contains(optionList, arg) {
					i--
					break
				}
				ret.ExecArgs = append(ret.ExecArgs, arg)
			}
		case "--compile-args":
			if i == len(args)-1 {
				return ret, fmt.Errorf("`--compile-args` with no argument")
			}
			ret.CompileArgs = []string{}
			for i < len(args)-1 {
				i++
				var arg = args[i]
				if slices.Contains(optionList, arg) {
					i--
					break
				}
				ret.CompileArgs = append(ret.CompileArgs, arg)
			}
		default:
			if ret.Source != "" {
				return ret, fmt.Errorf("more than one sources specified: [ %v, %v ]", ret.Source, arg)
			}
			ret.Source = arg
		}
	}

	if ret.Source == "" {
		return ret, fmt.Errorf("no source specified")
	}

	return ret, nil

}
