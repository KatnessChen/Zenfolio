package prompts

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed *.txt
var promptFiles embed.FS

// LoadPrompt loads a prompt file by name and returns its content
func LoadPrompt(filename string) (string, error) {
	content, err := promptFiles.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file '%s': %w", filename, err)
	}
	return strings.TrimSpace(string(content)), nil
}
