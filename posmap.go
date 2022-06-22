package buildposmap

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type PosMap map[token.Pos][]ast.Node

var Analyzer = &analysis.Analyzer{
	Name:             "buildposmap",
	Doc:              "build mapping between token.Pos and ast.Node",
	Run:              run,
	RunDespiteErrors: true,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	ResultType:       reflect.TypeOf(PosMap(nil)),
}

func run(pass *analysis.Pass) (interface{}, error) {
	return New(pass), nil
}

func New(pass *analysis.Pass) PosMap {
	posMap := make(PosMap)
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			posMap[i] = append(posMap[i], node)
		}
	})
	return posMap
}

func Node[T ast.Node](posMap PosMap, pos token.Pos) (node T, ok bool) {
	for i := pos; i > 0; i-- {
		stack := posMap[i]
		if len(stack) == 0 {
			break
		}
		for j := range stack {
			node := stack[len(stack)-1-j]
			ident, ok := node.(T)
			if ok {
				return ident, true
			}
		}
	}
	var zero T
	return zero, false
}
