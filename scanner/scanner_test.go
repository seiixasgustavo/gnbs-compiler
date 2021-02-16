package scanner

import (
	token2 "GNBS/token"
	"testing"
)

func TestNewScanner(t *testing.T) {
	src := []byte(
		`func main() {
			var x = 10
        }`)

	s := NewScanner(src, nil)

	for {
		tk := s.Scan()
		t.Logf("%s %q", tk.Token, tk.LitName)
		if tk.Token == token2.Eof {
			break
		}
	}
}
