package unhandledyield

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "unhandledyield",
	Doc:              "unhandledyield discover yield in range over func is not handled.",
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	ResultType:       nil,
	FactTypes:        nil,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// 以下の関数を対象にする
	// 1. 関数定義
	// 2. 無名関数
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			discoverFuncDecl(fn)
		case *ast.FuncLit:
			discoverFuncLit(fn)
		}
	})
	return nil, nil
}

func discoverFuncDecl(fn *ast.FuncDecl) {
	// TODO: 以下のrange over func のシグネチャに該当する関数かどうかを見る
	// func(func()bool)
	// func(func(V)bool)
	// func(func(K, V)bool)
	// TODO: 関数内にyield呼び出しがあるかどうか見る
	// TODO: yield呼び出しの返り値を評価し、値の返戻が関数内で行われないか見る
}

func discoverFuncLit(fn *ast.FuncLit) {
	// TODO: 以下のrange over func のシグネチャに該当する関数かどうかを見る
	// func(func()bool)
	// func(func(V)bool)
	// func(func(K, V)bool)
}
