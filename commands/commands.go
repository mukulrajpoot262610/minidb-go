package commands

import (
	"fmt"
	"minidb/constants"
	"minidb/utils"
	"os"
	"strconv"
	"strings"
)

func Select(command string) {
	fmt.Println("Executing SELECT")

	lower := strings.ToLower(command)
	if !strings.HasPrefix(lower, "select * from ") {
		fmt.Println("Only 'SELECT * FROM <table>' is supported.")
		return
	}

	tableName := strings.TrimSpace(command[len("select * from "):])
	schema, ok := constants.Schemas[tableName]
	if !ok {
		fmt.Println("Table not found:", tableName)
		return
	}

	// Load rows if not already in memory
	if _, ok := constants.TableRows[tableName]; !ok {
		utils.LoadRowsFromFile(tableName, schema)
	}

	rows := constants.TableRows[tableName]
	if len(rows) == 0 {
		fmt.Println("No rows found.")
		return
	}

	// Print header
	for _, col := range schema.Columns {
		fmt.Printf("%-*s ", col.Size, col.Name)
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for _, cell := range row {
			fmt.Printf("%-*v ", 16, cell) // Fixed width columns
		}
		fmt.Println()
	}

}

func Insert(command string) {

	// insert into users values (1, "", "")

	lower := strings.ToLower(command)
	if !strings.Contains(lower, "values") {
		fmt.Println("Invalid insert syntax. Use: insert into <table_name> values (<values>)")
		return
	}

	parts := strings.SplitN(command, "values", 2)
	intoPart := strings.TrimSpace(parts[0])
	valuesPart := strings.TrimSpace(parts[1])

	intoTokens := strings.Fields(intoPart)
	if len(intoTokens) < 3 {
		fmt.Println("Invalid INSERT INTO syntax. Use: insert into <table_name> values (<values>)")
		return
	}

	tableName := intoTokens[2]
	schema, ok := constants.Schemas[tableName]
	if !ok {
		fmt.Println("Table does not exist:", tableName)
		return
	}

	valuesRaw := strings.Trim(valuesPart, "()")
	valuesTokens := utils.ParseCSV(valuesRaw)
	if len(valuesTokens) != len(schema.Columns) {
		fmt.Printf("Expected %d values, got %d\n", len(schema.Columns), len(valuesTokens))
		return
	}

	var row []any
	for i, col := range schema.Columns {
		val := strings.Trim(valuesTokens[i], `"`)
		switch col.Type {
		case constants.INT:
			num, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("Expected INT value for", col.Name)
				return
			}
			row = append(row, int32(num))

		case constants.TEXT:
			if len(val) > col.Size {
				fmt.Printf("Value too long for column %s (max %d)\n", col.Name, col.Size)
				return
			}
			row = append(row, val)
		}
	}


	constants.TableRows[tableName] = append(constants.TableRows[tableName], row)
	utils.AppendRowToFile(tableName, row, schema)

	fmt.Println("Row inserted.")
}

func Delete(command string) {
	fmt.Println("Executing DELETE")

	parts := strings.Split(command, " ")
	if len(parts) != 3 {
		fmt.Println("Syntax error. Use: delete from <table_name> where <id>")
		return
	}

	tableName := parts[2]
	id, err := strconv.Atoi(parts[4])
	if err != nil {
		fmt.Println("Invalid ID. Must be a number.")
		return
	}

	found := false
	for i, row := range constants.TableRows[tableName] {
		if row[0].(int32) == int32(id) {
				constants.TableRows[tableName] = append(constants.TableRows[tableName][:i], constants.TableRows[tableName][i+1:]...)
				found = true
			break
		}
	}

	if !found {
		fmt.Printf("No row found with ID %d\n", id)
		return
	}

	fmt.Printf("Deleted row with ID %d\n", id)
}

func CreateTable(command string) {
	fmt.Println("Executing CREATE TABLE")

	// CREATE TABLE users (id INT, name TEXT(32), email TEXT(255));

	def := strings.TrimPrefix(command, "create table ")
	def = strings.TrimSpace(def)
	openIdx := strings.Index(def, "(")
	closeIdx := strings.LastIndex(def, ")")

	if (openIdx == -1 || closeIdx == -1 || openIdx > closeIdx) {
		fmt.Println("Syntax error. Use: create table <table_name> (<column_name> <column_type>(<column_size>))")
		return
	}

	tableName := strings.TrimSpace(def[:openIdx])
	columnsDef := def[openIdx+1:closeIdx]
	columnsParts := strings.Split(columnsDef, ",")

	var columns []constants.Column
	for _, colDef := range columnsParts {
		colDef = strings.TrimSpace(colDef)
		tokens := strings.Fields(colDef)
		if len(tokens) < 2 {
			fmt.Println("Invalid column definition:", colDef)
			return
		}

		colName := tokens[0]
		typeRaw := tokens[1]

		var colType constants.ColumnType
		var size int

		if strings.HasPrefix(typeRaw, "TEXT") {
			colType = constants.TEXT
			open := strings.Index(typeRaw, "(")
			close := strings.Index(typeRaw, ")")
			if open == -1 || close == -1 {
				fmt.Println("TEXT type must specify size.")
				return
			}
			sizeVal := typeRaw[open+1 : close]
			s, err := strconv.Atoi(sizeVal)
			if err != nil {
				fmt.Println("Invalid TEXT size.")
				return
			}
			size = s
		} else if typeRaw == "INT" {
			colType = constants.INT
			size = 4
		} else {
			fmt.Println("Unsupported column type:", typeRaw)
			return
		}

		columns = append(columns, constants.Column{
			Name: colName,
			Type: colType,
			Size: size,
		})
	}

	constants.Schemas[tableName] = constants.Schema{
		TableName: tableName,
		Columns: columns,
	}

	_, err := os.Create("db/"+tableName + ".db")
	if err != nil {
		fmt.Println("Error creating db file:", err)
		return
	}
	utils.SaveSchemasToFile()

	fmt.Printf("Created table %s\n", tableName)
}

func ShowTables() {
	fmt.Println("Executing SHOW TABLES")
	for tableName := range constants.Schemas {
		fmt.Println(tableName)
	}
}

func DropTable(command string) {
	fmt.Println("Executing DROP TABLE")

	parts := strings.Split(command, " ")

	if len(parts) != 2 {
		fmt.Println("Syntax error. Use: drop_table <table_name>")
		return
	}

	tableName := parts[1]

	constants.TableRows[tableName] = make([]constants.Row, 0)

	fmt.Printf("Dropped table %s\n", tableName)
}