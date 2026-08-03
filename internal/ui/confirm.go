package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirm asks a yes/no question on stdin.
func Confirm(prompt string, defaultYes bool) (bool, error) {
	return confirmWithReader(os.Stdin, prompt, defaultYes)
}

func confirmWithReader(r io.Reader, prompt string, defaultYes bool) (bool, error) {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}
	fmt.Printf("%s %s ", prompt, suffix)
	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	line = strings.ToLower(strings.TrimSpace(line))
	if line == "" {
		return defaultYes, nil
	}
	return line == "y" || line == "yes", nil
}
