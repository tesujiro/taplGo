package vm

import (
	"fmt"
	"reflect"

	"github.com/tesujiro/taplGo/ch03/ast"
)

func Depth(term ast.Term) (int, error) {
	return depth(term)
}

// Ch.3.3.2
func depth(term ast.Term) (int, error) {
	switch term.(type) {
	case *ast.BoolTerm:
		return 1, nil
	case *ast.ConstTerm:
		return 1, nil
	case *ast.NextNumTerm:
		s, err := Depth(term.(*ast.NextNumTerm).Arg)
		if err != nil {
			return 0, err
		}
		return s + 1, nil
	case *ast.IsZeroTerm:
		s, err := Depth(term.(*ast.IsZeroTerm).Arg)
		if err != nil {
			return 0, err
		}
		return s + 1, nil
	case *ast.IfTerm:
		cond, err := Depth(term.(*ast.IfTerm).Cond)
		if err != nil {
			return 0, err
		}
		thn, err := Depth(term.(*ast.IfTerm).Then)
		if err != nil {
			return 0, err
		}
		els, err := Depth(term.(*ast.IfTerm).Else)
		if err != nil {
			return 0, err
		}
		max := cond
		list := []int{cond, thn, els}
		for _, e := range list {
			if e > max {
				max = e
			}
		}
		return max, nil
	default:
		return 0, fmt.Errorf("invalid expression: %v", reflect.TypeOf(term))
	}
}
