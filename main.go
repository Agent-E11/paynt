package main

import (
	"bufio"
	"bytes"

	//"errors"
	"fmt"
	"os"

	//"regexp"
	//"strconv"
	"strings"

	colexp "github.com/agent-e11/paynt/colorexpression"
	"github.com/spf13/pflag"
)

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
  Expressions are in the form [selector/]<regex>/[color].
  Any number of expressions can be provided. If multiple
  expressions match a line, the leftmost expression will
  be applied.

  selector
    A number defining what the regex should match against.
    If it is greater than -1, then the field of that index
    is matched against.
    If it is -1, then the entire line is matched.
    If it is omitted, then it will default to -1.

  regex
    Regular expression to match against the selected text.
    Can be any valid syntax as defined in Go's regexp/syntax
    package:
    https://pkg.go.dev/regexp/syntax

  color
    Color to color the line if the regex matches.
    Can be black, red, green, yellow, blue, magenta,
    cyan, white, or default.
    If it is omitted, then it will default to default.`)
	}
	pflag.Parse()

	exps := []colexp.ColorExpression{}

	for _, arg := range pflag.Args() {
		exp, err := colexp.ParseColorExpression(arg)
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
		line := scanner.Text()
		for _, exp := range exps {
			matchText := line
			if exp.Selector > -1 {
				matchText = strings.Split(line, " ")[exp.Selector]
			}
			if exp.Pattern.MatchString(matchText) {
				line = fmt.Sprint(exp.ColorCode, line, colexp.Reset)
				break
			}
		}

		fmt.Print(line, separator)
	}
}

