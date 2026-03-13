package analyzer_test

import (
	"testing"

	"github.com/m4kson/loglint/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer_lowercase(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "lowercase")
}

func TestAnalyzer_englishOnly(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "englishonly")
}

func TestAnalyzer_specialChars(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "specialchars")
}

func TestAnalyzer_emoji(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "emoji")
}

func TestAnalyzer_sensitive(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "sensitive")
}

func TestAnalyzer_clean(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), analyzer.New(), "clean")
}
