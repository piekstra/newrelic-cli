package confirm

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"yes lowercase", "y\n", true},
		{"yes full", "yes\n", true},
		{"Yes uppercase", "Y\n", true},
		{"YES all caps", "YES\n", true},
		{"no lowercase", "n\n", false},
		{"no full", "no\n", false},
		{"No uppercase", "No\n", false},
		{"empty default no", "\n", false},
		{"random input", "maybe\n", false},
		{"whitespace yes", "  y  \n", true},
		{"whitespace no", "  n  \n", false},
		{"partial yes", "ye\n", false},
		{"yess typo", "yess\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompter{
				In:  strings.NewReader(tt.input),
				Out: io.Discard,
			}
			result := p.Confirm("Test?")
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfirm_EOF(t *testing.T) {
	p := &Prompter{
		In:  strings.NewReader(""), // EOF
		Out: io.Discard,
	}
	result := p.Confirm("Test?")
	assert.False(t, result, "EOF should return false")
}

func TestConfirmDanger(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		confirmWord string
		expected    bool
	}{
		{"exact match", "delete\n", "delete", true},
		{"wrong word", "remove\n", "delete", false},
		{"empty input", "\n", "delete", false},
		{"partial match", "delet\n", "delete", false},
		{"extra chars", "deletee\n", "delete", false},
		{"case sensitive", "DELETE\n", "delete", false},
		{"whitespace trimmed", "  delete  \n", "delete", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompter{
				In:  strings.NewReader(tt.input),
				Out: io.Discard,
			}
			result := p.ConfirmDanger("Test?", tt.confirmWord)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfirmDanger_EOF(t *testing.T) {
	p := &Prompter{
		In:  strings.NewReader(""), // EOF
		Out: io.Discard,
	}
	result := p.ConfirmDanger("Test?", "delete")
	assert.False(t, result, "EOF should return false")
}

func TestConfirm_OutputPrompt(t *testing.T) {
	var output strings.Builder
	p := &Prompter{
		In:  strings.NewReader("y\n"),
		Out: &output,
	}

	p.Confirm("Delete this item?")

	assert.Equal(t, "Delete this item? [y/N]: ", output.String())
}

func TestConfirmDanger_OutputPrompt(t *testing.T) {
	var output strings.Builder
	p := &Prompter{
		In:  strings.NewReader("delete\n"),
		Out: &output,
	}

	p.ConfirmDanger("This will permanently delete all data.", "delete")

	expected := "This will permanently delete all data.\nType 'delete' to confirm: "
	assert.Equal(t, expected, output.String())
}
