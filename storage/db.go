package storage

import (
	"database/sql"
	"expense-tracker/models"
	"expense-tracker/utils"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

type CategoryStat struct {
	Name  string
	Total float64
}

func (d *DB) AverageDailySpending()int{
	d.Query(`
	SELECT 
	`)
	return 0
}
func (d *DB) DropTable() error {
	createTableSQL := `
	DROP TABLE expenses;
	`

	_, err := d.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("could not create tables: %v", err)
	}

	return nil
}



func (d *DB) GetMonthStats(m int) ([]CategoryStat, float64, error) {
	rows, err := d.Query(`
		SELECT category, SUM(amount) 
		FROM expenses 
		WHERE strftime('%m', date) = printf('%02d', ?)
		GROUP BY category 
		ORDER BY SUM(amount) ASC 
	`, m)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stats []CategoryStat
	var grandTotal float64

	for rows.Next() {
		var s CategoryStat
		if err := rows.Scan(&s.Name, &s.Total); err != nil {
			return nil, 0, err
		}
		if s.Total < 0 {
			s.Total = -s.Total
		}
		grandTotal += s.Total
		stats = append(stats, s)
	}

	return stats, grandTotal, nil
}


func (d *DB) GetCategoryStats() ([]CategoryStat, float64, error) {
	rows, err := d.Query(`
		SELECT category, SUM(amount) 
		FROM expenses 
		GROUP BY category 
		ORDER BY SUM(amount) ASC
	`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stats []CategoryStat
	var grandTotal float64

	for rows.Next() {
		var s CategoryStat
		if err := rows.Scan(&s.Name, &s.Total); err != nil {
			return nil, 0, err
		}
		s.Name = utils.Capitalize(s.Name)
		if s.Total < 0 {
			s.Total = -s.Total
		}
		grandTotal += s.Total
		stats = append(stats, s)
	}

	return stats, grandTotal, nil
}

type DB struct {
	*sql.DB
}

func InitDB(filepath string) (*DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	// Create Tables
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT NOT NULL DEFAULT 'Uncategorized',
		source_file TEXT,
		label TEXT
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create tables: %v", err)
	}

	return &DB{db}, nil
}

func (d *DB) SaveExpense(e models.Expense) error {
	stmt, err := d.Prepare("INSERT INTO expenses(date, description, amount, category, source_file, label) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	var cat string
	if e.Category != "" {
		cat = e.Category
	}else{
		cat = "Uncategorized"
	}

	_, err = stmt.Exec(e.Date.Format(time.RFC3339), e.Description, e.Amount, cat, e.SourceFile, e.Label)
	if err != nil {
		return fmt.Errorf("insert error: %v", err)
	}
	return nil
}

func (d *DB) UpdateCategory(id int, newCategory string) error {
	_, err := d.Exec("UPDATE expenses SET category = ? WHERE id = ?", newCategory, id)
	return err
}

func (d *DB) GetCategorized() ([]models.Expense, error) {
	rows, err := d.Query("SELECT id, date, description, amount, category, source_file, label FROM expenses WHERE category != 'Uncategorized' ORDER BY date ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return convertRowsToModel(rows), nil
}

func (d *DB) GetUncategorized() ([]models.Expense, error) {
	rows, err := d.Query("SELECT id, date, description, amount, category, source_file, label FROM expenses WHERE category = 'Uncategorized' ORDER BY date ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return convertRowsToModel(rows), nil
}

func (d *DB) GetTopExpenses(month int) ([]models.Expense, error) {
	rows, err := d.Query(`
		SELECT id, date, description, amount, category, source_file, label
		FROM expenses 
		WHERE strftime('%m', date) = printf('%02d', ?)
		ORDER BY amount ASC
		LIMIT 5`, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return convertRowsToModel(rows), nil
}

func (d *DB) GetAllExpenses() ([]models.Expense, error) {
	rows, err := d.Query("SELECT id, date, description, amount, category, source_file, label FROM expenses ORDER BY date ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return convertRowsToModel(rows), nil
}

func convertRowsToModel(rows *sql.Rows)[]models.Expense{
	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var dateStr string
			if err := rows.Scan(&e.ID, &dateStr, &e.Description, &e.Amount, &e.Category, &e.SourceFile, &e.Label); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		e.Date, _ = time.Parse(time.RFC3339, dateStr)
		expenses = append(expenses, e)
	}
	return expenses 
}  


