package ast

type Term interface{}

type BoolTerm struct {
	Value bool
}

type ConstTerm struct {
	Value int
}

type IfTerm struct {
	Cond Term
	Then Term
	Else Term
}

type NextNumTerm struct {
	Arg  Term
	Diff int
}

type IsZeroTerm struct {
	Arg Term
}
