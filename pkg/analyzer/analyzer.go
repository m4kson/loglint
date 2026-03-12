package analyzer

import (
	"github.com/m4kson/loglint/pkg/analyzer/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/m4kson/loglint/pkg/analyzer/detector"
)

const (
	Name = "loglint"
	Doc  = "checks that log messages follow the project conventions"
)

type rule interface {
	Name() string
	Check(pass *analysis.Pass, call detector.Call)
}

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     Name,
		Doc:      Doc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	calls := detector.Detect(pass)
	if len(calls) == 0 {
		return nil, nil
	}

	for _, r := range registeredRules() {
		for _, call := range calls {
			r.Check(pass, call)
		}
	}

	return nil, nil
}

func registeredRules() []rule {
	return []rule{
		&rules.LowercaseRule{},
		&rules.EnglishOnlyRule{},
		&rules.NoSpecialCharsRule{},
	}
}
