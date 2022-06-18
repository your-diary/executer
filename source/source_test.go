package source

import "testing"
import "fmt"
import "strings"

func Test_misc(t *testing.T) {

	t.Run("with extension", func(t *testing.T) {

		var s = New("main.go")
		fmt.Println(s)

		if !((s.Original == "main.go") && (strings.HasPrefix(s.Path, "/") && strings.HasSuffix(s.Path, ".go")) && (strings.HasPrefix(s.PathWoExt, "/") && strings.HasSuffix(s.PathWoExt, "main")) && (s.Ext == "go") && (s.Base == "main.go") && strings.HasPrefix(s.Dir, "/") && (s.Name == "main")) {
			t.FailNow()
		}

	})

	t.Run("without extension", func(t *testing.T) {

		var s = New("main")
		fmt.Println(s)

		if !((s.Original == "main") && (strings.HasPrefix(s.Path, "/") && strings.HasSuffix(s.Path, "main")) && (s.PathWoExt == s.Path) && (s.Ext == "") && (s.Base == "main") && strings.HasPrefix(s.Dir, "/") && (s.Name == s.Base)) {
			t.FailNow()
		}

	})

}
