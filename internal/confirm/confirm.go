package confirm

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Prompter handles interactive confirmations
type Prompter struct {
	In  io.Reader
	Out io.Writer
}

// Confirm prompts the user with a yes/no question
// Returns true if user confirms, false otherwise
// Default (empty input) returns false for safety
func (p *Prompter) Confirm(message string) bool {
	_, _ = fmt.Fprintf(p.Out, "%s [y/N]: ", message)

	reader := bufio.NewReader(p.In)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response := strings.ToLower(strings.TrimSpace(input))
	return response == "y" || response == "yes"
}

// ConfirmDanger prompts for dangerous operations with explicit typing
// User must type the confirmWord exactly to confirm
func (p *Prompter) ConfirmDanger(message, confirmWord string) bool {
	_, _ = fmt.Fprintf(p.Out, "%s\nType '%s' to confirm: ", message, confirmWord)

	reader := bufio.NewReader(p.In)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	return strings.TrimSpace(input) == confirmWord
}
