package scanner

import (
	"go/token"
	"testing"
)

func TestNewScanner(t *testing.T) {
	src := []byte("const value = 10\nvar x = 10.0")

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s := NewScanner(file, src, nil)

	for {
		tk := s.Scan()
		t.Logf("%s %s %q", fset.Position(tk.Position), tk.Token, tk.LitName)
		if tk.Token == Eof {
			break
		}
	}
}
