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

func LoadFile(filePath string) ([]models.Expense, error) {
	var cleanExpenses []models.Expense

	// 1. Read the specific file provided
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	// 2. Unmarshal JSON
	var rawList []RawExpense
	if err := json.Unmarshal(content, &rawList); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %v", filePath, err)
	}

	// 3. Convert Raw to Clean Model
	for _, raw := range rawList {
		// Skip NaNs
		if raw.Amount == "nan" {
			continue
		}

		// Parse Amount
		amount, err := strconv.ParseFloat(raw.Amount, 64)
		if err != nil || math.IsNaN(amount) {
			continue
		}

		// Skip positive amounts (Income/Refunds)
		if amount > 0 {
			continue 
		}

		// Parse Date
		t, err := time.Parse("2006-01-02 15:04:05", raw.DateRaw)
		if err != nil {
			continue
		}

		cleanExpenses = append(cleanExpenses, models.Expense{
			Date:        t,
			Description: raw.Description,
			Label:       raw.Label,
			Amount:      amount,
			// Use filepath.Base to get just "data.json" from "path/to/data.json"
			SourceFile:  filepath.Base(filePath),
		})
	}

	return cleanExpenses, nil
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
		exps, _ := LoadFile(fullPath)
		cleanExpenses = append(cleanExpenses, exps...)
	}
	return cleanExpenses, nil
}
