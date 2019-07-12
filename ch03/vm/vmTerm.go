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
	case *ast.NextNumTerm:
		arg := term.(*ast.NextNumTerm).Arg
		diff := term.(*ast.NextNumTerm).Diff
		res, err := evalTerm(arg)
		if err != nil {
			return nil, fmt.Errorf("if condition error: %v", err)
		}
		if val, ok := res.(int); ok {
			return val + diff, nil
		} else {
			return diff, nil
		}
	case *ast.IsZeroTerm:
		arg := term.(*ast.IsZeroTerm).Arg
		res, err := evalTerm(arg)
		if err != nil {
			return nil, fmt.Errorf("if condition error: %v", err)
		}
		if val, ok := res.(int); ok {
			return val == 0, nil
		} else {
			return false, nil
		}

	default:
		return nil, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
}
