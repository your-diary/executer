package option

import "testing"
import "fmt"

import "golang.org/x/exp/slices"

func Test_misc(t *testing.T) {

	t.Run("No option.", func(t *testing.T) {

		var args = []string{"$0"}

		var ret, _ = Parse(args)

		if !((ret.Source == "") && (ret.CompileArgs == nil) && (ret.ExecArgs == nil) && (ret.IsOnlyCompileMode == false) && (ret.IsOnlyExecuteMode == false) && (ret.ShouldMeasureTime == false)) {
			t.FailNow()
		}

	})

	t.Run("All options.", func(t *testing.T) {

		var args = []string{"$0", "main.go", "--time", "--only-compile", "--only-execute", "--args", "a", "b", "--compile-args", "c", "d"}

		var ret, err = Parse(args)

		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if !((ret.Source == "main.go") && slices.Equal(ret.ExecArgs, []string{"a", "b"}) && slices.Equal(ret.CompileArgs, []string{"c", "d"}) && (ret.IsOnlyCompileMode == true) && (ret.IsOnlyExecuteMode == true) && (ret.ShouldMeasureTime == true)) {
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

	t.Run("`err != nil` when `--args` has no argument", func(t *testing.T) {

		var args = []string{"$0", "--time", "--args"}

		var _, err = Parse(args)
		fmt.Println(err)

		if err == nil {
			t.FailNow()
		}

	})

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

	t.Run("`err != nil` when `--compile-args` has no argument", func(t *testing.T) {

		var args = []string{"$0", "--time", "--compile-args"}

		var _, err = Parse(args)
		fmt.Println(err)

		if err == nil {
			t.FailNow()
		}

	})

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

		if err == nil {
			t.FailNow()
		}

	})

	t.Run("No source specified.", func(t *testing.T) {

		var args = []string{"$0"}

		var _, err = Parse(args)
		fmt.Println(err)

		if err == nil {
			t.FailNow()
		}

	})

}
