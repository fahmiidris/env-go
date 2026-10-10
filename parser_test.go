package env

import (
	"strings"
	"testing"
)

func TestParseSimple(t *testing.T) {
	input := "KEY=VALUE"
	expected := map[string]string{"KEY": "VALUE"}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseCommentsAndEmptyLines(t *testing.T) {
	input := `
		# This is a comment
		KEY1=VALUE1

		KEY2=VALUE2 # inline comment
	`

	expected := map[string]string{
		"KEY1": "VALUE1",
		"KEY2": "VALUE2",
	}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseExportPrefix(t *testing.T) {
	input := "export KEY=VALUE"
	expected := map[string]string{"KEY": "VALUE"}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseSingleQuotes(t *testing.T) {
	input := "KEY='VALUE with # spaces'"
	expected := map[string]string{"KEY": "VALUE with # spaces"}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseDoubleQuotesWithEscapes(t *testing.T) {
	input := `KEY="VALUE with\nnewlines and \"quotes\""`
	expected := map[string]string{"KEY": "VALUE with\nnewlines and \"quotes\""}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseUnquotedWithInlineComment(t *testing.T) {
	input := "KEY=VALUE # comment"
	expected := map[string]string{"KEY": "VALUE"}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func TestParseVariableExpansion(t *testing.T) {
	input := `
		DB_USER=postgres
		DB_PASS=postgres
		DB_HOST=127.0.0.1
		DB_PORT=5432
		DB_NAME=postgres
		DATABASE_URL="postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME"
	`

	expected := map[string]string{
		"DB_USER":      "postgres",
		"DB_PASS":      "postgres",
		"DB_HOST":      "127.0.0.1",
		"DB_PORT":      "5432",
		"DB_NAME":      "postgres",
		"DATABASE_URL": "postgres://postgres:postgres@127.0.0.1:5432/postgres",
	}

	got, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Parse() error = %v.", err)
	}

	assertMap(t, got, expected)
}

func assertMap(t *testing.T, got, want map[string]string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got %d keys, want %d keys.", len(got), len(want))
	}

	for k, v := range want {
		if got[k] != v {
			t.Errorf("key %s: got %v, want %v.", k, got[k], v)
		}
	}
}
