package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

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

type ColorExpression struct {
	Pattern   *regexp.Regexp
	ColorCode string
}

func parseColorExpression(expStr string) (ColorExpression, error) {
	idx := strings.LastIndex(expStr, ":")

	exp := ColorExpression{}

	if idx == -1 {
		return exp, errors.New("`:` not found in color expression")
	}

	patternStr := expStr[:idx]
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

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return exp, err
	}

	exp.Pattern = pattern
	exp.ColorCode = colorCode

	return exp, nil
}

func main() {
	pflag.Parse()

	exps := []ColorExpression{}

	for i, arg := range pflag.Args() {
		fmt.Printf("Arg %d: %s\n", i, arg)
		exp, err := parseColorExpression(arg)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		exps = append(exps, exp)
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Ensure stdin is from pipe, and not terminal
	// I am not sure how this works...
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("please pipe something from stdin")
		os.Exit(1)
	}

	for scanner.Scan() {
		text := scanner.Text()
		for _, exp := range exps {
			if exp.Pattern.MatchString(text) {
				text = fmt.Sprint(exp.ColorCode, text, Reset)
				break
			}
		}

		fmt.Print(text, "\n")
	}
}
