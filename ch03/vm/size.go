package vm

import (
	"fmt"
	"reflect"

	"github.com/tesujiro/taplGo/ch03/ast"
)

func Size(term ast.Term) (int, error) {
	return size(term)
}

// Ch.3.3.2
func size(term ast.Term) (int, error) {
	switch term.(type) {
	case *ast.BoolTerm:
		return 1, nil
	case *ast.ConstTerm:
		return 1, nil
	case *ast.NextNumTerm:
		s, err := Size(term.(*ast.NextNumTerm).Arg)
		if err != nil {
			return 0, err
		}
		return s + 1, nil
	case *ast.IsZeroTerm:
		s, err := Size(term.(*ast.IsZeroTerm).Arg)
		if err != nil {
			return 0, err
		}
		return s + 1, nil
	case *ast.IfTerm:
		cond, err := Size(term.(*ast.IfTerm).Cond)
		if err != nil {
			return 0, err
		}
		thn, err := Size(term.(*ast.IfTerm).Then)
		if err != nil {
			return 0, err
		}
		els, err := Size(term.(*ast.IfTerm).Else)
		if err != nil {
			return 0, err
		}
		return cond + thn + els, nil
	default:
		return 0, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
}
