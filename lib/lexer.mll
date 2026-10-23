{
    open Parser
}

let digit = ['0'-'9']
let sym_init = ['a'-'z' 'A'-'Z' '%' '&' '*' '+' '/' '-' '<' '=' '>' '?' '_']
let sym_subseq = sym_init | digit | ['.' '@']
let symbol = sym_init sym_subseq*
let int = '-'? digit+

rule token = parse
    | [' ' '\t' '\r']+ { token lexbuf }
    | '\n' { Lexing.new_line lexbuf; token lexbuf }
    | '(' { LeftParen }
    | ')' { RightParen }
    | '\'' { Quote }
    | int as n { Int (Z.of_string n) }
    | symbol as s { Symbol s }
    | eof { Eof }