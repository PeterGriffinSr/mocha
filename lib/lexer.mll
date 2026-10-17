{
    open Parser
}

let digit = ['0'-'9']
let letter = ['a'-'z' 'A'-'Z']
let sym_init = letter | ['%' '&' '*' '+' '/' '-' '<' '=' '>' '?' '_']
let sym_subseq = sym_init | digit | ['.' '@']
let symbol = sym_init sym_subseq*
let float = '-'? digit+ ('.' digit+)? (['e' 'E'] ['+' '-']? digit+)?

rule token = parse 
    | [' ' '\t' '\r']+ { token lexbuf }
    | '\n' { Lexing.new_line lexbuf; token lexbuf }
    | '(' { LeftParen }
    | ')' { RightParen }
    | '\'' { Quote }
    | '&' { Rest }
    | float as n { Float (float_of_string n) }
    | symbol as s { Symbol s }
    | eof { Eof }