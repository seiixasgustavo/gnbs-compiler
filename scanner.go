package GNBS

import (
	"GNBS/token"
	"strconv"
)

var keywords map[string]token.Type

func init() {
	keywords = make(map[string]token.Type)

	keywords[token.Lexeme[token.And]] = token.And
	keywords[token.Lexeme[token.Or]] = token.Or
	keywords[token.Lexeme[token.If]] = token.If
	keywords[token.Lexeme[token.Else]] = token.Else
	keywords[token.Lexeme[token.Var]] = token.Var
	keywords[token.Lexeme[token.Func]] = token.Func
	keywords[token.Lexeme[token.Struct]] = token.Struct
	keywords[token.Lexeme[token.Return]] = token.Return
	keywords[token.Lexeme[token.True]] = token.True
	keywords[token.Lexeme[token.False]] = token.False
}

type Scanner struct {
	source string
	Lexeme []token.Token

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

func (s *Scanner) peekNext() string {
	return strconv.Itoa(int(s.source[s.current+1]))
}

func (s *Scanner) match(char string) bool {
	if s.isAtEnd() {
		return false
	}
	if strconv.Itoa(int(s.source[s.current])) != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		Compiler.Error(s.line, "Unterminated string")
	}

	s.advance()
	value := s.source[s.start+1 : s.current]
	s.addTokenToList(token.String, value)
}

func (s *Scanner) number() {
	var isFloat = false
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == "." && s.isDigit(s.peekNext()) {
		isFloat = true

		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	if isFloat {
		value, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
		s.addTokenToList(token.Float, value)
	} else {
		value, _ := strconv.ParseInt(s.source[s.start:s.current], 10, 64)
		s.addTokenToList(token.Integer, value)
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	s.addToken(token.Identifier)
}

func (s *Scanner) isAlpha(char string) bool {
	return (char >= "a" && char <= "z") || (char >= "A" && char <= "Z") || char == "_"
}

func (s *Scanner) isDigit(char string) bool {
	return char >= "0" && char <= "9"
}

func (s *Scanner) isAlphaNumeric(char string) bool {
	return s.isDigit(char) || s.isAlpha(char)
}

func (s *Scanner) isAtEnd() bool {
	return uint(len(token.Lexeme)) <= s.current
}

func (s *Scanner) addToken(token token.Type) {
	s.addTokenToList(token, nil)
}

func (s *Scanner) addTokenToList(tk token.Type, literal interface{}) {
	text := s.source[s.start:s.current]
	s.Lexeme = append(s.Lexeme, token.NewToken(tk, text, literal, s.line))
}

func (s *Scanner) ScanLexeme() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.Lexeme = append(s.Lexeme, token.NewToken(token.Eof, "", nil, s.line))
	return s.Lexeme, nil
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case token.Lexeme[token.LParentheses]:
		s.addToken(token.LParentheses)
		break
	case token.Lexeme[token.RParentheses]:
		s.addToken(token.RParentheses)
		break
	case token.Lexeme[token.LBrace]:
		s.addToken(token.LBrace)
		break
	case token.Lexeme[token.RBrace]:
		s.addToken(token.RBrace)
		break
	case token.Lexeme[token.Comma]:
		s.addToken(token.Comma)
		break
	case token.Lexeme[token.Dot]:
		s.addToken(token.Dot)
		break
	case token.Lexeme[token.Minus]:
		s.addToken(token.Minus)
		break
	case token.Lexeme[token.Plus]:
		s.addToken(token.Plus)
		break
	case token.Lexeme[token.Semicolon]:
		s.addToken(token.Semicolon)
		break
	case token.Lexeme[token.Star]:
		s.addToken(token.Star)
		break

	case token.Lexeme[token.Not]:
		if s.match(token.Lexeme[token.Equal]) {
			s.addToken(token.NotEqual)
		} else {
			s.addToken(token.Not)
		}
		break
	case token.Lexeme[token.Equal]:
		if s.match(token.Lexeme[token.Equal]) {
			s.addToken(token.EqualEqual)
		} else {
			s.addToken(token.Equal)
		}
		break
	case token.Lexeme[token.Greater]:
		if s.match(token.Lexeme[token.Equal]) {
			s.addToken(token.GreaterEqual)
		} else {
			s.addToken(token.Greater)
		}
		break
	case token.Lexeme[token.Less]:
		if s.match(token.Lexeme[token.LessEqual]) {
			s.addToken(token.LessEqual)
		} else {
			s.addToken(token.Less)
		}
		break
	case token.Lexeme[token.Slash]:
		if s.match(token.Lexeme[token.Slash]) {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.Slash)
		}
		break

	case " ":
	case "\r":
	case "\t":
		break

	case "\n":
		s.line++
		break
	case "\"":
		s.string()
		break
	default:
		if s.isDigit(char) {
			s.number()
		} else if s.isAlpha(char) {
			s.identifier()
		} else {
			Compiler.Error(s.line, "Unexpected Character.")
		}
	}
}
