package mochalib

import (
	"fmt"
	"sort"
)

type Span struct {
	Start int
	End   int
}

func (span Span) String() string {
	return fmt.Sprintf("%d:%d", span.Start, span.End)
}

func (span Span) Len() int {
	return max(span.End-span.Start, 0)
}

func (span Span) Empty() bool {
	return span.Start == span.End
}

type Position struct {
	Line   int
	Column int
}

func (position Position) String() string {
	return fmt.Sprintf("%d:%d", position.Line, position.Column)
}

type Source struct {
	runes      []rune
	lineStarts []int
}

func NewSource(text string) Source {
	runes := []rune(text)

	lineStarts := []int{0}

	for i, ch := range runes {
		if ch == '\n' {
			lineStarts = append(lineStarts, i+1)
		}
	}

	return Source{
		runes:      runes,
		lineStarts: lineStarts,
	}
}

func (source Source) Runes() []rune {
	return source.runes
}

func (source Source) Len() int {
	return len(source.runes)
}

func (source Source) PositionOf(offset int) Position {
	offset = clamp(offset, 0, len(source.runes))

	lineIndex := sort.Search(len(source.lineStarts), func(i int) bool {
		return source.lineStarts[i] > offset
	}) - 1

	lineStart := source.lineStarts[lineIndex]

	return Position{
		Line:   lineIndex + 1,
		Column: offset - lineStart + 1,
	}
}

func (source Source) Line(lineNum int) []rune {
	if lineNum < 1 || lineNum > len(source.lineStarts) {
		return nil
	}

	start := source.lineStarts[lineNum-1]

	end := len(source.runes)
	if lineNum < len(source.lineStarts) {
		end = source.lineStarts[lineNum] - 1
	}

	return source.runes[start:end]
}

func PositionOf(source []rune, offset int) Position {
	return NewSource(string(source)).PositionOf(offset)
}

func clamp(value, lower, upper int) int {
	return min(max(value, lower), upper)
}
