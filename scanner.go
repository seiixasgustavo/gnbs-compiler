package GNBS

import (
	"strconv"
)

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)

	keywords[tokens[And]] = And
	keywords[tokens[Or]] = Or
	keywords[tokens[If]] = If
	keywords[tokens[Else]] = Else
	keywords[tokens[Var]] = Var
	keywords[tokens[Func]] = Func
	keywords[tokens[Struct]] = Struct
	keywords[tokens[Return]] = Return
	keywords[tokens[True]] = True
	keywords[tokens[False]] = False
}

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
	s.addTokenToList(String, value)
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
		s.addTokenToList(Float, value)
	} else {
		value, _ := strconv.ParseInt(s.source[s.start:s.current], 10, 64)
		s.addTokenToList(Integer, value)
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	s.addToken(Identifier)
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
	return uint(len(tokens)) <= s.current
}

func (s *Scanner) addToken(token TokenType) {
	s.addTokenToList(token, nil)
}

func (s *Scanner) addTokenToList(token TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(token, text, literal, s.line))
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, *NewToken(Eof, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case "(":
		s.addToken(LParentheses)
		break
	case ")":
		s.addToken(RParentheses)
		break
	case "{":
		s.addToken(LBrace)
		break
	case "}":
		s.addToken(RBrace)
		break
	case ",":
		s.addToken(Comma)
		break
	case ".":
		s.addToken(Dot)
		break
	case "-":
		s.addToken(Minus)
		break
	case "+":
		s.addToken(Plus)
		break
	case ";":
		s.addToken(Semicolon)
		break
	case "*":
		s.addToken(Star)
		break
	case "!":
		if s.match("=") {
			s.addToken(NotEqual)
		} else {
			s.addToken(Not)
		}
		break
	case "=":
		if s.match("=") {
			s.addToken(EqualEqual)
		} else {
			s.addToken(Equal)
		}
		break
	case ">":
		if s.match("=") {
			s.addToken(GreaterEqual)
		} else {
			s.addToken(Greater)
		}
		break
	case "<":
		if s.match("=") {
			s.addToken(LessEqual)
		} else {
			s.addToken(Less)
		}
		break
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash)
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
