package GNBS

import "fmt"

var Compiler *gnbs

type gnbs struct {
	hadError bool
}

func (g *gnbs) Error(line uint, message string) {
	g.report(line, "", message)
}

func (g *gnbs) report(line uint, where, message string) {
	fmt.Printf("[line: %d] Error %s: %s", line, where, message)
	g.hadError = true
}
