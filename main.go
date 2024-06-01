package main

import (
	"bufio"
	"bytes"
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

type colorExpression struct {
	Pattern   *regexp.Regexp
	ColorCode string
}

func main() {
	var separator string
	pflag.StringVarP(&separator, "separator", "s", "\n", "test")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... EXPRESSION...\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n") // Newline

		fmt.Fprintln(os.Stderr,
`Expressions:
  Expressions are in the form <regex>:<color>.
  Valid colors are black, red, green, yellow,
  blue, magenta, cyan, white, and default.`,
		)
	}
	pflag.Parse()

	exps := []colorExpression{}

	for _, arg := range pflag.Args() {
		exp, err := parseColorExpression(arg)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		exps = append(exps, exp)
	}

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

	scanner := bufio.NewScanner(os.Stdin)
	if separator != "\n" {
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			// If we are at the end of file, and there is no more data, return nothing.
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.Index(data, []byte(separator)); i >= 0 {
				// We have a full token, return it and skip the length of the separator
				return i + len(separator), data[0:i], nil
			}
			// If we're at EOF, we have a final, non-terminated line. Return it.
			if atEOF {
				return len(data), data, nil
			}
			// Request more data.
			return 0, nil, nil
		})
	}

	for scanner.Scan() {
		text := scanner.Text()
		for _, exp := range exps {
			if exp.Pattern.MatchString(text) {
				text = fmt.Sprint(exp.ColorCode, text, Reset)
				break
			}
		}

		fmt.Print(text, separator)
	}
}

func parseColorExpression(expStr string) (colorExpression, error) {
	idx := strings.LastIndex(expStr, ":")

	exp := colorExpression{}

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

