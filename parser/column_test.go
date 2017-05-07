package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutCreateTableQuery(t *testing.T) {
	var docHoge = `<table name="hoge">
<column name="id" type="int" nullable="false"/>
<column name="name" type="string" nullable="true"/>
</table>
`
	hoge, err := Parse([]byte(docHoge))
	if err != nil {
		t.Fatal(err)
	}

	query := hoge.PutCreateTableQuery()

	assert.Equal(t, "CREATE TABLE `hoges` (\n  `id` INT NOT NULL,\n  `name` VARCHAR(255)\n);", query)
}
