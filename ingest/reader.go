package ingest

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"expense-tracker/models"
)

type RawExpense struct {
	ID          int       `json:"id"` // Database ID
	DateRaw     string `json:"date"`
	Description string `json:"description"`
	Amount      string `json:"amount"` 
	Label       string `json:"label"`
}

// LoadFiles reads all JSON files in the dir and returns clean Expense objects
func LoadFiles(dirPath string) ([]models.Expense, error) {
	var cleanExpenses []models.Expense

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".txt" && filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		fullPath := filepath.Join(dirPath, file.Name())
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", file.Name(), err)
			continue
		}

		var rawList []RawExpense
		if err := json.Unmarshal(content, &rawList); err != nil {
			fmt.Printf("Skipping %s: invalid JSON\n", file.Name())
			continue
		}

		// Convert Raw to Clean Model
		for _, raw := range rawList {
			// 1. Skip NaNs
			if raw.Amount == "nan" {
				continue
			}

			// 2. Parse Amount
			amount, err := strconv.ParseFloat(raw.Amount, 64)
			if err != nil || math.IsNaN(amount) {
				continue
			}

			// 3. Skip positive amounts (Income/Refunds) - optional, based on your preference
			if amount > 0 {
				continue 
			}

			// 4. Parse Date
			// Adjust layout if your input date format differs
			t, err := time.Parse("2006-01-02 15:04:05", raw.DateRaw)
			if err != nil {
				// Fallback or skip
				continue
			}

			cleanExpenses = append(cleanExpenses, models.Expense{
				Date:        t,
				Description: raw.Description,
				Label: 		 raw.Label,
				Amount:      amount, // Keeping it negative to show expense
				SourceFile:  file.Name(),
			})
		}
	}
	return cleanExpenses, nil
}
