package rules

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"

	"github.com/m4kson/loglint/pkg/analyzer/detector"
)

type LowercaseRule struct{}

func (r *LowercaseRule) Name() string { return "lowercase" }

func (r *LowercaseRule) Check(pass *analysis.Pass, call detector.Call) {
	msg := call.MsgValue
	if msg == "" {
		return
	}

	firstRune, size := utf8.DecodeRuneInString(msg)
	if firstRune == utf8.RuneError {
		return
	}

	if !unicode.IsUpper(firstRune) {
		return
	}

	lowered := strings.ToLower(string(firstRune))

	msgContentStart := call.MsgPos + 1

	pass.Report(analysis.Diagnostic{
		Pos:     call.MsgPos,
		End:     call.MsgLit.End(),
		Message: `log message must start with a lowercase letter`,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     msgContentStart,
						End:     msgContentStart + token.Pos(size),
						NewText: []byte(lowered),
					},
				},
			},
		},
	})
}
