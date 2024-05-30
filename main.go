package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

const (
	// Escape key
	Esc string = "\u001b"
	// Reset color
	Reset = Esc + "[0m"

	// Color codes
	Black   = Esc + "[30m"
	Red     = Esc + "[31m"
	Green   = Esc + "[32m"
	Yellow  = Esc + "[33m"
	Blue    = Esc + "[34m"
	Magenta = Esc + "[35m"
	Cyan    = Esc + "[36m"
	White   = Esc + "[37m"
	Default = Esc + "[39m"
)

type ColorExpression struct {
	Pattern   string
	ColorCode string
}

func parseColorExpression(expStr string) (ColorExpression, error) {
	idx := strings.LastIndex(expStr, ":")

	exp := ColorExpression{}

	if idx == -1 {
		return exp, errors.New("`:` not found in color expression")
	}

	pattern := expStr[:idx]
	color := expStr[idx+1:]

	// Map color to color code
	var colorCode string
	switch strings.ToLower(color) {
	case "black":
		colorCode = Black
	case "red":
		colorCode = Red
	case "green":
		colorCode = Green
	case "yellow":
		colorCode = Yellow
	case "blue":
		colorCode = Blue
	case "magenta":
		colorCode = Magenta
	case "cyan":
		colorCode = Cyan
	case "white":
		colorCode = White
	case "default":
		colorCode = Default
	case "":
		colorCode = Default
	default:
		return exp, errors.New("invalid color name")
	}

	exp.Pattern = pattern
	exp.ColorCode = colorCode

	return exp, nil
}

func main() {
	pflag.Parse()

	for i, arg := range pflag.Args() {
		fmt.Printf("Arg %d: %s\n", i, arg)
		exp, err := parseColorExpression(arg)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("Expression: %s%s%s\n", exp.ColorCode, exp.Pattern, Reset)
	}
}
