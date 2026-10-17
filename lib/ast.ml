open Sexplib.Conv

type exp =
  | Nil
  | Float of float
  | Sym of string
  | Cons of exp * exp
  | Closure of (closure[@printer fun fmt _ -> Format.fprintf fmt "<closure>"])
  | Prim of string * ((exp list -> exp)[@equal fun _ _ -> false])
  | Macro of (closure[@printer fun fmt _ -> Format.fprintf fmt "<macro>"])

and closure = { params : exp; body : exp; env : (string * exp) list }
[@@deriving show, sexp, eq]
