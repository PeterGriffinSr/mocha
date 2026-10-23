%{
    open Ast
%}

%token LeftParen RightParen Quote Eof
%token <Z.t> Int
%token <string> Symbol

%start program
%type <Ast.exp list> program

%%

program:
    es = sexp* Eof { es }

sexp:
    | n = Int { Int n }
    | s = Symbol { Sym s }
    | Quote e = sexp { List [ Sym "quote"; e ] }
    | LeftParen es = sexp* RightParen { List es }