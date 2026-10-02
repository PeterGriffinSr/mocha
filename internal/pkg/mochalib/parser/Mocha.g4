grammar Mocha;

program: NEWLINE* (topLevelDecl (NEWLINE+ | EOF))* EOF;

topLevelDecl: functionDecl | structDecl;

functionDecl:
	FUNCTION IDENT LPAREN parameterList? RPAREN type? block;

parameterList:
	NEWLINE* parameter (COMMA NEWLINE* parameter)* COMMA? NEWLINE*;

parameter: IDENT type;

structDecl: STRUCT IDENT LBRACE NEWLINE* fieldList? RBRACE;

fieldList:
	fieldDecl (COMMA NEWLINE* fieldDecl)* COMMA? NEWLINE*;

fieldDecl: IDENT COLON type;

variableDecl: LET IDENT COLON type ASSIGN NEWLINE* expression;

block:
	LBRACE NEWLINE* (statement (NEWLINE+ statement)* NEWLINE*)? RBRACE;

statement: variableDecl | returnStmt | assignment | expression;

returnStmt: RETURN expression?;

assignment: assignTarget ASSIGN NEWLINE* expression;

assignTarget: IDENT (DOT IDENT)*;

type: primitiveType | arrayType | IDENT;

primitiveType:
	I8
	| I16
	| I32
	| I64
	| F16
	| F32
	| F64
	| F128
	| STRING_TYPE;

arrayType: ARRAY LT type GT;

expression:
	expression DOT IDENT
	| op = (NOT | MINUS) expression
	| expression op = (STAR | SLASH | PERCENT) NEWLINE* expression
	| expression op = (PLUS | MINUS) NEWLINE* expression
	| expression op = (LT | GT | LE | GE) NEWLINE* expression
	| expression op = (EQ | NEQ) NEWLINE* expression
	| expression AND NEWLINE* expression
	| expression OR NEWLINE* expression
	| primary;

primary:
	literal
	| functionCall
	| structLiteral
	| arrayLiteral
	| IDENT
	| LPAREN NEWLINE* expression NEWLINE* RPAREN;

functionCall: IDENT LPAREN argumentList? RPAREN;

argumentList:
	NEWLINE* expression (COMMA NEWLINE* expression)* COMMA? NEWLINE*;

structLiteral: IDENT LBRACE NEWLINE* fieldInitList? RBRACE;

fieldInitList:
	fieldInit (COMMA NEWLINE* fieldInit)* COMMA? NEWLINE*;

fieldInit: IDENT COLON expression;

arrayLiteral:
	LBRACK NEWLINE* (
		expression (COMMA NEWLINE* expression)* COMMA? NEWLINE*
	)? RBRACK;

literal: INT_LITERAL | FLOAT_LITERAL | STRING_LITERAL;

FUNCTION: 'function';
STRUCT: 'struct';
LET: 'let';
RETURN: 'return';

I8: 'i8';
I16: 'i16';
I32: 'i32';
I64: 'i64';
F16: 'f16';
F32: 'f32';
F64: 'f64';
F128: 'f128';
STRING_TYPE: 'string';
ARRAY: 'array';

AND: '&&';
OR: '||';
EQ: '==';
NEQ: '!=';
LE: '<=';
GE: '>=';
LT: '<';
GT: '>';
NOT: '!';
ASSIGN: '=';
PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
PERCENT: '%';

COMMA: ',';
COLON: ':';
DOT: '.';
LPAREN: '(';
RPAREN: ')';
LBRACE: '{';
RBRACE: '}';
LBRACK: '[';
RBRACK: ']';

FLOAT_LITERAL: DIGIT+ '.' DIGIT+ EXPONENT? | DIGIT+ EXPONENT;

INT_LITERAL: '0' [xX] HEX_DIGIT+ | '0' [bB] [01]+ | DIGIT+;

STRING_LITERAL: '"' ( ESCAPE | ~["\\\r\n])* '"';

IDENT: [A-Za-z_] [A-Za-z0-9_]*;

NEWLINE: '\r'? '\n';

COMMENT: '#' ~[\r\n]* -> skip;
WS: [ \t]+ -> skip;

fragment DIGIT: [0-9];
fragment HEX_DIGIT: [0-9a-fA-F];
fragment EXPONENT: [eE] [+\-]? DIGIT+;
fragment ESCAPE: '\\' .;
