package widgets

import (
	"strings"

	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
)

// ref. github.com/evertras/bubble-table@v0.14.6/table/strlimit.go.
func StringWithTail(str string, maxLen int) string {
	if maxLen == 0 {
		return ""
	}

	newLineIndex := strings.Index(str, "\n")
	if newLineIndex > -1 {
		str = str[:newLineIndex] + "…"
	}

	if ansi.PrintableRuneWidth(str) > maxLen {
		return truncate.StringWithTail(str, uint(maxLen), "…")
	}

	return str
}
