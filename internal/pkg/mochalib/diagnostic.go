package mochalib

import (
	"fmt"
	"strings"
)

const tabWidth = 4

type Severity uint8

const (
	SeverityError Severity = iota
	SeverityWarning
)

func (severity Severity) String() string {
	switch severity {
	case SeverityError:
		return "Error"

	case SeverityWarning:
		return "Warning"

	default:
		return "Unknown"
	}
}

type Diagnostic struct {
	Severity Severity
	Span     Span
	Message  string
}

func Errorf(span Span, format string, args ...any) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Span:     span,
		Message:  fmt.Sprintf(format, args...),
	}
}

func Warningf(span Span, format string, args ...any) Diagnostic {
	return Diagnostic{
		Severity: SeverityWarning,
		Span:     span,
		Message:  fmt.Sprintf(format, args...),
	}
}

type Diagnostics []Diagnostic

func (diagnostics *Diagnostics) AddError(span Span, format string, args ...any) {
	*diagnostics = append(*diagnostics, Errorf(span, format, args...))
}

func (diagnostics *Diagnostics) AddWarning(span Span, format string, args ...any) {
	*diagnostics = append(*diagnostics, Warningf(span, format, args...))
}

func (diagnostics Diagnostics) Empty() bool {
	return len(diagnostics) == 0
}

func (diagnostics Diagnostics) HasErrors() bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}

	return false
}

func (diagnostics Diagnostics) HasWarnings() bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == SeverityWarning {
			return true
		}
	}

	return false
}

func (diagnostics Diagnostics) Format(filename string, source string) string {
	return diagnostics.FormatSource(filename, NewSource(source))
}

func (diagnostics Diagnostics) FormatSource(filename string, source Source) string {
	var b strings.Builder

	for i, diagnostic := range diagnostics {
		if i > 0 {
			b.WriteByte('\n')
		}

		b.WriteString(FormatDiagnostic(filename, source, diagnostic))
	}

	return b.String()
}

func FormatDiagnostic(filename string, source Source, diagnostic Diagnostic) string {
	position := source.PositionOf(diagnostic.Span.Start)
	line := source.Line(position.Line)

	lineStart := source.lineStarts[position.Line-1]
	lineEnd := lineStart + len(line)

	start := clamp(diagnostic.Span.Start, lineStart, lineEnd)
	end := clamp(diagnostic.Span.End, start, lineEnd)

	startColumn := start - lineStart
	endColumn := end - lineStart

	prefixWidth := displayWidth(line[:startColumn])
	underlineWidth := displayWidth(line[startColumn:endColumn])

	if underlineWidth == 0 {
		underlineWidth = 1
	}

	displayedLine := expandTabs(line)

	lineNumStr := fmt.Sprintf("%d", position.Line)
	gutter := strings.Repeat(" ", len(lineNumStr))

	var b strings.Builder

	fmt.Fprintf(
		&b,
		"%s:%s: %s: %s\n",
		filename,
		position,
		diagnostic.Severity,
		diagnostic.Message,
	)

	fmt.Fprintf(&b, "%s |\n", gutter)
	fmt.Fprintf(&b, "%s | %s\n", lineNumStr, displayedLine)
	fmt.Fprintf(
		&b,
		"%s | %s%s",
		gutter,
		strings.Repeat(" ", prefixWidth),
		strings.Repeat("^", underlineWidth),
	)

	return b.String()
}

func expandTabs(line []rune) string {
	var b strings.Builder

	column := 0

	for _, ch := range line {
		if ch != '\t' {
			b.WriteRune(ch)
			column++
			continue
		}

		spaces := tabWidth - column%tabWidth
		b.WriteString(strings.Repeat(" ", spaces))
		column += spaces
	}

	return b.String()
}

func displayWidth(text []rune) int {
	width := 0
	column := 0

	for _, ch := range text {
		if ch == '\t' {
			spaces := tabWidth - column%tabWidth
			width += spaces
			column += spaces
			continue
		}

		width++
		column++
	}

	return width
}
