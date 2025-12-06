package models

import "time"


var fileList = []string{
	"files/2025-08.txt",
	"files/2025-09.txt",
	// Add more files here as needed
}


// --- VISUAL CONSTANTS (ANSI Colors) ---
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	StyleBold   = "\033[1m"
)

// --- CONFIGURATION ---
var Categories = []Category{
	{Name: "kafe", Aliases: []string{"KAFE", "CAFE", "COFFE", "STARBUCKS", "ESPRESSO"}},
	{Name: "yeme / içme", Aliases: []string{"RESTORAN", "BURGER", "KFC", "PIZZA"}},
	{Name: "market", Aliases: []string{"MARKET", "FIRIN", "MİGROS", "BİM", "A101", "CARREFOUR"}},
	{Name: "alkol", Aliases: []string{"ALKOL", "BAR", "TEKEL"}},
	{Name: "sigara", Aliases: []string{"SIGARA"}},
	{Name: "ulaşım", Aliases: []string{"ULAŞIM", "EGO", "TAKSI", "UBER", "MARTI"}},
	{Name: "seyahat", Aliases: []string{"SEYAHAT", "PEGASUS", "THY", "OTEL", "BNB"}},
	{Name: "hediye", Aliases: []string{"HEDIYE", "CICEK"}},
	{Name: "alışveriş", Aliases: []string{"ALIŞVERIŞ", "GİYİM", "ZARA", "MANGO", "TRENDYOL", "HEPSIBURADA"}},
	{Name: "bakım", Aliases: []string{"BAKIM", "KUAFOR", "BERBER", "ECZANE"}},
	{Name: "faturalar", Aliases: []string{"FATURA", "TURKCELL", "VODAFONE", "ENERJISA", "ASKI"}},
	{Name: "eğlence", Aliases: []string{"EĞLENCE", "NETFLIX", "SPOTIFY", "STEAM", "SINEMA"}},
	{Name: "kira", Aliases: []string{"KIRA"}},
	{Name: "konaklama", Aliases: []string{"KONAKLAMA"}},
}

// --- STRUCTS ---

type Expense struct {
	ID          int    
	Date        time.Time
	Description string
	Amount      float64
	Label       string
	Category 	string
	SourceFile 	string
}

type Category struct {
	ID      int    
	Name    string
	Aliases []string
}


