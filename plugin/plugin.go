package main

import (
	"github.com/ziv-nagish/commentspace"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		commentspace.NewAnalyzer(),
	}, nil
}
