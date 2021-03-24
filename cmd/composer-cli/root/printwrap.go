// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package root

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// PrintWrap prints a string wrapped at columns and subsequent lines indented
func PrintWrap(indent, columns int, s string) {
	words := strings.Split(strings.ReplaceAll(s, "\n", " "), " ")
	fmt.Print(words[0])
	col := indent + utf8.RuneCountInString(words[0])
	for _, w := range words[1:] {
		if col+utf8.RuneCountInString(w)+1 > columns {
			fmt.Printf("\n%*s%s", indent, "", w)
			col = indent + utf8.RuneCountInString(w)
		} else if len(w) > 0 {
			fmt.Print(" ", w)
			col = col + 1 + utf8.RuneCountInString(w)
		}
	}
	fmt.Println()
}
