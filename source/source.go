package source

import "strings"
import "path/filepath"

type Source struct {
	Original  string //original input
	Path      string //`/home/user/build/main.py`
	PathWoExt string //`/home/user/build/main`
	Ext       string //`py`
	Base      string //`main.py`
	Dir       string //`/home/user/build`
	Name      string //`main`
}

func (s Source) IsEmpty() bool {
	return s.Original == ""
}

func New(p string) Source {

	var ret Source

	var absPath, _ = filepath.Abs(p)

	ret.Original = p
	ret.Path = absPath
	ret.Ext = filepath.Ext(absPath)
	if ret.Ext != "" {
		ret.Ext = ret.Ext[1:]
	}
	ret.Base = filepath.Base(absPath)
	ret.Dir = filepath.Dir(absPath)

	if ret.Ext == "" {
		ret.Name = ret.Base
		ret.PathWoExt = ret.Path
	} else {

		ret.Name = func() string {
			var l = strings.Split(ret.Base, ".")
			return strings.Join(l[:len(l)-1], ".")
		}()

		ret.PathWoExt = func() string {
			var l = strings.Split(ret.Path, ".")
			return strings.Join(l[:len(l)-1], ".")
		}()

	}

	return ret

}
