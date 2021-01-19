package parser

import "GNBS/token"

type Expr interface {
	accept(Visitor) interface{}
}

type Visitor interface {
	visitUnaryExpr(unary) interface{}
	visitBinaryExpr(binary) interface{}
	visitGroupingExpr(grouping) interface{}
	visitLiteralExpr(literal) interface{}
}

type unary struct {
	operator token.Token
	right    Expr
}

func (u unary) accept(visitor Visitor) interface{} {
	return visitor.visitUnaryExpr(u)
}

type binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (b binary) accept(visitor Visitor) interface{} {
	return visitor.visitBinaryExpr(b)
}

type grouping struct {
	expressions Expr
}

func (g grouping) accept(visitor Visitor) interface{} {
	return visitor.visitGroupingExpr(g)
}

type literal struct {
	value interface{}
}

func (l literal) accept(visitor Visitor) interface{} {
	return visitor.visitLiteralExpr(l)
}
