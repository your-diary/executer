package option

import "testing"
import "fmt"
import "strings"

import "golang.org/x/exp/slices"

func Test_misc(t *testing.T) {

	t.Run("No option.", func(t *testing.T) {

		//mock
		var exitStatus int = -1
		exit = func(i int) {
			exitStatus = i
		}

		var args = []string{"$0"}

		var _, _ = Parse(args)

		if exitStatus != 0 {
			t.FailNow()
		}

	})

	t.Run("All options.", func(t *testing.T) {

		var args = []string{
			"$0",
			"main.go",
			"--time",
			"--only-compile",
			"--only-execute",
			"--args",
			"a",
			"b",
			"--compile-args",
			"c",
			"d",
		}

		var ret, err = Parse(args)

		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if !((ret.Source.Original == "main.go") && slices.Equal(ret.ExecArgs, []string{"a", "b"}) && slices.Equal(ret.CompileArgs, []string{"c", "d"}) && (ret.IsOnlyCompileMode == true) && (ret.IsOnlyExecuteMode == true) && (ret.ShouldMeasureTime == true)) {
			t.FailNow()
		}

	})

}

func Test_help(t *testing.T) {

	for _, option := range []string{"-h", "--help"} {

		t.Run(fmt.Sprintf("`os.Exit(0)` is called when `%s`", option), func(t *testing.T) {

			//mock
			var exitStatus int = -1
			exit = func(i int) {
				exitStatus = i
			}

			var args = []string{"$0", "--time", option}
			Parse(args)

			if exitStatus != 0 {
				t.FailNow()
			}

		})

	}

}

func Test_args(t *testing.T) {

	t.Run("`--args` with arguments", func(t *testing.T) {

		var args = []string{"$0", "--args", "a", "b", "--time", "main.go"}

		var ret, err = Parse(args)

		if err != nil {
			t.FailNow()
		}

		if !slices.Equal(ret.ExecArgs, []string{"a", "b"}) {
			fmt.Println(ret)
			t.FailNow()
		}

	})

}

func Test_compileArgs(t *testing.T) {

	t.Run("`--compile-args` with arguments", func(t *testing.T) {

		var args = []string{"$0", "--compile-args", "a", "b", "--time", "main.go"}

		var ret, err = Parse(args)

		if err != nil {
			t.FailNow()
		}

		if !slices.Equal(ret.CompileArgs, []string{"a", "b"}) {
			fmt.Println(ret)
			t.FailNow()
		}

	})

}

func Test_source(t *testing.T) {

	t.Run("More than one sources are specified.", func(t *testing.T) {

		var args = []string{"$0", "main.go", "main.py"}

		var _, err = Parse(args)
		fmt.Println(err)

		if !strings.HasPrefix(err.Error(), "more than one sources specified") {
			t.Fatal(err)
		}

	})

	t.Run("Unknown option.", func(t *testing.T) {

		var args = []string{"$0", "main.go", "--main.py"}

		var _, err = Parse(args)
		fmt.Println(err)

		if !strings.HasPrefix(err.Error(), "unknown option") {
			t.Fatal(err)
		}

	})

	t.Run("No source specified.", func(t *testing.T) {

		var args = []string{"$0"}

		var _, err = Parse(args)
		fmt.Println(err)

		if !strings.HasPrefix(err.Error(), "no source specified") {
			t.FailNow()
		}

	})

}
