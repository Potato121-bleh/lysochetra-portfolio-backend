package sqlbuilder

import (
	"strconv"
	"strings"
)

type InsertSqlBuilder struct {
	colArr     []string
	col        string
	tbName     string
	identifier string
}

func (s *InsertSqlBuilder) AddTable(tbName string) SqlBuilderI {
	s.tbName = tbName
	return s
}

func (s *InsertSqlBuilder) AddColumn(colArr []string) SqlBuilderI {
	prepStatement := " ("
	if colArr == nil {
		return nil
	}
	s.colArr = colArr
	for i := 0; i < len(colArr); i++ {
		prepStatement += " " + colArr[i] + " ,"
	}
	prepStatementArr := strings.Split(prepStatement, " ")
	prepStatementformatted := strings.Join(
		prepStatementArr[:len(prepStatementArr)-1],
		" ")
	prepStatementformatted += " ) "
	s.col = prepStatementformatted
	return s
}

func (s *InsertSqlBuilder) AddIdentifier(identifier string) SqlBuilderI {
	return s
}

func (s *InsertSqlBuilder) Build() string {
	prepStatement := "INSERT INTO "
	if s.tbName == "" || s.col == "" {
		return ""
	}
	prepStatement += s.tbName + s.col + "VALUES ("
	for i := 1; i <= len(s.colArr); i++ {
		prepStatement += " $" + strconv.Itoa(i) + " ,"
	}

	prepStatementArr := strings.Split(prepStatement, " ")
	prepStatement = strings.Join(prepStatementArr[:len(prepStatementArr)-1], " ")
	prepStatement += " )"

	return prepStatement

}
