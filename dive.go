package dive

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "dive finds low readability if-blocks"

var Analyzer = &analysis.Analyzer{
	Name: "dive",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var (
	flagBlockLength int
	flagReturns     int
	flagDepthNest   int
)

func init() {
	Analyzer.Flags.IntVar(&flagBlockLength, "len", 5, "max length of if block")
	Analyzer.Flags.IntVar(&flagReturns, "ret", 2, "max number of returns in a if block")
	Analyzer.Flags.IntVar(&flagDepthNest, "nest", 2, "max depth of nest")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			checkTopBlock(pass, n.Body)
		case *ast.FuncLit:
			checkTopBlock(pass, n.Body)
		}
	})

	return nil, nil
}

func checkTopBlock(pass *analysis.Pass, body *ast.BlockStmt) {
	if body == nil {
		return
	}

	for _, stmt := range body.List {
		switch stmt := stmt.(type) {
		case *ast.IfStmt:
			checkIf(pass, stmt)
		}
	}
}

func checkIf(pass *analysis.Pass, ifstmt *ast.IfStmt) {
	checkBlock(pass, ifstmt.Body)
	switch elsestmt := ifstmt.Else.(type) {
	case *ast.BlockStmt:
		checkBlock(pass, elsestmt)
	case *ast.IfStmt:
		checkIf(pass, elsestmt)
	}
}

func checkBlock(pass *analysis.Pass, body *ast.BlockStmt) {
	if body == nil || len(body.List) == 0 {
		return
	}

	checkLongBlock(pass, body)
	checkManyReturns(pass, body)
	checkHasLoop(pass, body)
	checkDeeplyNest(pass, body)
}

func checkLongBlock(pass *analysis.Pass, body *ast.BlockStmt) {
	if body != nil && len(body.List) > flagBlockLength {
		pass.Reportf(body.Pos(), "too long block")
	}
}

func countLength(body *ast.BlockStmt) int {
	if body == nil {
		return 0
	}

	var count int

	for _, stmt := range body.List {
		switch stmt := stmt.(type) {
		case *ast.ForStmt:
			count += 1 + countLength(stmt.Body)
		case *ast.IfStmt:
			count += 1 + countLength(stmt.Body)
		case *ast.SwitchStmt:
			count += 1 + countLength(stmt.Body)
		case *ast.TypeSwitchStmt:
			count += 1 + countLength(stmt.Body)
		case *ast.SelectStmt:
			count += 1 + countLength(stmt.Body)
		default:
			count++
		}
	}

	return count
}

func checkManyReturns(pass *analysis.Pass, body *ast.BlockStmt) {
	var count int
	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit, *ast.SelectStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			return false
		case *ast.ReturnStmt:
			count++
		}
		return true
	})
	if count > flagReturns {
		pass.Reportf(body.Pos(), "too many returns in the block")
	}
}

func checkHasLoop(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ForStmt:
			pass.Reportf(body.Pos(), "loop in if block")
		}
		return true
	})
}

func checkDeeplyNest(pass *analysis.Pass, body *ast.BlockStmt) {
	if 1+countDepth(body) > flagDepthNest {
		pass.Reportf(body.Pos(), "too deeply nest")
	}
}

func countDepth(body *ast.BlockStmt) int {
	if body == nil {
		return 0
	}

	for _, stmt := range body.List {
		switch stmt := stmt.(type) {
		case *ast.ForStmt:
			return 1 + countDepth(stmt.Body)
		case *ast.IfStmt:
			return 1 + countDepth(stmt.Body)
		case *ast.SwitchStmt:
			return 1 + countDepth(stmt.Body)
		case *ast.TypeSwitchStmt:
			return 1 + countDepth(stmt.Body)
		case *ast.SelectStmt:
			return 1 + countDepth(stmt.Body)
		}
	}

	return 0
}
