package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type AstPrinter struct {
	Visitor
}

func (a AstPrinter) visitBinaryExpr(expr binary) string {
	return a.parenthesize(expr.operator.ToString(), expr.left, expr.right)
}

func (a AstPrinter) visitGroupingExpr(expr grouping) string {
	return a.parenthesize("group", expr.expressions)
}

func (a AstPrinter) visitUnaryExpr(expr unary) string {
	return a.parenthesize(expr.operator.ToString(), expr.right)
}

func (a AstPrinter) visitLiteralExpr(expr literal) string {
	if expr.value == nil {
		return "null"
	}

	var s string
	var i int32
	var f float32

	switch reflect.TypeOf(expr.value) {
	case reflect.TypeOf(s):
		return fmt.Sprintf("\"%s\"", expr.value)
	case reflect.TypeOf(i):
		return fmt.Sprintf("%d", expr.value)
	case reflect.TypeOf(f):
		return fmt.Sprintf("%f", expr.value)
	}
	return fmt.Sprintf("%v", expr.value)
}

func (a AstPrinter) print(expr Expr) string {
	return expr.accept(a)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.Write([]byte("("))
	builder.Write([]byte(name))

	for _, expr := range exprs {
		builder.Write([]byte(" "))
		builder.Write([]byte(expr.accept(a)))
	}
	builder.Write([]byte(")"))

	return builder.String()
}
