package parser

// ColumnType is enum of column type
type ColumnType int

// Definitions of ColumnType
const (
	Char ColumnType = iota
	Varchar
	Binary
	VarBinary
	Text
	Bool
	Int8
	Int16
	Int32
	Int64
	Uint8
	Uint16
	Uint32
	Uint64
	Float
	Reference
)

var columnTypeNames = map[string]ColumnType{
	"char":      Char,
	"varchar":   Varchar,
	"binary":    Binary,
	"varbinary": VarBinary,
	"text":      Text,
	"bool":      Bool,
	"int8":      Int8,
	"int16":     Int16,
	"int32":     Int32,
	"int64":     Int64,
	"uint8":     Uint8,
	"uint16":    Uint16,
	"uint32":    Uint32,
	"uint64":    Uint64,
	"float":     Float,

	// Alias
	"string": Varchar,
	"int":    Int32,
	"uint":   Uint32,
}

type columnTypeAttr struct {
	TypeName string
	Length   bool
	Unsigned bool
}

var columnTypes = map[ColumnType]columnTypeAttr{
	Char:      columnTypeAttr{TypeName: "CHAR", Length: true},
	Varchar:   columnTypeAttr{TypeName: "VARCHAR", Length: true},
	Binary:    columnTypeAttr{TypeName: "BINARY", Length: true},
	VarBinary: columnTypeAttr{TypeName: "VARBINARY", Length: true},
	Text:      columnTypeAttr{TypeName: "TEXT"},
	Bool:      columnTypeAttr{TypeName: "TINYINT"},
	Int8:      columnTypeAttr{TypeName: "TINYINT"},
	Int16:     columnTypeAttr{TypeName: "SMALLINT"},
	Int32:     columnTypeAttr{TypeName: "INT"},
	Int64:     columnTypeAttr{TypeName: "BIGINT"},
	Uint8:     columnTypeAttr{TypeName: "TINYINT", Unsigned: true},
	Uint16:    columnTypeAttr{TypeName: "SMALLINT", Unsigned: true},
	Uint32:    columnTypeAttr{TypeName: "INT", Unsigned: true},
	Uint64:    columnTypeAttr{TypeName: "BIGINT", Unsigned: true},
	Float:     columnTypeAttr{TypeName: "FLOAT"},
	Reference: columnTypeAttr{},
}
