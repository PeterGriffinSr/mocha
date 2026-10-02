package mochalib

import (
	"mocha/internal/pkg/mochalib/parser"
	"sort"

	"github.com/antlr4-go/antlr/v4"
)

type syntaxErrorListener struct {
	*antlr.DefaultErrorListener
	source      Source
	diagnostics *Diagnostics
}

func (listener *syntaxErrorListener) SyntaxError(
	_ antlr.Recognizer,
	offendingSymbol any,
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	span := listener.spanFor(offendingSymbol, line, column)

	listener.diagnostics.AddError(span, "%s", msg)
}

func (listener *syntaxErrorListener) spanFor(offendingSymbol any, line, column int) Span {
	if token, ok := offendingSymbol.(antlr.Token); ok && token != nil {
		return SpanOfToken(token)
	}

	offset := listener.source.OffsetOf(line, column)

	return Span{
		Start: offset,
		End:   offset + 1,
	}
}

func SpanOfToken(token antlr.Token) Span {
	start, stop := token.GetStart(), token.GetStop()

	if token.GetTokenType() == antlr.TokenEOF || stop < start {
		return Span{
			Start: start,
			End:   start,
		}
	}

	return Span{
		Start: start,
		End:   stop + 1,
	}
}

func SpanOfContext(ctx antlr.ParserRuleContext) Span {
	first, last := ctx.GetStart(), ctx.GetStop()

	if first == nil {
		return Span{}
	}

	if last == nil {
		return SpanOfToken(first)
	}

	return Span{Start: SpanOfToken(first).Start, End: SpanOfToken(last).End}
}

type ParseResult struct {
	Tree        parser.IProgramContext
	Source      Source
	Diagnostics Diagnostics
}

func Parse(text string) ParseResult {
	source := NewSource(text)

	var diagnostics Diagnostics

	listener := &syntaxErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		source:               source,
		diagnostics:          &diagnostics,
	}

	lexer := parser.NewMochaLexer(antlr.NewInputStream(text))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(listener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewMochaParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(listener)

	tree := p.Program()

	sort.SliceStable(diagnostics, func(i, j int) bool {
		return diagnostics[i].Span.Start < diagnostics[j].Span.Start
	})

	return ParseResult{Tree: tree, Source: source, Diagnostics: diagnostics}
}
