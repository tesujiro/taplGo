package vm

import (
	"fmt"
	"reflect"

	"github.com/tesujiro/taplGo/ch03/ast"
)

func Consts(term ast.Term) ([]interface{}, error) {
	return consts(term)
}

// Ch.3.3.1
func consts(term ast.Term) ([]interface{}, error) {
	switch term.(type) {
	case *ast.BoolTerm:
		return []interface{}{term.(*ast.BoolTerm).Value}, nil
	case *ast.ConstTerm:
		return []interface{}{term.(*ast.ConstTerm).Value}, nil
	case *ast.NextNumTerm:
		return Consts(term.(*ast.NextNumTerm).Arg)
	case *ast.IsZeroTerm:
		return Consts(term.(*ast.IsZeroTerm).Arg)
	case *ast.IfTerm:
		cond, err := Consts(term.(*ast.IfTerm).Cond)
		if err != nil {
			return nil, err
		}
		thn, err := Consts(term.(*ast.IfTerm).Then)
		if err != nil {
			return nil, err
		}
		els, err := Consts(term.(*ast.IfTerm).Else)
		if err != nil {
			return nil, err
		}
		ret := []interface{}{}
		union := func(a, b []interface{}) []interface{} {
			var result []interface{}
			result = append(result, a...)
		loop:
			for _, e1 := range b {
				for _, e2 := range a {
					if e1 == e2 {
						continue loop
					}
				}
				result = append(result, e1)
			}

			return result
		}
		ret = union(ret, cond)
		ret = union(ret, thn)
		ret = union(ret, els)
		return ret, nil
	default:
		return nil, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
}
