package constants

type ColumnType string

const (
	INT  ColumnType = "INT"
	TEXT ColumnType = "TEXT"
)

type Column struct {
	Name string
	Type ColumnType
	Size int
}

type Schema struct {
	TableName string
	Columns   []Column
}

type Row []interface{}
var Schemas = map[string]Schema{}
var TableRows = make(map[string][]Row)
