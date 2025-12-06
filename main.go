package main

import (
	"expense-tracker/classify"
	"expense-tracker/ingest"
	"expense-tracker/storage"
	"expense-tracker/utils"
	"expense-tracker/view"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 1. Setup Paths
	// We assume the app is run from the root expense-cli/ folder
	dataDir := "./data"
	filesDir := filepath.Join(dataDir, "files")
	dbPath := filepath.Join(dataDir, "expenses.db")

	// 2. Initialize Database
	db, err := storage.InitDB(dbPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()


	// 3. Handle CLI Commands
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	classifyCmd := flag.NewFlagSet("classify", flag.ExitOnError)
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	viewMonth := viewCmd.String("month", "", "Group expenses by: 'month' or 'category'")

	switch os.Args[1] {
	case "import":
		importCmd.Parse(os.Args[2:])
		HandleImport(db, filesDir)

	case "classify":
		classifyCmd.Parse(os.Args[2:])
		classify.ProcessExpenses(db)
	
	case "stats":
		cat, tot, _ := db.GetCategoryStats()
		view.CategoryExpenses(cat, tot)

	case "view":
		viewCmd.Parse(os.Args[2:])
		if *viewMonth != "" {
			m := utils.MonthStringToInt(viewMonth)
			expList, total,_ := db.GetMonthStats(m)
			topExp, _ := db.GetTopExpenses(m)
			view.CategoryExpenses(expList, total)
			view.AllExpenses(topExp)
		} else {
			i := 0
			for i < 12{
				expList, total,_ := db.GetMonthStats(i)
				if len(expList) > 0{
					view.CategoryExpenses(expList, total)
				}
				i++
			}
		}


	default:
		fmt.Println("Expected 'import', 'classify', or 'stats' subcommands")
	}
}


func HandleImport(db *storage.DB, dir string) {
	fmt.Println("Scanning for files in:", dir)
	
	expenses, err := ingest.LoadFiles(dir)
	if err != nil {
		fmt.Printf("Error loading files: %v\n", err)
		return
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses found in files.")
		return
	}

	fmt.Printf("Found %d items. syncing to DB...\n", len(expenses))

	added := 0
	duplicates := 0

	for _, e := range expenses {
		err := db.SaveExpense(e)
		if err != nil {
			// If error is due to unique constraint, it's a duplicate
			duplicates++
		} else {
			added++
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("New Items Added: %d\n", added)
	fmt.Printf("Duplicates Skipped: %d\n", duplicates)
	fmt.Printf("Total in DB: %d\n", added+duplicates) // Approximation
}
