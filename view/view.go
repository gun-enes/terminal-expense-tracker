package view

import (
	"expense-tracker/models"
	"expense-tracker/storage"
	"fmt"
	"strings"
	"text/tabwriter"
	"math"
	"os"
"unicode/utf8"
)

// Assuming your models package looks something like this, 
// I defined them locally for the snippet to work standalone.
var (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	StyleBold   = "\033[1m"
	ColorMagenta  = "\033[35m"
	ColorGray     = "\033[90m" // Bright Black (Gray) for dates
)

func CategoryExpenses(expList []storage.CategoryStat, total float64) {
	// 1. Initialize Tabwriter for perfect column alignment
	// minwidth, tabwidth, padding, padchar, flags
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// 2. Print Header
	fmt.Fprintf(w, "%sCATEGORY\tAMOUNT\t%%\tDISTRIBUTION%s\n", StyleBold, ColorReset)
	
	// Create a separator based on estimated width
	fmt.Fprintln(w, strings.Repeat("-", 60))

	// 3. Iterate and Print Rows
	for _, value := range expList {
		// Calculate percentage logic
		percentage := 0.0
		if total > 0 {
			percentage = (value.Total / total) * 100
		}

		// Generate Visual Bar (Max 20 chars width)
		barLength := int(math.Round(percentage / 100 * 20))
		bar := strings.Repeat("█", barLength)
		empty := strings.Repeat("░", 20-barLength)

		// Format the amount (Right aligned logic handled by logic or explicit spacing, 
		// but tabwriter handles the columns)
		amountStr := fmt.Sprintf("%.2f", value.Total)

		// Print the row
		// %s\t denotes a new tab column
		fmt.Fprintf(w, "%s%s\t%s%s\t%.1f%%\t%s%s%s%s\n", 
			ColorBlue, value.Name,         // Col 1: Name (Blue)
			ColorGreen, amountStr,         // Col 2: Amount (Green)
			percentage,                    // Col 3: Percentage
			ColorYellow, bar, empty, ColorReset, // Col 4: Visual Bar
		)
	}

	// 4. Footer (Total)
	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintf(w, "%sTOTAL\t%.2f\t100%%\t%s\n", StyleBold, total, ColorReset)

	// Flush buffer to terminal
	w.Flush()
}

// Define your visual styles
const (

)

func AllExpenses(expList []models.Expense) {
	// 1. Handle Empty State
	if len(expList) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	// 2. Initialize Tabwriter
	// minwidth, tabwidth, padding, padchar, flags
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// 3. Header
	// We use standard colors for the header
	fmt.Fprintf(w, "%sDATE\tCATEGORY\tDESCRIPTION\tAMOUNT%s\n", StyleBold, ColorReset)
	
	// A cleaner separator
	fmt.Fprintln(w, strings.Repeat("-", 80))

	totalAmount := 0.0

	// 4. Loop Rows
	for _, value := range expList {
		totalAmount += value.Amount

		// Format Date: Use Gray so it doesn't distract
		dateStr := fmt.Sprintf("%s%s%s", ColorGray, value.Date.Format("02 Jan"), ColorReset)

		// Format Category: Colored for quick scanning
		catStr := fmt.Sprintf("%s%s%s", ColorBlue, value.Category, ColorReset)

		// Format Description: Truncate intelligently
		cleanDesc := truncateText(value.Description, 35)
		
		// Format Amount: Bold and Green
		amountStr := fmt.Sprintf("%s%s%.2f%s", StyleBold, ColorGreen, value.Amount, ColorReset)

		// Write the row using tabs
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", dateStr, catStr, cleanDesc, amountStr)
	}

	// 5. Footer (Total)
	fmt.Fprintln(w, strings.Repeat("-", 80))
	// Empty tabs for the first two columns to align Total under Description/Amount
	fmt.Fprintf(w, "\t\t%sTOTAL\t%.2f%s\n", StyleBold, totalAmount, ColorReset)

	w.Flush()
}

// Helper to truncate string safely (handling emojis/utf8 correctly)
func truncateText(text string, maxLen int) string {
	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}
	// Convert to runes to slice characters, not bytes
	runes := []rune(text)
	return string(runes[:maxLen-3]) + "..."
}
