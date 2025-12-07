package view

import (
	"expense-tracker/storage"
	"fmt"
	"math"
	"sort"
)

// ANSI Color codes for terminal output
const (
	ColorBold   = "\033[1m"
)
func CompareMonth(m int, db *storage.DB) {
	// [Fetching data logic remains same as before...] 
	curStat, curTotal, err := db.GetMonthStats(m)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if len(curStat) == 0 {
		fmt.Println("No records for current month.")
		return
	}

	prevMonth := m - 1
	if prevMonth == 0 { prevMonth = 12 }
	prevStat, prevTotal, _ := db.GetMonthStats(prevMonth)

	// Prepare Map
	prevMap := make(map[string]float64)
	for _, s := range prevStat {
		prevMap[s.Name] = s.Total
	}

	// --- HEADERS ---
	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("%sMONTHLY OVERVIEW (Month %d vs %d)%s\n", ColorBold, m, prevMonth, ColorReset)
	fmt.Println("---------------------------------------------------------------")

	// Print Totals
	diff := curTotal - prevTotal
	percentChange := 0.0
	if prevTotal > 0 {
		percentChange = (diff / prevTotal) * 100
	}

	totalColor := ColorReset
	indicator := "="
	if diff > 0 {
		totalColor = ColorRed
		indicator = "↑"
	} else if diff < 0 {
		totalColor = ColorGreen
		indicator = "↓"
	}

	fmt.Printf("Total:      ₺%-10.2f\n", curTotal)
	fmt.Printf("Previous:   ₺%-10.2f\n", prevTotal)
	// We format the color specifically around the diff values, not the label
	fmt.Printf("Difference: %s%s ₺%.2f (%.1f%%)%s\n", totalColor, indicator, math.Abs(diff), math.Abs(percentChange), ColorReset)
	fmt.Println()

	// --- TABLE HEADER ---
	// %-15s : Left-align text in a 15-char wide space
	// %10s  : Right-align text in a 10-char wide space
	fmt.Printf("%-15s %-12s %-12s %-12s %-10s\n", "CATEGORY", "CURRENT", "PREVIOUS", "DIFF", "CHANGE")
	fmt.Println("---------------------------------------------------------------")

	// Sort Current
	sort.Slice(curStat, func(i, j int) bool {
		return curStat[i].Total > curStat[j].Total
	})

	// --- TABLE ROWS ---
	for _, cat := range curStat {
		prevAmount := prevMap[cat.Name]
		catDiff := cat.Total - prevAmount
		
		catPercent := 0.0
		if prevAmount > 0 {
			catPercent = (catDiff / prevAmount) * 100
		} else if prevAmount == 0 && cat.Total > 0 {
			catPercent = 100.0
		}

		// Pick Color
		rowColor := ColorReset
		sign := " "
		if catDiff < 0 {
			rowColor = ColorGreen
			sign = "-"
		} else if catDiff > 0 {
			rowColor = ColorRed
			sign = "+"
		}

		delete(prevMap, cat.Name)

		// ALIGNMENT FIX:
		// We print the Color Code (%s), then the Data Columns, then Reset (%s).
		// Note the widths: %-15s for name, %10.2f for money.
		// Since '₺' is multi-byte, it can sometimes shift alignment by 1 char in some terminals. 
		// If so, add a space after the symbol or remove it from the width calculation.
		
		fmt.Printf("%s%-15s ₺%-10.2f ₺%-10.2f %s₺%-10.2f %6.1f%%%s\n",
			rowColor,            // 1. Set Color
			cat.Name,            // 2. Category Name (padded to 15)
			cat.Total,           // 3. Current (padded to 10)
			prevAmount,          // 4. Previous (padded to 10)
			sign,                // 5. +/- Sign
			math.Abs(catDiff),   // 6. Diff Amount
			math.Abs(catPercent),// 7. Percentage
			ColorReset,          // 8. Reset Color
		)
	}

	// --- STOPPED SPENDING ---
	if len(prevMap) > 0 {
		fmt.Println("\nStopped Spending:")
		for name, amount := range prevMap {
			fmt.Printf("%s%-15s ₺0.00        ₺%-10.2f -₺%-10.2f -100.0%%%s\n", 
				ColorGreen, name, amount, amount, ColorReset)
		}
	}
}
