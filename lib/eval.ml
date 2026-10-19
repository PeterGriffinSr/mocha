open Ast

exception Eval_error of string

module Globals = Map.Make (String)

type globals = exp Globals.t

let err fmt = Printf.ksprintf (fun s -> raise (Eval_error s)) fmt
let truthy = function Data ("Nil", []) -> false | _ -> true
let truth b = if b then Sym "t" else nil

let list_of v =
  let rec go acc = function
    | Data ("Nil", []) -> List.rev acc
    | Data ("Cons", [ a; d ]) -> go (a :: acc) d
    | v -> err "improper list ending in %s" (show_exp v)
  in
  go [] v

let to_list xs = List.fold_right cons xs nil

let lookup genv env s =
  match List.assoc_opt s env with
  | Some v -> v
  | None -> (
      match Globals.find_opt s genv with
      | Some v -> v
      | None -> err "unbound symbol: %s" s)

let rec bind params args env =
  match (params, args) with
  | Data ("Nil", []), [] -> env
  | Sym p, rest -> (p, to_list rest) :: env
  | Data ("Cons", [ Sym p; ps ]), a :: rest -> bind ps rest ((p, a) :: env)
  | _ -> err "arity mismatch"

let rec eval genv env = function
  | (Float _ | Closure _ | Prim _ | Macro _) as v -> v
  | Sym s -> lookup genv env s
  | Data ("Cons", [ f; args ]) -> (
      match (f, list_of args) with
      | Sym "quote", [ x ] -> x
      | Sym "if", [ c; t; e ] ->
          eval genv env (if truthy (eval genv env c) then t else e)
      | Sym "lambda", [ params; body ] -> Closure { params; body; env }
      | _, args -> (
          match eval genv env f with
          | Macro m -> eval genv env (apply genv (Closure m) args)
          | fv -> apply genv fv (List.map (eval genv env) args)))
  | Data _ as v -> v

and apply genv f args =
  match f with
  | Closure c -> eval genv (bind c.params args c.env) c.body
  | Prim (_, p) -> p args
  | v -> err "not a function: %s" (show_exp v)

let p1 f = function [ a ] -> f a | _ -> err "expected 1 argument"
let p2 f = function [ a; b ] -> f a b | _ -> err "expected 2 arguments"

let num = function
  | Float f -> f
  | v -> err "expected a number, got %s" (show_exp v)

let float2 op a b = Float (op (num a) (num b))
let cmp2 op a b = truth (op (num a) (num b))

let prims =
  [
    ("cons", p2 cons);
    ("eq?", p2 (fun a b -> truth (equal_exp a b)));
    ("+", p2 (float2 ( +. )));
    ("-", p2 (float2 ( -. )));
    ("*", p2 (float2 ( *. )));
    ("<", p2 (cmp2 ( < )));
    ( "number?",
      p1 (fun v -> truth (match v with Float _ -> true | _ -> false)) );
    ("data?", p1 (fun v -> truth (match v with Data _ -> true | _ -> false)));
    ( "make-data",
      function
      | Sym tag :: fields -> Data (tag, fields)
      | _ -> err "make-data expects a tag symbol first" );
    ( "data-tag",
      p1 (function
        | Data (t, _) -> Sym t
        | _ -> err "data-tag expects a data value") );
    ( "data-ref",
      p2 (fun d i ->
          match (d, i) with
          | Data (_, fs), Float i -> (
              match List.nth_opt fs (int_of_float i) with
              | Some v -> v
              | None -> err "data-ref: index %g out of range" i)
          | _ -> err "data-ref expects (data index)") );
  ]

let initial_globals : globals =
  List.fold_left
    (fun g (name, f) -> Globals.add name (Prim (name, f)) g)
    (Globals.of_seq
       (List.to_seq [ ("t", Sym "t"); ("else", Sym "t"); ("nil", nil) ]))
    prims

let define genv name v =
  if Globals.mem name genv then err "already defined: %s" name;
  (Globals.add name v genv, Sym name)

let global_macro genv = function
  | Sym s -> (
      match Globals.find_opt s genv with Some (Macro m) -> Some m | _ -> None)
  | _ -> None

let rec toplevel genv form : globals * exp =
  match form with
  | Data ("Cons", [ Sym "begin"; body ]) ->
      List.fold_left (fun (g, _) f -> toplevel g f) (genv, nil) (list_of body)
  | Data ("Cons", [ Sym (("define" | "defmacro") as f); args ]) -> (
      match (f, list_of args) with
      | "define", [ Sym n; e ] -> define genv n (eval genv [] e)
      | "defmacro", [ Sym n; params; body ] ->
          define genv n (Macro { params; body; env = [] })
      | _ -> err "malformed %s: %s" f (show_exp form))
  | Data ("Cons", [ head; args ]) -> (
      match global_macro genv head with
      | Some m -> toplevel genv (apply genv (Closure m) (list_of args))
      | None -> (genv, eval genv [] form))
  | e -> (genv, eval genv [] e)

let run_string genv src =
  Parser.program Lexer.token (Lexing.from_string src)
  |> List.fold_left (fun (g, _) f -> toplevel g f) (genv, nil)

let load_file genv path =
  run_string genv (In_channel.with_open_bin path In_channel.input_all)
