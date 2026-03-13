package linters

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/m4kson/loglint/pkg/analyzer"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct{}

type Plugin struct {
	settings Settings
}

func New(raw any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](raw)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: settings}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.New(),
	}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
