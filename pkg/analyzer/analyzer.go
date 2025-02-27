package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Name = "commentspace"

// Analyzer is the commentspace analyzer
var Analyzer = &analysis.Analyzer{
	Name:     Name,
	Doc:      "Checks if single-line comments have a space after the // characters",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Process all comments in the file directly
	for _, file := range pass.Files {
		// Process file comments (package documentation)
		checkCommentGroup(file.Doc, pass)

		// Process all comments in the file
		for _, commentGroup := range file.Comments {
			checkCommentGroup(commentGroup, pass)
		}
	}

	// Additional specific node inspection for documentation comments
	// that might be interesting for deeper analysis
	nodeTypes := []ast.Node{
		(*ast.GenDecl)(nil),
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeTypes, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.GenDecl:
			checkCommentGroup(n.Doc, pass)
		case *ast.FuncDecl:
			checkCommentGroup(n.Doc, pass)
		}
	})

	return nil, nil
}

// checkCommentGroup processes a group of comments
func checkCommentGroup(group *ast.CommentGroup, pass *analysis.Pass) {
	if group == nil {
		return
	}

	for _, comment := range group.List {
		// Only check single-line comments (those starting with //)
		if !strings.HasPrefix(comment.Text, "//") {
			continue
		}

		// Skip directives like //go:generate
		if strings.HasPrefix(comment.Text, "//go:") {
			continue
		}

		// Skip URLs in comments (http:// or https://)
		if strings.HasPrefix(comment.Text, "//http://") || strings.HasPrefix(comment.Text, "//https://") {
			continue
		}

		// Check if there's a space after the // characters
		if len(comment.Text) > 2 && comment.Text[2] != ' ' {
			pass.Report(analysis.Diagnostic{
				Pos:     comment.Pos(),
				End:     comment.End(),
				Message: "comment should have a space after //",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Add space after //",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     token.Pos(int(comment.Pos()) + 2),
								End:     token.Pos(int(comment.Pos()) + 2),
								NewText: []byte(" "),
							},
						},
					},
				},
			})
		}
	}
}
