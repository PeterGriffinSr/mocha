%{
open Ast

let wrap tag x = cons (Sym tag) (cons x nil)
%}

%token LeftParen RightParen Rest
%token Quote
%token <float> Float
%token <string> Symbol
%token Eof

%start program
%type <Ast.exp list> program

%%

program:
    | Eof { [] }
    | sexp program { $1 :: $2 }
    ;

sexp:
    | Float { Float $1 }
    | Symbol { Sym $1 }
    | Quote sexp { wrap "quote" $2 }
    | LeftParen RightParen { nil }
    | LeftParen Rest sexp RightParen { $3 }
    | LeftParen sexp tail { cons $2 $3 }

tail:
    | RightParen { nil }
    | Rest sexp RightParen { $2 }
    | sexp tail { cons $1 $2 }
    ;