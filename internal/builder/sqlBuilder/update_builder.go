package sqlbuilder

import (
	"strconv"
	"strings"
)

// UPDATE tableName SET col1 = $1, col2 = $2 WHERE identifier = $3

// DELETE FROM tableName WHERE identifier

/**
update we can use map[key]val
we going to accept a normally but we do a map during build and the loop inside to fill out.
*/

type updateSqlBuilder struct {
	colArr     []string
	col        string
	identifier string
	tbName     string
}

// type sqlBuilderI interface {
// 	addColumn(colArr []string) sqlBuilderI
// 	addTable(tbName string) sqlBuilderI
// 	addIdentifier(identifier string) sqlBuilderI
// 	addValue(valArr []string) sqlBuilderI
// 	build() string
// }

func (s *updateSqlBuilder) AddTable(tbName string) sqlBuilderI {
	s.tbName = tbName
	return s
}

func (s *updateSqlBuilder) AddColumn(colArr []string) sqlBuilderI {
	if len(colArr) == 0 {
		return s
	}

	// UPDATE tableName SET col1 = $1, col2 = $2 WHERE identifier = $3

	prepStatement := ""

	for i := 1; i <= len(colArr); i++ {
		prepStatement += " " + colArr[i-1] + " = $" + strconv.Itoa(i) + " ,"
	}

	prepStatementArr := strings.Split(prepStatement, " ")
	prepStatement = strings.Join(prepStatementArr[:len(prepStatementArr)-1], " ")

	s.col = prepStatement

	return s
}

func (s *updateSqlBuilder) AddIdentifier(identifier string) sqlBuilderI {
	s.identifier = identifier
	return s
}

func (s *updateSqlBuilder) Build() string {
	prepStatement := "UPDATE "
	if s.tbName == "" || s.col == "" {
		return ""
	}

	prepStatement += s.tbName + " SET"
	for i := 0; i < len(s.colArr); i++ {
		prepStatement += " " + s.colArr[i] + " = $" + strconv.Itoa(i+1) + " ,"
	}

	prepStatementArr := strings.Split(prepStatement, " ")
	prepStatement = strings.Join(prepStatementArr[:len(prepStatementArr)-1], " ")

	if s.identifier != "" {
		prepStatement += " WHERE " + s.identifier + " = $" + strconv.Itoa(len(s.colArr)+1)
	}

	return prepStatement
}
