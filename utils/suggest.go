package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SuggestivePrint(validOptions []string) string {

	disableInputBuffering()
	
	defer restoreTerminal()
	
	result := autoCompleteInput(validOptions)
	
	restoreTerminal() // Restore before printing final result
	return result
}

// autoCompleteInput handles the keystroke logic and visual rendering
func autoCompleteInput(options []string) string {
	var input string

	for {
		// 2. Read exactly one byte (one key press)
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			break
		}
		char := b[0]

		// 3. Handle Special Keys

		// Tab key (9) - Autocomplete without submitting
		if char == 9 {
			suggestion := findSuggestion(input, options)
			if suggestion != "" {
				input = suggestion
			}
		}

		// Enter key (10 or 13)
		if char == 10 || char == 13 {
			// If we have a valid match available, autocomplete it on Enter
			suggestion := findSuggestion(input, options)
			if suggestion != "" {
				input = suggestion
			}
			break
		}

		// Backspace (127)
		if char == 127 {
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		} else if char >= 32 && char <= 126 {
			// Normal printable characters
			input += string(char)
		} else if char == 3 {
			// Ctrl+C to exit
			restoreTerminal()
			os.Exit(0)
		}

		// 4. Find the best match
		suggestion := findSuggestion(input, options)

		// 5. Render the line
		renderLine(input, suggestion)
	}

	return input
}

// findSuggestion looks for the first string in the list that starts with input
func findSuggestion(input string, options []string) string {
	if input == "" {
		return ""
	}
	for _, opt := range options {
		if strings.HasPrefix(opt, input) {
			return opt
		}
	}
	return ""
}

// renderLine clears the console line and prints the input + ghost text
func renderLine(input, suggestion string) {
	// ANSI Escape Codes:
	// \r       -> Move cursor to start of line
	// \033[K   -> Clear everything from cursor to end of line
	// \033[2m  -> Dim/Faint text (Gray)
	// \033[0m  -> Reset text formatting
	
	fmt.Print("\r\033[K") // Clean the line
	
	// Print what the user actually typed
	fmt.Print(input)

	// If there is a suggestion that is longer than the input, print the rest as "ghost text"
	if suggestion != "" && len(suggestion) > len(input) {
		remainder := suggestion[len(input):]
		fmt.Printf("\033[2m%s\033[0m", remainder)
		
		// Move cursor back to where the user is typing (before the ghost text)
		// \033[XD moves cursor left by X positions
		fmt.Printf("\033[%dD", len(remainder))
	}
}

// Utility to enable raw mode (Unix/Linux/Mac specific)
func disableInputBuffering() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

// Utility to restore terminal (Unix/Linux/Mac specific)
func restoreTerminal() {
	exec.Command("stty", "-F", "/dev/tty", "-cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}
