type exp =
  | Float of float
  | Sym of string
  | Data of string * exp list
  | Closure of
      (closure
      [@printer fun fmt _ -> Format.fprintf fmt "<closure>"]
      [@equal fun a b -> a == b])
  | Prim of string * ((exp list -> exp)[@equal fun _ _ -> true])
  | Macro of
      (closure
      [@printer fun fmt _ -> Format.fprintf fmt "<macro>"]
      [@equal fun a b -> a == b])

and closure = { params : exp; body : exp; env : (string * exp) list }
[@@deriving show, eq]

let nil = Data ("Nil", [])
let cons a d = Data ("Cons", [ a; d ])
