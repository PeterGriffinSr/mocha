open Mocha

let () =
  let ast = Parser.program Lexer.token (Lexing.from_channel stdin) in
  List.iter (fun e -> print_endline (Ast.show_sexp e)) ast
