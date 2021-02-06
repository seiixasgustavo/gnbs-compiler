package scanner

import (
	"bytes"
	"fmt"
	"go/token"
	"path/filepath"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type ErrorHandler func(pos token.Position, msg string)

type Scanner struct {
	file *token.File
	dir  string
	src  []byte
	err  ErrorHandler

	ch         rune
	offset     int
	rdOffset   int
	lineOffset int

	insertSemi bool

	ErrorCount int
}

type Token struct {
	Position token.Pos
	Token    TokenType
	LitName  string
}

func NewScanner(file *token.File, src []byte, err ErrorHandler) *Scanner {
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}

	dir, _ := filepath.Split(file.Name())
	scanner := &Scanner{
		file:       file,
		dir:        dir,
		src:        src,
		err:        err,
		ch:         ' ',
		offset:     0,
		rdOffset:   0,
		lineOffset: 0,
		insertSemi: false,
		ErrorCount: 0,
	}
	scanner.next()
	return scanner
}

func (s *Scanner) Scan() (token *Token) {
	token = &Token{}

scanAgain:

	s.skipWhitespace()

	token.Position = s.file.Pos(s.offset)
	insertSemi := false

	switch ch := s.ch; {

	case isLetter(ch):
		token.LitName = s.scanIdentifier()
		if len(token.LitName) > 1 {
			token.Token = Lookup(token.LitName)

			switch token.Token {
			case Identifier, Return, Break:
				insertSemi = true
			}
		} else {
			insertSemi = true
			token.Token = Identifier
		}

	case isDecimal(ch) || ch == '.' && isDecimal(rune(s.peek())):
		insertSemi = true
		token.Token, token.LitName = s.scanNumber()

	default:
		s.next()

		switch ch {

		case -1:
			if s.insertSemi {
				s.insertSemi = false
				token.Token, token.LitName = Semicolon, "\n"
				return token
			}
			token.Token = Eof
		case '\n':
			s.insertSemi = false
			token.Token, token.LitName = Semicolon, "\n"
			return token
		case '"':
			insertSemi = true
			token.Token = String
			token.LitName = s.scanString()
		case '.':
			token.Token = Dot
		case ',':
			token.Token = Comma
		case ';':
			token.Token = Semicolon
			token.LitName = ";"
		case '(':
			token.Token = LParentheses
		case ')':
			insertSemi = true
			token.Token = RParentheses
		case '+':
			token.Token = Plus
		case '-':
			token.Token = Minus
		case '*':
			token.Token = Star
		case '/':
			if s.ch == '/' || s.ch == '*' {
				if s.insertSemi && s.findLineEnd() {
					s.ch = '/'
					s.offset = s.file.Offset(token.Position)
					s.rdOffset = s.offset + 1
					s.insertSemi = false
					return token
				}
				_ = s.scanComments()
				s.insertSemi = false
				goto scanAgain
			} else {
				token.Token = Slash
			}
		case '<':
			token.Token = s.switch2(Less, LessEqual)
		case '=':
			token.Token = s.switch2(Equal, EqualEqual)
		case '!':
			token.Token = s.switch2(Not, NotEqual)
		default:
			insertSemi = s.insertSemi
			token.Token = Illegal
			token.LitName = string(ch)
		}
	}
	s.insertSemi = insertSemi

	return
}

// Scan Functions

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanComments() string {
	offs := s.offset - 1
	next := -1
	numCR := 0

	// Comments using //
	if s.ch == '/' {
		s.next()
		for s.ch != '\n' && s.ch >= 0 {
			if s.ch == '\r' {
				numCR++
			}
			s.next()
		}
		next = s.offset
		if s.ch == '\n' {
			next++
		}
		//Comments using */
	} else {
		s.next()
		for s.ch >= 0 {
			ch := s.ch
			if ch == '\r' {
				numCR++
			}
			s.next()
			if ch == '*' && s.ch == '/' {
				s.next()
				next = s.offset
			} else {
				s.error(offs, "comment not terminated")
			}
		}
	}

	lit := s.src[offs:s.offset]
	if numCR > 0 && len(lit) >= 2 && lit[1] == '/' && lit[len(lit)-1] == '\r' {
		lit = lit[:len(lit)-1]
		numCR--
	}

	if next >= 0 && (lit[1] == '*' || offs == s.lineOffset) && bytes.HasPrefix(lit[2:], prefix) {
		s.updateLineInfo(next, offs, lit)
	}
	return string(lit)
}

func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset

	var n int
	var base, max uint32

	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanString() string {
	offs := s.offset - 1

	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}
		s.next()
		if ch == '"' {
			break
		}
		if ch == '\\' {
			s.scanEscape('"')
		}
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber() (TokenType, string) {
	offs := s.offset
	tk := &Token{Token: Illegal}

	base := 10
	prefix := rune(0)
	digsep := 0
	invalid := -1

	if s.ch != '.' {
		tk.Token = Integer
		if s.ch == '0' {
			s.next()
			switch lower(s.ch) {
			case 'x':
				s.next()
				base, prefix = 16, 'x'
			case 'o':
				s.next()
				base, prefix = 8, 'o'
			case 'b':
				s.next()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				digsep = 1
			}
		}
		digsep |= s.digits(base, &invalid)
	}

	if s.ch == '.' {
		tk.Token = Float
		if prefix == 'o' || prefix == 'b' {
			s.error(s.offset, "invalid radix point")
		}
		s.next()
		digsep |= s.digits(base, &invalid)
	}

	if digsep&1 == 0 {
		s.error(s.offset, litname(prefix)+" has no digits")
	}

	lit := string(s.src[offs:s.offset])
	if tk.Token == Integer && invalid >= 0 {
		s.errorf(invalid, "invalid digit %q in %s", lit[invalid-offs], litname(prefix))
	}
	if digsep&2 != 0 {
		if i := invalidSep(lit); i >= 0 {
			s.error(offs+i, "'_' must separate successive digits")
		}
	}
	return tk.Token, lit
}

// Aux Functions

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset

		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}

		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			r, w = utf8.DecodeLastRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal utf8 encoding")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.file.AddLine(s.offset)
		}
		s.ch = -1
	}
}

func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}

func (s *Scanner) error(offset int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offset)), msg)
	}
	s.ErrorCount++
}
func (s *Scanner) errorf(offs int, format string, args ...interface{}) {
	s.error(offs, fmt.Sprintf(format, args...))
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) updateLineInfo(next, offs int, text []byte) {
	if text[1] == '*' {
		text = text[:len(text)-2]
	}
	text = text[7:]
	offs += 7

	i, n, ok := trailingDigits(text)
	if i == 0 {
		return
	}
	if !ok {
		s.error(offs+i, "invalid line number: "+string(text[i:]))
		return
	}

	var line, col int
	i2, n2, ok2 := trailingDigits(text[:i-1])

	if ok2 {
		i, i2 = i2, i
		line, col = n2, n
		if col == 0 {
			s.error(offs+i2, "invalid column number: "+string(text[i2:]))
			return
		}
		text = text[:i2-1]
	} else {
		line = n
	}

	if line == 0 {
		s.error(offs+i, "invalid line number: "+string(text[i:]))
		return
	}

	filename := string(text[:i-1])
	if filename == "" && ok2 {
		filename = s.file.Position(s.file.Pos(offs)).Filename
	} else if filename != "" {
		filename = filepath.Clean(filename)
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(s.dir, filename)
		}
	}

	s.file.AddLineColumnInfo(next, filename, line, col)
}
func (s *Scanner) findLineEnd() bool {
	defer func(offs int) {
		s.ch = '/'
		s.offset = offs
		s.rdOffset = offs + 1
		s.next()
	}(s.offset - 1)

	for s.ch == '/' || s.ch == '*' {
		if s.ch == '/' {
			return true
		}

		s.next()
		for s.ch >= 0 {
			ch := s.ch
			if ch == '\n' {
				return true
			}
			s.next()
			if ch == '*' && s.ch == '/' {
				s.next()
				break
			}
		}
		s.skipWhitespace()
		if s.ch < 0 || s.ch == '\n' {
			return true
		}
		if s.ch != '/' {
			return false
		}
		s.next()
	}

	return false
}

var prefix = []byte("line ")

func (s *Scanner) digits(base int, invalid *int) (digsep int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			} else if s.ch >= max && *invalid < 0 {
				*invalid = s.offset
			}
			digsep |= ds
			s.next()
		}
	} else {
		for isHex(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			}
			digsep |= ds
			s.next()
		}
	}
	return
}
func (s *Scanner) switch2(tok0, tok1 TokenType) TokenType {
	if s.ch == '=' {
		s.next()
		return tok1
	}
	return tok0
}

// Utils/Standalone Functions

func trailingDigits(text []byte) (int, int, bool) {
	i := bytes.LastIndexByte(text, ':')
	if i < 0 {
		return 0, 0, false
	}
	n, err := strconv.ParseUint(string(text[i+1:]), 10, 0)
	return i + 1, int(n), err == nil
}
func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}
func lower(ch rune) rune     { return ('a' - 'A') | ch }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }
func litname(prefix rune) string {
	switch prefix {
	case 'x':
		return "hexadecimal literal"
	case 'o', '0':
		return "octal literal"
	case 'b':
		return "binary literal"
	}
	return "decimal literal"
}
func invalidSep(x string) int {
	x1 := ' ' // prefix char, we only care if it's 'x'
	d := '.'  // digit, one of '_', '0' (a digit), or '.' (anything else)
	i := 0

	// a prefix counts as a digit
	if len(x) >= 2 && x[0] == '0' {
		x1 = lower(rune(x[1]))
		if x1 == 'x' || x1 == 'o' || x1 == 'b' {
			d = '0'
			i = 2
		}
	}

	// mantissa and exponent
	for ; i < len(x); i++ {
		p := d // previous digit
		d = rune(x[i])
		switch {
		case d == '_':
			if p != '0' {
				return i
			}
		case isDecimal(d) || x1 == 'x' && isHex(d):
			d = '0'
		default:
			if p == '_' {
				return i - 1
			}
			d = '.'
		}
	}
	if d == '_' {
		return len(x) - 1
	}

	return -1
}
func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}
