package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var docHoge = `<table name="great_hoge">
<column name="id" type="int" nullable="false"/>
<column name="name" type="string" nullable="true"/>
<column name="piyo_id" refer="piyo.id"/>
</table>
`

func TestParse(t *testing.T) {
	hoge, err := Parse([]byte(docHoge))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, &Table{
		Name:   "great_hoge",
		DBName: "great_hoges",
		GoName: "GreatHoge",
		Columns: []*Column{
			&Column{
				Name:   "id",
				DBName: "id",
				GoName: "ID",
				Type:   Int32,
			},
			&Column{
				Name:     "name",
				DBName:   "name",
				GoName:   "Name",
				Type:     Varchar,
				Length:   255,
				Nullable: true,
			},
			&Column{
				Name:   "piyo_id",
				DBName: "piyo_id",
				GoName: "PiyoID",
				Type:   Reference,
				Refer:  ColumnName{Table: "piyo", Column: "id"},
			},
		},
	}, hoge)
}

func TestResolveReference(t *testing.T) {
	hoge, err := Parse([]byte(docHoge))
	if err != nil {
		t.Fatal(err)
	}

	docPiyo := `<table name="piyo">
<column name="id" type="int"/>
</table>
`

	piyo, err := Parse([]byte(docPiyo))
	if err != nil {
		t.Fatal(err)
	}

	err = ResolveReference([]*Table{hoge, piyo})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, Int32, hoge.Columns[2].Type)

	// loop
	hoge, err = Parse([]byte(docHoge))
	if err != nil {
		t.Fatal(err)
	}

	docPiyo = `<table name="piyo">
<column name="id" refer="hoge.piyo_id"/>
</table>
`

	piyo, err = Parse([]byte(docPiyo))
	if err != nil {
		t.Fatal(err)
	}

	err = ResolveReference([]*Table{hoge, piyo})
	assert.NotNil(t, err)
}
