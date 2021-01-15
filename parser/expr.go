package parser

import "GNBS/token"

type Expr interface {
	accept(Visitor) string
}

type Visitor interface {
	visitUnaryExpr(unary) string
	visitBinaryExpr(binary) string
	visitGroupingExpr(grouping) string
	visitLiteralExpr(literal) string
}

type unary struct {
	operator token.Token
	right    Expr
}

func (u unary) accept(visitor Visitor) string {
	return visitor.visitUnaryExpr(u)
}

type binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (b binary) accept(visitor Visitor) string {
	return visitor.visitBinaryExpr(b)
}

type grouping struct {
	expressions Expr
}

func (g grouping) accept(visitor Visitor) string {
	return visitor.visitGroupingExpr(g)
}

type literal struct {
	value interface{}
}

func (l literal) accept(visitor Visitor) string {
	return visitor.visitLiteralExpr(l)
}
