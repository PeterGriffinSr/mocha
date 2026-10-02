package mochalib

import "testing"

func TestSeverityString(t *testing.T) {
	tests := []struct {
		name     string
		severity Severity
		want     string
	}{
		{
			name:     "error",
			severity: SeverityError,
			want:     "Error",
		},
		{
			name:     "warning",
			severity: SeverityWarning,
			want:     "Warning",
		},
		{
			name:     "unknown",
			severity: Severity(255),
			want:     "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.want {
				t.Fatalf("Severity.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	diagnostic := Errorf(Span{Start: 2, End: 5}, "unexpected %s", "token")

	if diagnostic.Severity != SeverityError {
		t.Fatalf("Severity = %v, want %v", diagnostic.Severity, SeverityError)
	}

	if diagnostic.Span != (Span{Start: 2, End: 5}) {
		t.Fatalf("Span = %+v, want %+v", diagnostic.Span, Span{Start: 2, End: 5})
	}

	if diagnostic.Message != "unexpected token" {
		t.Fatalf("Message = %q, want %q", diagnostic.Message, "unexpected token")
	}
}

func TestWarningf(t *testing.T) {
	diagnostic := Warningf(Span{Start: 1, End: 4}, "unused %s", "variable")

	if diagnostic.Severity != SeverityWarning {
		t.Fatalf("Severity = %v, want %v", diagnostic.Severity, SeverityWarning)
	}

	if diagnostic.Span != (Span{Start: 1, End: 4}) {
		t.Fatalf("Span = %+v, want %+v", diagnostic.Span, Span{Start: 1, End: 4})
	}

	if diagnostic.Message != "unused variable" {
		t.Fatalf("Message = %q, want %q", diagnostic.Message, "unused variable")
	}
}

func TestDiagnosticsAddError(t *testing.T) {
	var diagnostics Diagnostics

	diagnostics.AddError(Span{Start: 0, End: 3}, "invalid %s", "syntax")

	if len(diagnostics) != 1 {
		t.Fatalf("len(diagnostics) = %d, want 1", len(diagnostics))
	}

	diagnostic := diagnostics[0]

	if diagnostic.Severity != SeverityError {
		t.Fatalf("Severity = %v, want %v", diagnostic.Severity, SeverityError)
	}

	if diagnostic.Span != (Span{Start: 0, End: 3}) {
		t.Fatalf("Span = %+v, want %+v", diagnostic.Span, Span{Start: 0, End: 3})
	}

	if diagnostic.Message != "invalid syntax" {
		t.Fatalf("Message = %q, want %q", diagnostic.Message, "invalid syntax")
	}
}

func TestDiagnosticsAddWarning(t *testing.T) {
	var diagnostics Diagnostics

	diagnostics.AddWarning(Span{Start: 4, End: 7}, "unused %s", "import")

	if len(diagnostics) != 1 {
		t.Fatalf("len(diagnostics) = %d, want 1", len(diagnostics))
	}

	diagnostic := diagnostics[0]

	if diagnostic.Severity != SeverityWarning {
		t.Fatalf("Severity = %v, want %v", diagnostic.Severity, SeverityWarning)
	}

	if diagnostic.Span != (Span{Start: 4, End: 7}) {
		t.Fatalf("Span = %+v, want %+v", diagnostic.Span, Span{Start: 4, End: 7})
	}

	if diagnostic.Message != "unused import" {
		t.Fatalf("Message = %q, want %q", diagnostic.Message, "unused import")
	}
}

func TestDiagnosticsEmpty(t *testing.T) {
	var diagnostics Diagnostics

	if !diagnostics.Empty() {
		t.Fatal("Diagnostics.Empty() = false, want true")
	}

	diagnostics.AddError(Span{}, "error")

	if diagnostics.Empty() {
		t.Fatal("Diagnostics.Empty() = true, want false")
	}
}

func TestDiagnosticsHasErrors(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics Diagnostics
		want        bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name: "only warnings",
			diagnostics: Diagnostics{
				Warningf(Span{}, "warning"),
			},
			want: false,
		},
		{
			name: "one error",
			diagnostics: Diagnostics{
				Errorf(Span{}, "error"),
			},
			want: true,
		},
		{
			name: "error after warning",
			diagnostics: Diagnostics{
				Warningf(Span{}, "warning"),
				Errorf(Span{}, "error"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diagnostics.HasErrors(); got != tt.want {
				t.Fatalf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticsHasWarnings(t *testing.T) {
	tests := []struct {
		name        string
		diagnostics Diagnostics
		want        bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name: "only errors",
			diagnostics: Diagnostics{
				Errorf(Span{}, "error"),
			},
			want: false,
		},
		{
			name: "one warning",
			diagnostics: Diagnostics{
				Warningf(Span{}, "warning"),
			},
			want: true,
		},
		{
			name: "warning after error",
			diagnostics: Diagnostics{
				Errorf(Span{}, "error"),
				Warningf(Span{}, "warning"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diagnostics.HasWarnings(); got != tt.want {
				t.Fatalf("HasWarnings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticsFormat(t *testing.T) {
	diagnostics := Diagnostics{
		Errorf(Span{Start: 6, End: 11}, "unexpected identifier"),
	}

	got := diagnostics.Format("main.mocha", "hello world")

	want := "main.mocha:1:7: Error: unexpected identifier\n" +
		"  |\n" +
		"1 | hello world\n" +
		"  |       ^^^^^"

	if got != want {
		t.Fatalf("Diagnostics.Format() =\n%q\nwant:\n%q", got, want)
	}
}

func TestDiagnosticsFormatMultiple(t *testing.T) {
	diagnostics := Diagnostics{
		Errorf(Span{Start: 0, End: 1}, "first error"),
		Warningf(Span{Start: 2, End: 3}, "second warning"),
	}

	got := diagnostics.Format("main.mocha", "abc")

	want := "main.mocha:1:1: Error: first error\n" +
		"  |\n" +
		"1 | abc\n" +
		"  | ^\n" +
		"main.mocha:1:3: Warning: second warning\n" +
		"  |\n" +
		"1 | abc\n" +
		"  |   ^"

	if got != want {
		t.Fatalf("Diagnostics.Format() =\n%q\nwant:\n%q", got, want)
	}
}

func TestDiagnosticsFormatSource(t *testing.T) {
	source := NewSource("first line\nsecond line")

	diagnostics := Diagnostics{
		Errorf(Span{Start: 11, End: 17}, "bad expression"),
	}

	got := diagnostics.FormatSource("test.mocha", source)

	want := "test.mocha:2:1: Error: bad expression\n" +
		"  |\n" +
		"2 | second line\n" +
		"  | ^^^^^^"

	if got != want {
		t.Fatalf("Diagnostics.FormatSource() =\n%q\nwant:\n%q", got, want)
	}
}

func TestFormatDiagnosticEmptySpan(t *testing.T) {
	source := NewSource("hello")

	diagnostic := Errorf(Span{Start: 2, End: 2}, "expected expression")

	got := FormatDiagnostic("main.mocha", source, diagnostic)

	want := "main.mocha:1:3: Error: expected expression\n" +
		"  |\n" +
		"1 | hello\n" +
		"  |   ^"

	if got != want {
		t.Fatalf("FormatDiagnostic() =\n%q\nwant:\n%q", got, want)
	}
}

func TestFormatDiagnosticClampsSpan(t *testing.T) {
	source := NewSource("hello")

	tests := []struct {
		name string
		span Span
		want string
	}{
		{
			name: "start before line",
			span: Span{-10, 2},
			want: "main.mocha:1:1: Error: test\n" +
				"  |\n" +
				"1 | hello\n" +
				"  | ^^",
		},
		{
			name: "end past line",
			span: Span{2, 100},
			want: "main.mocha:1:3: Error: test\n" +
				"  |\n" +
				"1 | hello\n" +
				"  |   ^^^",
		},
		{
			name: "completely past line",
			span: Span{100, 200},
			want: "main.mocha:1:6: Error: test\n" +
				"  |\n" +
				"1 | hello\n" +
				"  |      ^",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDiagnostic(
				"main.mocha",
				source,
				Errorf(tt.span, "test"),
			)

			if got != tt.want {
				t.Fatalf("FormatDiagnostic() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		name string
		line []rune
		want string
	}{
		{
			name: "single tab",
			line: []rune("\tfoo"),
			want: "    foo",
		},
		{
			name: "tab after one character",
			line: []rune("a\tfoo"),
			want: "a   foo",
		},
		{
			name: "tab after four characters",
			line: []rune("abcd\tfoo"),
			want: "abcd    foo",
		},
		{
			name: "multiple tabs",
			line: []rune("\t\tfoo"),
			want: "        foo",
		},
		{
			name: "no tabs",
			line: []rune("hello"),
			want: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandTabs(tt.line); got != tt.want {
				t.Fatalf("expandTabs() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		name string
		text []rune
		want int
	}{
		{
			name: "empty",
			text: []rune(""),
			want: 0,
		},
		{
			name: "plain text",
			text: []rune("hello"),
			want: 5,
		},
		{
			name: "single tab",
			text: []rune("\t"),
			want: 4,
		},
		{
			name: "tab after one character",
			text: []rune("a\t"),
			want: 4,
		},
		{
			name: "tab after two characters",
			text: []rune("ab\t"),
			want: 4,
		},
		{
			name: "tab after three characters",
			text: []rune("abc\t"),
			want: 4,
		},
		{
			name: "tab after four characters",
			text: []rune("abcd\t"),
			want: 8,
		},
		{
			name: "unicode",
			text: []rune("hé"),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := displayWidth(tt.text); got != tt.want {
				t.Fatalf("displayWidth() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFormatDiagnosticTabs(t *testing.T) {
	source := NewSource("\tfoo := 42")

	diagnostic := Errorf(
		Span{Start: 1, End: 4},
		"unexpected token",
	)

	got := FormatDiagnostic("main.mocha", source, diagnostic)

	want := "main.mocha:1:2: Error: unexpected token\n" +
		"  |\n" +
		"1 |     foo := 42\n" +
		"  |     ^^^"

	if got != want {
		t.Fatalf("FormatDiagnostic() =\n%q\nwant:\n%q", got, want)
	}
}

func TestFormatDiagnosticUnicode(t *testing.T) {
	source := NewSource("héllo world")

	diagnostic := Errorf(
		Span{Start: 2, End: 7},
		"invalid identifier",
	)

	got := FormatDiagnostic("main.mocha", source, diagnostic)

	want := "main.mocha:1:3: Error: invalid identifier\n" +
		"  |\n" +
		"1 | héllo world\n" +
		"  |   ^^^^^"

	if got != want {
		t.Fatalf("FormatDiagnostic() =\n%q\nwant:\n%q", got, want)
	}
}

func TestFormatDiagnosticMultilineSpan(t *testing.T) {
	source := NewSource("first\nsecond\nthird")

	diagnostic := Errorf(
		Span{Start: 2, End: 10},
		"invalid construct",
	)

	got := FormatDiagnostic("main.mocha", source, diagnostic)

	want := "main.mocha:1:3: Error: invalid construct\n" +
		"  |\n" +
		"1 | first\n" +
		"  |   ^^^"

	if got != want {
		t.Fatalf("FormatDiagnostic() =\n%q\nwant:\n%q", got, want)
	}
}
