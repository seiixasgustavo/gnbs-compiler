package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type AstPrinter struct {
	Visitor
}

func (a AstPrinter) visitBinaryExpr(expr binary) interface{} {
	return a.parenthesize(expr.operator.ToString(), expr.left, expr.right)
}

func (a AstPrinter) visitGroupingExpr(expr grouping) interface{} {
	return a.parenthesize("group", expr.expressions)
}

func (a AstPrinter) visitUnaryExpr(expr unary) interface{} {
	return a.parenthesize(expr.operator.ToString(), expr.right)
}

func (a AstPrinter) visitLiteralExpr(expr literal) interface{} {
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
	return expr.accept(a).(string)
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.Write([]byte("("))
	builder.Write([]byte(name))

	for _, expr := range exprs {
		builder.Write([]byte(" "))
		builder.Write([]byte(expr.accept(a).(string)))
	}
	builder.Write([]byte(")"))

	return builder.String()
}
