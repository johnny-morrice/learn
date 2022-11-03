package parser

import (
	"reflect"
	"testing"
)

func TestParserCombinators(t *testing.T) {
	type testCase struct {
		input    ParseContext
		comb     ParseCombinator
		expected ParseContext
	}

	testCases := map[string]testCase{
		"EOF WhenMatch": {
			comb: EOF(),
		},
		"EOF WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "  ",
			},
			comb: EOF(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "  ",
			},
		},

		"TextEq WhenMatch": {
			input: ParseContext{
				RemainingInput: "hello there",
			},
			comb: TextEq("hello"),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: " there",
			},
		},
		"TextEq WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "hello there",
			},
			comb: TextEq("hi"),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "hello there",
			},
		},

		"Number WhenMatch": {
			input: ParseContext{
				RemainingInput: "123 apples",
			},
			comb: Number(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: " apples",
			},
		},
		"Number WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "one two three apples",
			},
			comb: Number(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "one two three apples",
			},
		},

		"VarName WhenMatch": {
			input: ParseContext{
				RemainingInput: "foo bar",
			},
			comb: VarName(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: " bar",
			},
		},
		"VarName WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "123 bar",
			},
			comb: VarName(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "123 bar",
			},
		},

		"Whitespace WhenMatch": {
			input: ParseContext{
				RemainingInput: "    foo",
			},
			comb: Whitespace(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: "foo",
			},
		},
		"Whitespace WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "foo",
			},
			comb: Whitespace(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "foo",
			},
		},

		"OpStmt WhenMatch": {
			input: ParseContext{
				RemainingInput: "push foo 123",
			},
			comb: OpStmt(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: "",
			},
		},
		"OpStmt WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "var foo 123",
			},
			comb: OpStmt(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "var foo 123",
			},
		},

		"VarStmt WhenMatch": {
			input: ParseContext{
				RemainingInput: "var foo bar",
			},
			comb: VarStmt(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: "",
			},
		},
		"VarStmt WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "push foo 123",
			},
			comb: VarStmt(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "push foo 123",
			},
		},

		"LabelStmt WhenMatch": {
			input: ParseContext{
				RemainingInput: "foo:",
			},
			comb: LabelStmt(),
			expected: ParseContext{
				Failed:         false,
				RemainingInput: "",
			},
		},
		"LabelStmt WhenNotMatch": {
			input: ParseContext{
				RemainingInput: "foo",
			},
			comb: LabelStmt(),
			expected: ParseContext{
				Failed:         true,
				RemainingInput: "foo",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := tc.comb(tc.input)
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected: %v\n\nbut was: %v", tc.expected, actual)
			}
		})
	}
}
