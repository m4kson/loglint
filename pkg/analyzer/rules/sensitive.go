package rules

import (
	"strings"

	"github.com/m4kson/loglint/pkg/analyzer/detector"
	"golang.org/x/tools/go/analysis"
)

type NoSensitiveDataRule struct{}

func (r *NoSensitiveDataRule) Name() string { return "no-sensitive-data" }

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"apikey",
	"api_key",
	"api-key",
	"auth",
	"credential",
	"private_key",
	"privatekey",
	"access_key",
	"accesskey",
	"bearer",
	"jwt",
	"ssn",
	"credit_card",
	"creditcard",
	"cvv",
	"pin",
}

func (r *NoSensitiveDataRule) Check(pass *analysis.Pass, call detector.Call) {
	msg := call.MsgValue
	if msg == "" {
		return
	}

	lower := strings.ToLower(msg)

	keywords := defaultSensitiveKeywords

	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			pass.Reportf(
				call.MsgPos,
				"log message may contain sensitive data: keyword %q found in message",
				kw,
			)
			return
		}
	}
}
