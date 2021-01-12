package GNBS

import (
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   uint
	current uint
	line    uint
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1}
}

func (s *Scanner) advance() string {
	s.current++
	return strconv.Itoa(int(s.source[s.current-1]))
}

func (s *Scanner) peek() string {
	return strconv.Itoa(int(s.source[s.current]))
}

func (s *Scanner) isAtEnd() bool {
	return uint(len(tokens)) <= s.current
}

func (s *Scanner) addToken(token TokenType) {
	s.addTokenToList(token, nil)
}

func (s *Scanner) addTokenToList(token TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(token, text, literal, s.line))
}

func (s *Scanner) scanToken() error {
	char := s.advance()
	switch char {
	case "(":
		s.addToken(LParenteses)
		break
	case ")":
		s.addToken(RParenteses)
		break
	}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, *NewToken(Eof, "", nil, s.line))
	return s.tokens, nil
}
