%{
	package parser
	import (
		//"fmt"
		"github.com/tesujiro/taplGo/ch03/ast"
	)
%}

%union{
	token	ast.Token
	t	ast.Term
}

%type	<t>		program
%type	<t>		t

%token	<token>		TRUE FALSE
%token	<token>		ZERO
%token	<token>		IF THEN ELSE
%token	<token>		SUCC PRED
%token	<token>		ISZERO

%%

program
	: t ';'
	{
		$$ = $1
		yylex.(*Lexer).result = $$
	}

t
	: TRUE
	{
		$$ = &ast.BoolTerm{ Value: true }
	}
	| FALSE
	{
		$$ = &ast.BoolTerm{ Value: false }
	}
	| IF t THEN t ELSE t
	{
		$$ = &ast.IfTerm{Cond: $2, Then: $4, Else: $6}
	}
	| ZERO
	{
		$$ = &ast.ConstTerm{ Value: 0 }
	}
	| SUCC t
	{
		$$ = &ast.SuccTerm{ Arg: $2 }
	}
	| PRED t
	{
		$$ = &ast.PredTerm{ Arg: $2 }
	}
	| ISZERO t
	{
		$$ = &ast.IsZeroTerm{ Arg: $2 }
	}
	
/*
term
	: semi
	| term semi

semi
	: ';'
*/

%%
