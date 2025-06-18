package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"minidb/constants"
	"os"
	"strings"
)

func ExitProgram() {
	os.Exit(0)
}

func serializeRow(row constants.Row, schema constants.Schema) []byte {
	buf := new(bytes.Buffer)

	for i, col := range schema.Columns {
		switch col.Type {
		case constants.INT:
			_ = binary.Write(buf, binary.LittleEndian, row[i].(int32))
		case constants.TEXT:
			txt := make([]byte, col.Size)
			copy(txt, []byte(row[i].(string)))
			buf.Write(txt)
		}
	}
	return buf.Bytes()
}

func deserializeRow(data []byte, schema constants.Schema) constants.Row {
	row := constants.Row{}
	offset := 0

	for _, col := range schema.Columns {
		switch col.Type {
		case constants.INT:
			val := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
			row = append(row, val)
			offset += 4
		case constants.TEXT:
			str := string(data[offset : offset+col.Size])
			row = append(row, string(bytes.TrimRight([]byte(str), "\x00")))
			offset += col.Size
		}
	}
	return row
}

func SaveSchemasToFile() {
	file, err := os.Create("schemas/schema.json")
	if err != nil {
		fmt.Println("Error saving schema:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(constants.Schemas)
	if err != nil {
		fmt.Println("Error encoding schema:", err)
	}
}

func LoadSchemasFromFile() {
	file, err := os.Open("schemas/schema.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error loading schema:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&constants.Schemas)
	if err != nil {
		fmt.Println("Error decoding schema:", err)
		return
	}
}

func AppendRowToFile(table string, row constants.Row, schema constants.Schema) {
	f, err := os.OpenFile("db/"+table+".db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	defer f.Close()

	bin := serializeRow(row, schema)
	f.Write(bin)
}

func LoadRowsFromFile(table string, schema constants.Schema) {
	data, err := os.ReadFile("db/"+table + ".db")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var rows []constants.Row
	rowSize := 0
	for _, col := range schema.Columns {
		rowSize += col.Size
	}

	for i := 0; i+rowSize <= len(data); i += rowSize {
		row := deserializeRow(data[i:i+rowSize], schema)
		rows = append(rows, row)
	}

	constants.TableRows[table] = rows
	fmt.Printf("Loaded %d rows from %s.db\n", len(rows), table)
}

func ParseCSV(input string) []string {
	var out []string
	var current strings.Builder
	inQuotes := false

	for _, r := range input {
		switch r {
		case ',':
			if inQuotes {
				current.WriteRune(r)
			} else {
				out = append(out, strings.TrimSpace(current.String()))
				current.Reset()
			}
		case '"':
			inQuotes = !inQuotes
		default:
			current.WriteRune(r)
		}
	}
	out = append(out, strings.TrimSpace(current.String()))
	return out
}

func PrintHelp() {
	fmt.Println()
	fmt.Println("Supported commands:")
	fmt.Println("select * from <table_name>")
	fmt.Println("insert into <table_name> values (<values>)")
	fmt.Println("delete from <table_name> where <column_name> = <value>")
	fmt.Println("create table <table_name> (<column_name> <column_type>)")
	fmt.Println("show tables")
	fmt.Println("drop table <table_name>")
	fmt.Println()
}
