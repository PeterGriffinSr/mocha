open Mocha

let () =
  In_channel.input_all stdin |> Lexing.from_string |> Parser.program Lexer.token
  |> List.iter (fun e -> print_endline (Ast.show_exp e))
