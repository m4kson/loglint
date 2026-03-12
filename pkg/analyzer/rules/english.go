package rules

import (
	"go/token"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"

	"github.com/m4kson/loglint/pkg/analyzer/detector"
)

type EnglishOnlyRule struct{}

// Name returns the rule identifier.
func (r *EnglishOnlyRule) Name() string { return "english-only" }

func (r *EnglishOnlyRule) Check(pass *analysis.Pass, call detector.Call) {
	msg := call.MsgValue
	if msg == "" {
		return
	}

	byteOffset := 0

	for _, ch := range msg {
		if ch == utf8.RuneError {
			break
		}

		if ch < 0x20 || ch > 0x7E {
			offendingPos := call.MsgPos + 1 + token.Pos(byteOffset)

			pass.Reportf(
				offendingPos,
				"log message must contain only English (ASCII) characters, found %q",
				string(ch),
			)
			return
		}

		byteOffset += utf8.RuneLen(ch)
	}
}
