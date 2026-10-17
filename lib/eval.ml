open Ast

exception Eval_error of string

let err fmt = Printf.ksprintf (fun s -> raise (Eval_error s)) fmt
let globals = Hashtbl.create 64

let lookup env s =
  match List.assoc_opt s env with
  | Some v -> v
  | None -> (
      match Hashtbl.find_opt globals s with
      | Some v -> v
      | None -> err "unbound symbol: %s" s)

let define name v =
  if Hashtbl.mem globals name then err "already defined: %s" name;
  Hashtbl.add globals name v;
  Sym name

let truthy = function Nil -> false | _ -> true
let truth b = if b then Sym "t" else Nil

let rec list_of = function
  | Nil -> []
  | Cons (a, d) -> a :: list_of d
  | v -> err "improper list ending in %s" (show_exp v)

let rec to_list = function [] -> Nil | x :: xs -> Cons (x, to_list xs)

let rec bind params args env =
  match (params, args) with
  | Nil, [] -> env
  | Sym p, rest -> (p, to_list rest) :: env
  | Cons (Sym p, ps), a :: rest -> bind ps rest ((p, a) :: env)
  | _ -> err "arity mismatch"

let rec eval env = function
  | (Float _ | Nil | Closure _ | Prim _ | Macro _) as v -> v
  | Sym s -> lookup env s
  | Cons (Sym "quote", Cons (x, Nil)) -> x
  | Cons (Sym "if", Cons (c, Cons (t, Cons (e, Nil)))) ->
      eval env (if truthy (eval env c) then t else e)
  | Cons (Sym "lambda", Cons (params, Cons (body, Nil))) ->
      Closure { params; body; env }
  | Cons (Sym "define", Cons (Sym name, Cons (e, Nil))) ->
      define name (eval env e)
  | Cons (Sym "defmacro", Cons (Sym name, Cons (params, Cons (body, Nil)))) ->
      define name (Macro { params; body; env = [] })
  | Cons (f, args) -> (
      match eval env f with
      | Macro m -> eval env (apply_closure m (list_of args))
      | fv -> apply fv (List.map (eval env) (list_of args)))

and apply_closure c args = eval (bind c.params args c.env) c.body

and apply f args =
  match f with
  | Closure c -> apply_closure c args
  | Prim (_, p) -> p args
  | v -> err "not a function: %s" (show_exp v)

let num = function
  | Float f -> f
  | v -> err "expected a number, got %s" (show_exp v)

let init_prims () =
  let bin name f =
    Hashtbl.replace globals name
      (Prim
         ( name,
           function [ a; b ] -> f a b | _ -> err "%s takes 2 arguments" name ))
  in
  let arith name f = bin name (fun a b -> f (num a) (num b)) in
  bin "cons" (fun a d -> Cons (a, d));
  bin "uncons" (fun xs f ->
      match xs with
      | Nil -> Nil
      | Cons (h, t) -> apply f [ h; t ]
      | v -> err "uncons expects a list, got %s" (show_exp v));
  bin "eq?" (fun a b -> truth (equal_exp a b));
  arith "-" (fun a b -> Float (a -. b));
  arith "*" (fun a b -> Float (a *. b));
  arith "<" (fun a b -> truth (a < b))

let run_string src =
  Parser.program Lexer.token (Lexing.from_string src)
  |> List.to_seq
  |> Seq.map (eval [])
  |> List.of_seq

let load_file path =
  run_string (In_channel.with_open_bin path In_channel.input_all)
