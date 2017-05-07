package parser

import "fmt"

// Table is what is parsed from table schema
type Table struct {
	Name    string
	DBName  string
	GoName  string
	Columns []*Column
}

// Column is Table's column
type Column struct {
	Name     string
	DBName   string
	GoName   string
	Type     ColumnType
	Nullable bool
	Length   int
	Refer    ColumnName
}

// PutCreateTableQuery puts CREATE TABLE query of the table
func (t *Table) PutCreateTableQuery() (query string) {
	query += fmt.Sprintf("CREATE TABLE `%s` (\n", t.DBName)

	for i, c := range t.Columns {
		if i != 0 {
			query += ",\n" // separator
		}

		query += "  " // indent
		query += c.putCreateTableQuery()
	}

	query += "\n);"
	return
}

func (c *Column) putCreateTableQuery() (query string) {
	if c.Type == Reference {
		panic(fmt.Errorf(`reference is not resolved`))
	}

	attr, ok := columnTypes[c.Type]
	if !ok {
		panic(fmt.Errorf(`no such type %s`, c.Type))
	}

	query += fmt.Sprintf("`%s` %s", c.DBName, attr.TypeName)

	if attr.Length {
		query += fmt.Sprintf("(%d)", c.Length) // char length
	}

	if !c.Nullable {
		query += " NOT NULL"
	}

	if attr.Unsigned {
		query += " UNSIGNED"
	}
	return
}
