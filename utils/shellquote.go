package utils

import (
	"strconv"
	"strings"
)

// ShellQuoteArgs returns a shell-like representation of args for printing/logging.
// Note: this is meant for display, not for feeding back into exec.Command.
func ShellQuoteArgs(args ...string) string {
	quoted := make([]string, 0, len(args))
	for _, a := range args {
		// strconv.Quote produces a Go-style quoted string; for basic CLI display it's fine.
		// If there are no spaces/quotes, keep it as-is to stay readable.
		if strings.IndexFunc(a, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '"' || r == '\''
		}) >= 0 {
			quoted = append(quoted, strconv.Quote(a))
		} else {
			quoted = append(quoted, a)
		}
	}
	return strings.Join(quoted, " ")
}
