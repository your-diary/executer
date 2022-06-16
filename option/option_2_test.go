package option

import "testing"

import "golang.org/x/exp/slices"

func Test_extractArgumentsToOption(t *testing.T) {

	t.Run("1", func(t *testing.T) {

		var args = []string{"--args", "a", "b", "--time"}
		var i = 0

		ret, i := extractArgumentsToOption(args, i)

		if !slices.Equal(ret, []string{"a", "b"}) {
			t.Fatal(ret)
		}

		if i != 2 {
			t.Fatal(i)
		}

	})

	for _, args := range [][]string{[]string{"--args"}, []string{"--args", "--time"}} {

		t.Run("2", func(t *testing.T) {

			var i = 0

			ret, i := extractArgumentsToOption(args, i)

			if len(ret) != 0 {
				t.Fatal(ret)
			}

			if i != 0 {
				t.Fatal(i)
			}

		})

	}

}
