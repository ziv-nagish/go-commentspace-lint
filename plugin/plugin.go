package main

import (
	"github.com/ziv-nagish/go-commentspace-lint"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		commentspace.NewAnalyzer(),
	}, nil
}
