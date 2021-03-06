package compiler

import "GNBS/token"

type Precedence int

const (
	None Precedence = iota
	Assignment
	Or
	And
	Equality
	Comparison
	Term
	Factor
	Unary
	Call
	Primary
)

type ParseRule struct {
	prefix     func(bool)
	infix      func(bool)
	precedence Precedence
}

var rules []ParseRule

func init() {
	rules = []ParseRule{
		token.LParentheses: {grouping, callFn, Call},
		token.RParentheses: {nil, nil, None},
		token.LBrace:       {nil, nil, None},
		token.RBrace:       {nil, nil, None},
		token.Comma:        {nil, nil, None},
		token.Dot:          {nil, nil, None},
		token.Minus:        {unary, binary, Term},
		token.Plus:         {nil, grouping, Term},
		token.Semicolon:    {nil, nil, None},
		token.Slash:        {nil, binary, Factor},
		token.Star:         {nil, binary, Factor},
		token.Not:          {unary, nil, None},
		token.NotEqual:     {nil, binary, Equality},
		token.Equal:        {nil, nil, None},
		token.EqualEqual:   {nil, binary, Equality},
		token.Less:         {nil, binary, Comparison},
		token.LessEqual:    {nil, binary, Comparison},
		token.Greater:      {nil, binary, Comparison},
		token.GreaterEqual: {nil, binary, Comparison},
		token.Identifier:   {variable, nil, None},
		token.String:       {stringvalue, nil, None},
		token.Float:        {floatnumber, nil, None},
		token.Integer:      {intnumber, nil, None},
		token.And:          {nil, and_, And},
		token.Or:           {nil, or_, Or},
		token.Class:        {nil, nil, None},
		token.Function:     {nil, nil, None},
		token.True:         {literal, nil, None},
		token.False:        {literal, nil, None},
		token.For:          {nil, nil, None},
		token.If:           {nil, nil, None},
		token.Else:         {nil, nil, None},
		token.Null:         {nil, nil, None},
		token.Return:       {nil, nil, None},
		token.This:         {nil, nil, None},
		token.Error:        {nil, nil, None},
		token.Eof:          {nil, nil, None},
	}
}

func getRule(tokenType token.TokenType) *ParseRule {
	return &rules[tokenType]
}
