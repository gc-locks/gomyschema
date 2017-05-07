package parser

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/gedex/inflector"
	"github.com/serenize/snaker"
)

type tableScheme struct {
	Name    string  `xml:"name,attr"`
	GoName  *string `xml:"goname,attr"`
	DBName  *string `xml:"dbname,attr"`
	Columns []struct {
		Name          string  `xml:"name,attr"`
		GoName        *string `xml:"goname,attr"`
		DBName        *string `xml:"dbname,attr"`
		Type          *string `xml:"type,attr"`
		Nullable      *bool   `xml:"nullable,attr"`
		Length        *int    `xml:"length,attr"`
		Refer         *string `xml:"refer,attr"`
		AutoIncrement *bool   `xml:"auto_increment,attr"`
	} `xml:"column"`
}

// ColumnName is used for reference of Column
type ColumnName struct {
	Table  string
	Column string
}

// Parse byte data to table object
func Parse(data []byte) (table *Table, err error) {
	t := tableScheme{}
	xml.Unmarshal(data, &t)

	table = &Table{Name: t.Name}
	if t.DBName != nil {
		table.DBName = *t.DBName
	} else {
		table.DBName = inflector.Pluralize(table.Name)
	}
	if t.GoName != nil {
		table.GoName = *t.GoName
	} else {
		table.GoName = snaker.SnakeToCamel(table.Name)
	}

	for _, c := range t.Columns {
		column := Column{Name: c.Name}
		if c.DBName != nil {
			column.DBName = *c.DBName
		} else {
			column.DBName = column.Name
		}
		if c.GoName != nil {
			column.GoName = *c.GoName
		} else {
			column.GoName = snaker.SnakeToCamel(column.Name)
		}

		if c.Nullable != nil {
			column.Nullable = *c.Nullable
		}

		if c.Refer != nil {
			if c.Type != nil {
				return nil, fmt.Errorf(`%s.%s: refer and type are cannot be set together`, t.Name, c.Name)
			}
			names := strings.SplitN(*c.Refer, ".", 2)
			if len(names) != 2 {
				return nil, fmt.Errorf(`%s.%s: refer attribute must be "Table.Column"`, t.Name, c.Name)
			}
			column.Type = Reference
			column.Refer = ColumnName{Table: names[0], Column: names[1]}
		} else {
			if c.Type == nil {
				return nil, fmt.Errorf(`%s.%s: no type is specified`, t.Name, c.Name)
			}
			typ, ok := columnTypeNames[*c.Type]
			if !ok {
				return nil, fmt.Errorf(`%s.%s: type %s is not found`, t.Name, c.Name, *c.Type)
			}
			column.Type = typ
		}

		attr, ok := columnTypes[column.Type]
		if !ok {
			panic(fmt.Errorf(`no such type %s`, column.Type))
		}
		if attr.Length {
			if c.Length != nil {
				column.Length = *c.Length
			} else {
				column.Length = 255
			}
		} else {
			if c.Length != nil {
				return nil, fmt.Errorf(`%s.%s: do not specify length for type %s`, t.Name, c.Name, column.Type)
			}
		}

		table.Columns = append(table.Columns, &column)
	}

	return
}

// ResolveReference resolve reference to set type
func ResolveReference(tables []*Table) (err error) {
	columns := make(map[ColumnName]*Column)
	for _, t := range tables {
		for _, c := range t.Columns {
			columns[ColumnName{Table: t.Name, Column: c.Name}] = c
		}
	}

	for _, t := range tables {
		for _, c := range t.Columns {
			if c.Type == Reference {
				err = resolve(c, []*Column{c}, t.Name, columns)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func resolve(c *Column, stack []*Column, tableName string, columns map[ColumnName]*Column) (err error) {
	ref, ok := columns[c.Refer]
	if !ok {
		return fmt.Errorf(`%s.%s: cannot resolve reference %s.%s`, tableName, c.Name, c.Refer.Table, c.Refer.Column)
	}

	for _, s := range stack {
		if s == ref {
			return fmt.Errorf(`%s.%s: infinite reference loop is detected`, tableName, c.Name)
		}
	}

	if ref.Type == Reference {
		err = resolve(ref, append(stack, ref), tableName, columns)
		if err != nil {
			return
		}
	}

	c.Type = ref.Type
	c.Length = ref.Length
	return
}
