package utils

import (
	"expense-tracker/models"
	"fmt"
	"unicode"
)

func MonthStringToInt(monthPtr *string) int {
		var m int
		switch *monthPtr {
		case "jan":
			m = 1
		case "feb":
			m = 2
		case "mar":
			m = 3
		case "apr":
			m = 4
		case "may":
			m = 5
		case "jun":
			m = 6
		case "jul":
			m = 7
		case "aug":
			m = 8
		case "sep":
			m = 9
		case "oct":
			m = 10
		case "nov":
			m = 11
		case "dec":
			m = 12
		}
		return m
}



func PrettyPrint(desc string, sug string) {
	const maxWidth = 30

	// 1. Capitalize safely
	sug = Capitalize(sug)

	// 2. Handle Description Truncation with Runes (safe for emojis/unicode)
	// We convert to a slice of runes to handle characters correctly
	descRunes := []rune(desc)
	var formattedDesc string

	if len(descRunes) > maxWidth {
		// Cut at maxWidth - 3 to make room for "..."
		formattedDesc = string(descRunes[:maxWidth-3]) + "..."
	} else {
		formattedDesc = string(descRunes)
	}

	// 3. Print with colors and alignment
	// %-30s pads the string with spaces on the right to fill 30 chars
	fmt.Printf(
		"%s%-30s%s %s--->%s    %s%s%s\n",
		models.ColorCyan, formattedDesc, models.ColorReset, // The Description (Cyan)
		models.ColorYellow, models.ColorReset,              // The Arrow (Yellow)
		models.ColorGreen, sug, models.ColorReset,          // The Suggestion (Green)
	)
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	// Convert to rune slice to handle non-ASCII characters correctly
	r := []rune(s)
	return string(unicode.ToUpper(r[0])) + string(r[1:])
}
