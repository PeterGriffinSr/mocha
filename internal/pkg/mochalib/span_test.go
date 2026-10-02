package mochalib

import "testing"

func TestSpanString(t *testing.T) {
	tests := []struct {
		name string
		span Span
		want string
	}{
		{
			name: "normal",
			span: Span{Start: 2, End: 7},
			want: "2:7",
		},
		{
			name: "empty",
			span: Span{Start: 4, End: 4},
			want: "4:4",
		},
		{
			name: "zero",
			span: Span{Start: 0, End: 0},
			want: "0:0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.String(); got != tt.want {
				t.Fatalf("Span.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSpanLen(t *testing.T) {
	tests := []struct {
		name string
		span Span
		want int
	}{
		{
			name: "normal",
			span: Span{Start: 2, End: 7},
			want: 5,
		},
		{
			name: "empty",
			span: Span{Start: 4, End: 4},
			want: 0,
		},
		{
			name: "reversed",
			span: Span{Start: 7, End: 2},
			want: 0,
		},
		{
			name: "zero",
			span: Span{Start: 0, End: 0},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.Len(); got != tt.want {
				t.Fatalf("Span.Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpanEmpty(t *testing.T) {
	tests := []struct {
		name string
		span Span
		want bool
	}{
		{
			name: "empty",
			span: Span{Start: 4, End: 4},
			want: true,
		},
		{
			name: "normal",
			span: Span{Start: 2, End: 7},
			want: false,
		},
		{
			name: "reversed",
			span: Span{Start: 7, End: 2},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.Empty(); got != tt.want {
				t.Fatalf("Span.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionString(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     string
	}{
		{
			name:     "normal",
			position: Position{Line: 3, Column: 5},
			want:     "3:5",
		},
		{
			name:     "first",
			position: Position{Line: 1, Column: 1},
			want:     "1:1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.String(); got != tt.want {
				t.Fatalf("Position.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewSource(t *testing.T) {
	source := NewSource("hello\nworld")

	if got := source.Len(); got != 11 {
		t.Fatalf("Source.Len() = %d, want 11", got)
	}

	if got := string(source.Runes()); got != "hello\nworld" {
		t.Fatalf("Source.Runes() = %q, want %q", got, "hello\nworld")
	}
}

func TestSourceLenUsesRunes(t *testing.T) {
	source := NewSource("é🙂")

	if got := source.Len(); got != 2 {
		t.Fatalf("Source.Len() = %d, want 2", got)
	}

	if got := len(source.Runes()); got != 2 {
		t.Fatalf("len(Source.Runes()) = %d, want 2", got)
	}
}

func TestSourcePositionOf(t *testing.T) {
	source := NewSource("hello\nworld")

	tests := []struct {
		name   string
		offset int
		want   Position
	}{
		{
			name:   "start",
			offset: 0,
			want:   Position{Line: 1, Column: 1},
		},
		{
			name:   "middle of first line",
			offset: 2,
			want:   Position{Line: 1, Column: 3},
		},
		{
			name:   "newline",
			offset: 5,
			want:   Position{Line: 1, Column: 6},
		},
		{
			name:   "start of second line",
			offset: 6,
			want:   Position{Line: 2, Column: 1},
		},
		{
			name:   "middle of second line",
			offset: 8,
			want:   Position{Line: 2, Column: 3},
		},
		{
			name:   "end of source",
			offset: 11,
			want:   Position{Line: 2, Column: 6},
		},
		{
			name:   "negative offset",
			offset: -10,
			want:   Position{Line: 1, Column: 1},
		},
		{
			name:   "offset past end",
			offset: 100,
			want:   Position{Line: 2, Column: 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := source.PositionOf(tt.offset); got != tt.want {
				t.Fatalf("Source.PositionOf(%d) = %+v, want %+v", tt.offset, got, tt.want)
			}
		})
	}
}

func TestSourcePositionOfUnicode(t *testing.T) {
	source := NewSource("aé🙂b")

	tests := []struct {
		offset int
		want   Position
	}{
		{offset: 0, want: Position{Line: 1, Column: 1}},
		{offset: 1, want: Position{Line: 1, Column: 2}},
		{offset: 2, want: Position{Line: 1, Column: 3}},
		{offset: 3, want: Position{Line: 1, Column: 4}},
		{offset: 4, want: Position{Line: 1, Column: 5}},
	}

	for _, tt := range tests {
		t.Run(
			source.PositionOf(tt.offset).String(),
			func(t *testing.T) {
				if got := source.PositionOf(tt.offset); got != tt.want {
					t.Fatalf(
						"Source.PositionOf(%d) = %+v, want %+v",
						tt.offset,
						got,
						tt.want,
					)
				}
			},
		)
	}
}

func TestSourceLine(t *testing.T) {
	source := NewSource("first\nsecond\n")

	tests := []struct {
		name    string
		lineNum int
		want    string
	}{
		{
			name:    "first line",
			lineNum: 1,
			want:    "first",
		},
		{
			name:    "second line",
			lineNum: 2,
			want:    "second",
		},
		{
			name:    "trailing empty line",
			lineNum: 3,
			want:    "",
		},
		{
			name:    "line before first",
			lineNum: 0,
			want:    "",
		},
		{
			name:    "line after end",
			lineNum: 4,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(source.Line(tt.lineNum))

			if got != tt.want {
				t.Fatalf("Source.Line(%d) = %q, want %q", tt.lineNum, got, tt.want)
			}
		})
	}
}

func TestSourceLineUnicode(t *testing.T) {
	source := NewSource("héllo\n世界")

	if got := string(source.Line(1)); got != "héllo" {
		t.Fatalf("Source.Line(1) = %q, want %q", got, "héllo")
	}

	if got := string(source.Line(2)); got != "世界" {
		t.Fatalf("Source.Line(2) = %q, want %q", got, "世界")
	}
}

func TestPositionOfFunction(t *testing.T) {
	source := []rune("hello\nworld")

	want := Position{Line: 2, Column: 3}

	if got := PositionOf(source, 8); got != want {
		t.Fatalf("PositionOf(%d) = %+v, want %+v", 8, got, want)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		lower int
		upper int
		want  int
	}{
		{
			name:  "below lower",
			value: -10,
			lower: 0,
			upper: 10,
			want:  0,
		},
		{
			name:  "inside range",
			value: 5,
			lower: 0,
			upper: 10,
			want:  5,
		},
		{
			name:  "above upper",
			value: 20,
			lower: 0,
			upper: 10,
			want:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clamp(tt.value, tt.lower, tt.upper); got != tt.want {
				t.Fatalf("clamp(%d, %d, %d) = %d, want %d",
					tt.value, tt.lower, tt.upper, got, tt.want)
			}
		})
	}
}
