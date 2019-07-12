package vm

import (
	"fmt"
	"reflect"

	"github.com/tesujiro/taplGo/ch03/ast"
)

func Run(term ast.Term) (interface{}, error) {
	return evalTerm(term)
}

func evalTerm(term ast.Term) (interface{}, error) {
	//fmt.Printf("evalExpr(%#v)\n", term)
	switch term.(type) {
	case *ast.BoolTerm:
		return term.(*ast.BoolTerm).Value, nil
	case *ast.ConstTerm:
		return term.(*ast.ConstTerm).Value, nil
	case *ast.IfTerm:
		cond := term.(*ast.IfTerm).Cond
		thn := term.(*ast.IfTerm).Then
		els := term.(*ast.IfTerm).Else
		res, err := evalTerm(cond)
		if err != nil {
			return nil, fmt.Errorf("if condition error: %v", err)
		}
		if val, ok := res.(bool); ok && val == true {
			return evalTerm(thn)
		} else {
			return evalTerm(els)
		}
	case *ast.SuccTerm:
	case *ast.PredTerm:
	case *ast.IsZeroTerm:

	default:
		return nil, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
	return nil, nil
}
