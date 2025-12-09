# Go Terminal Application To Track Expenses

## Build
```bash
go mod tidy
go build
```

## Commands
You can import files with this command. File has to be in the format:
[{"date": "2025-08-23 00:00:00", "description": "SOME EXPENSE", "amount": "100", "label": "General"}]
```bash
expense-tracker import file.txt/file.json
```

After importing, you need to classify the items. It auto suggests categories, and auto accepts some expenses by matching them with history.
```bash
expense-tracker classify 
```

You can view the expenses by below command. You can give it flags like --month to specify the month.
```bash
expense-tracker view
```


