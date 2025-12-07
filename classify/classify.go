package classify

import (
	"bufio"
	"expense-tracker/models"
	"expense-tracker/storage"
	"expense-tracker/utils"
	"fmt"
	"os"
	"slices"
	"strings"
)



func isInCategory(exp models.Expense) string {
	descUpper := strings.ToUpper(exp.Description)
	for _, cat := range models.Categories {
		for _, val := range cat.Aliases {
			if strings.Contains(descUpper, val) {
				return val
			}
		}
	}
	return ""
	}

func suggestedCategory(exp models.Expense, history []models.Expense) (string, bool) {
	// 1. Check History (Exact match)
	

	for _, val := range history {
		if val.Description == exp.Description {
			return val.Category, true // High confidence (Found in history)
		}
	}
	sigara := []string {"AKARYAKIT", "BUFE"}

	

	if (slices.Contains(sigara, exp.Label) && exp.Amount <= -90.00 && exp.Amount >= -105 && exp.Amount == float64(int(exp.Amount))){
		return "sigara", false
	}

	val := isInCategory(exp)
	if val != ""{
		return val, true
	}

	if exp.Label == "Yeme / İçme"{
		return "yemek", false
	}

	if exp.Label == "Turizm / Konaklama"{
		return "konaklama", false
	}

	if exp.Label == "Eğlence / Hobi"{
		return "eğlence", false
	}

	return strings.ToLower(exp.Label), false
}

func ProcessExpenses(db *storage.DB) []models.Expense {
	uncatItems, _ := db.GetUncategorized()
	allItems, _ := db.GetCategorized()
	var extList []models.Expense
	reader := bufio.NewReader(os.Stdin)

	totalItems := len(uncatItems)


	fmt.Println(models.StyleBold + models.ColorBlue + "--- STARTING CATEGORIZATION ---" + models.ColorReset)
	fmt.Printf("Processing %d transactions...\n\n", totalItems)

	// Create a header for the table
	fmt.Printf("%s%-12s %-30s %-10s %-20s%s\n", models.StyleBold, "DATE", "DESCRIPTION", "AMOUNT", "CATEGORY", models.ColorReset)
	fmt.Println(strings.Repeat("-", 80))

	for _, value := range uncatItems {
		sug, highConfidence := suggestedCategory(value, allItems)
		var finalCategory string

		// Format output variables
		dateStr := value.Date.Format("02 Jan")
		// Truncate description if too long for the table
		descStr := value.Description
		if len(descStr) > 28 {
			descStr = descStr[:25] + "..."
		}
		amountStr := fmt.Sprintf("%.2f", value.Amount)

		// 1. AUTO-ACCEPT: If we found this exact description in history, skip asking
		if highConfidence {
			finalCategory = sug
			// Print log line in gray so it doesn't distract
			fmt.Printf("%s%-12s %-30s %-10s %-20s [Auto]%s\n",
				models.ColorGray, dateStr, descStr, amountStr, sug, models.ColorReset)
			
			extList = append(extList, models.Expense{
				ID: value.ID,
				Description: value.Description, 
				Label: value.Label, 
				Amount: value.Amount, 
				Date: value.Date, 
				SourceFile: value.SourceFile, 
				Category: finalCategory})

			db.UpdateCategory(value.ID, finalCategory)

			continue
		}

		// 2. INTERACTIVE MODE
		// Highlight the current row
		fmt.Printf("%s%-12s %-30s %s%-10s%s ",
			models.ColorReset, dateStr, descStr, models.StyleBold, amountStr, models.ColorReset)

		// Show Suggestion Arrow
		fmt.Printf("%s---> %s%s%s ? ", models.ColorCyan, models.StyleBold, utils.Capitalize(sug), models.ColorReset)

		// 3. INPUT LOOP (Validation)
		for {
			// Read input
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "" {
				// User pressed Enter -> Accept Suggestion
				finalCategory = sug
				fmt.Printf("\033[1A\033[65C %s✓ Accepted%s\n", models.ColorGreen, models.ColorReset) // Move cursor up and right
				break
			}

			// Validate input against categories
			validCat := false
			foundName := ""
			for _, cat := range models.Categories {
				if strings.EqualFold(cat.Name, input) { // Case insensitive
					validCat = true
					foundName = cat.Name
					break
				}
			}

			if validCat {
				finalCategory = foundName
				fmt.Printf("\033[1A\033[65C %s✓ Set to: %s%s\n", models.ColorYellow, utils.Capitalize(foundName), models.ColorReset)

				break
			} else {
				// Error message
				fmt.Printf("%s   Invalid category! Try again: %s", models.ColorRed, models.ColorReset)
			}
		}
		exp := models.Expense{
			ID: value.ID,
			Description: value.Description, 
			Label: value.Label, 
			Amount: value.Amount, 
			Date: value.Date, 
			SourceFile: value.SourceFile, 
			Category: finalCategory}

		extList = append(extList, exp)
		allItems = append(allItems, exp)


		db.UpdateCategory(value.ID, finalCategory)
	}

	fmt.Println("\n" + models.StyleBold + "--- SUMMARY ---" + models.ColorReset)
	fmt.Printf("Categorized %d expenses.\n", len(extList))

	return extList
}


