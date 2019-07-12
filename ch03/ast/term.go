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

type SuccTerm struct {
	Arg Term
}

type PredTerm struct {
	Arg Term
}

type IsZeroTerm struct {
	Arg Term
}
