package option

import "os"
import "fmt"
import "strings"

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

func extractArgumentsToOption(args []string, i int) ([]string, int) {

	var ret = make([]string, 0)

	for i < len(args)-1 {
		i++
		var arg = args[i]
		if slices.Contains(optionList, arg) {
			i--
			break
		}
		ret = append(ret, arg)
	}

	return ret, i

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
			ret.ExecArgs, i = extractArgumentsToOption(args, i)
			if len(ret.ExecArgs) == 0 {
				return ret, fmt.Errorf("`--args` with no argument")
			}

		case "--compile-args":
			ret.CompileArgs, i = extractArgumentsToOption(args, i)
			if len(ret.CompileArgs) == 0 {
				return ret, fmt.Errorf("`--compile-args` with no argument")
			}

		default:
			if strings.HasPrefix(arg, "-") {
				return ret, fmt.Errorf("unknown option: [ %v ]", arg)
			}
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
