package colorexpression

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ANSI Escape Sequences
const (
	// Escape key
	Esc string = "\u001b"
	// Reset color
	Reset = Esc + "[0m"

	// Color codes
	Black   = Esc + "[30m" + Esc + "[47m"
	Red     = Esc + "[31m"
	Green   = Esc + "[32m"
	Yellow  = Esc + "[33m"
	Blue    = Esc + "[34m"
	Magenta = Esc + "[35m"
	Cyan    = Esc + "[36m"
	White   = Esc + "[37m"
	Default = Esc + "[39m"
)

type Selector int
const (
	Line Selector = -1
)

type ColorExpression struct {
	Selector Selector
	Pattern   *regexp.Regexp
	ColorCode string
}

func ParseColorExpression(expStr string) (ColorExpression, error) {
	firstIdx := strings.Index(expStr, "/")
	lastIdx := strings.LastIndex(expStr, "/")

	exp := ColorExpression{}

	if firstIdx == -1 {
		return exp, errors.New("`/` not found in color expression")
	}

	selector := Line

	if firstIdx == lastIdx {
		firstIdx = -1
	} else {
		var err error
		selectorInt, err := strconv.Atoi(expStr[0:firstIdx])
		if err != nil {
			return exp, errors.New("invalid selector")
		}
		selector = Selector(selectorInt)
	}

	patternStr := expStr[firstIdx+1 : lastIdx]
	color := expStr[lastIdx+1:]

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

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return exp, err
	}

	exp.Selector = selector
	exp.Pattern = pattern
	exp.ColorCode = colorCode

	return exp, nil
}
