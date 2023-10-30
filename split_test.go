package repl

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	for n, tc := range map[string]struct {
		in  string
		out []string
		err bool
	}{
		"simple":         {in: "foo bar baz", out: []string{"foo", "bar", "baz"}},
		"longSpaces":     {in: "foo  bar  baz", out: []string{"foo", "bar", "baz"}},
		"slashSpaces":    {in: `foo bar\ baz quux`, out: []string{"foo", "bar baz", "quux"}},
		"slash":          {in: `foo\\bar baz`, out: []string{`foo\bar`, "baz"}},
		"endingSlash":    {in: `foo bar\\`, out: []string{"foo", `bar\`}},
		"endingBadSlash": {in: `foo bar\`, err: true},
		"endingSpace":    {in: "foo bar ", out: []string{"foo", "bar"}},
		"hex":            {in: `foo\x{41}bar`, out: []string{"fooAbar"}},
		"hexOverflow":    {in: `\x{1ffffffff}`, err: true},
		"doubleHex":      {in: `\x{41}\x{42}`, out: []string{"AB"}},
		"hexLetters":     {in: `\x{7a}\x{7A}`, out: []string{"zz"}},
	} {
		t.Run(n, func(t *testing.T) {
			out, err := split(tc.in)
			t.Logf("actual: %s", strings.Join(out, "|"))
			t.Logf("expected: %s", strings.Join(tc.out, "|"))
			if (err != nil) != tc.err {
				if tc.err {
					t.Fatal("expected error but got none")
				} else {
					t.Fatalf("expected no error but got: %s", err)
				}
			}
			if len(out) != len(tc.out) {
				t.Fatalf("expected output len %d but got %d", len(tc.out), len(out))
			}
			for i := range out {
				if out[i] != tc.out[i] {
					t.Fatalf("expected %q but got %q", tc.out[i], out[i])
				}
			}
		})
	}
}
