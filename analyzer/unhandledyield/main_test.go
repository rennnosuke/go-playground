package unhandledyield_test

import (
	"testing"

	"github.com/rennnosuke/go-playground/analyzer/unhandledyield"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUnhandledYield(t *testing.T) {
	testData := analysistest.TestData()
	analysistest.Run(t, testData, unhandledyield.Analyzer, "a")
}
