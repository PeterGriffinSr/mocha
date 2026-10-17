(lang dune 3.22)

(name mocha)

(generate_opam_files true)

(source
 (github PeterGriffinSr/mocha))

(authors "Author Name <author@example.com>")

(maintainers "Maintainer Name <maintainer@example.com>")

(license GPL-2.0-or-later)

(documentation https://mocha.<domain>/docs)

(package
 (name mocha)
 (synopsis "Yet another lisp dialect.")
 (description "Yet another lisp dialect.")
 (depends ocaml dune ppx_deriving ppx_sexp_conv sexplib)
 (tags
  (lisp compiler interpreter)))

; See the complete stanza docs at https://dune.readthedocs.io/en/stable/reference/dune-project/index.html
