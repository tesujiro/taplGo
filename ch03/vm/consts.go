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
		ret = append(ret, cond)
		ret = append(ret, thn)
		ret = append(ret, els)
		return ret, nil
	default:
		return nil, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
}
