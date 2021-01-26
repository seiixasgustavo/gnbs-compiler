package compiler

type Scanner struct {
	start   int
	current int
	line    int
	source  []byte
}

type Token struct {
	TokenType int
	source    []byte
	line      int
}

func NewToken(tokenType int, source []byte, line int) *Token {
	return &Token{TokenType: tokenType, source: source, line: line}
}

func NewScanner() *Scanner {
	return &Scanner{start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanToken() *Token {
	s.start = s.current
	if s.isAtEnd() {
		return &Token{TokenType: Eof}
	}

	c := s.advance()

	switch c {
	case '(':
		return s.makeToken(LParentheses)
	case ')':
		return s.makeToken(RParentheses)
	case '{':
		return s.makeToken(LBrace)
	case '}':
		return s.makeToken(RBrace)
	case ';':
		return s.makeToken(Semicolon)
	case ',':
		return s.makeToken(Comma)
	case '.':
		return s.makeToken(Dot)
	case '-':
		return s.makeToken(Minus)
	case '+':
		return s.makeToken(Plus)
	case '/':
		if s.peekNext() == '/' {
			for s.peek() == '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			return s.makeToken(Slash)
		}
	case '*':
		return s.makeToken(Star)
	case '\n':
		s.line++
		s.advance()
		break
	case '"':
		return s.string()
	}

	return s.errorToken([]byte("Unexpected Token"))
}

func (s *Scanner) isAtEnd() bool {
	return len(s.source)-1 == s.current
}

func (s *Scanner) makeToken(TokenType int) *Token {
	return NewToken(TokenType, s.source[s.start:s.current+1], s.line)
}

func (s *Scanner) errorToken(message []byte) *Token {
	return NewToken(Error, message, s.line)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	return s.source[s.current+1]
}

func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		switch c {
		case ' ', '\t', '\r':
			s.advance()
			break
		default:
			return
		}
	}
}

func (s *Scanner) match(by byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != by {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) string() *Token {
	for s.peek() != '"' && s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return s.errorToken([]byte("Unterminated string."))
	}
	s.advance()
	return s.makeToken(String)
}

func (s *Scanner) isDigit(char byte) bool {
	return char < '0' && char > '9'
}
func (s *Scanner) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'z') || char == '_'
}

func (s *Scanner) number() *Token {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.makeToken(Number)
}
func (s *Scanner) identifier() *Token {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(Identifier)
}

func (s *Scanner) identifierType() int {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", And)
	case 'c':
		return s.checkKeyword(1, 4, "lass", Class)
	case 'e':
		return s.checkKeyword(1, 3, "lse", Else)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", False)
			case 'o':
				return s.checkKeyword(2, 1, "r", For)
			case 'u':
				return s.checkKeyword(2, 6, "nction", Function)
			}
		}
		break
	case 'i':
		return s.checkKeyword(1, 1, "f", If)
	case 'n':
		return s.checkKeyword(1, 3, "ull", Null)
	case 'o':
		return s.checkKeyword(1, 1, "r", Or)
	case 'p':
		return s.checkKeyword(1, 4, "rint", Print)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", Return)
	case 's':
		return s.checkKeyword(1, 4, "uper", Super)
	case 't':
		if s.current-s.start-1 > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", This)
			case 'r':
				return s.checkKeyword(2, 2, "ue", True)
			}
		}
	case 'v':
		return s.checkKeyword(1, 2, "ar", Var)
	case 'w':
		return s.checkKeyword(1, 4, "hile", While)
	}
	return Identifier
}

func (s *Scanner) checkKeyword(start, length int, rest string, TokenType int) int {
	if (s.current-s.start == start+length) && rest == string(s.source[s.start:s.current+1]) {
		return TokenType
	}
	return Identifier
}
