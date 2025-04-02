package sqlbuilder

import (
	"strconv"
	"strings"
)

type UpdateSqlBuilder struct {
	colArr     []string
	col        string
	identifier string
	tbName     string
}

func (s *UpdateSqlBuilder) AddTable(tbName string) SqlBuilderI {
	s.tbName = tbName
	return s
}

func (s *UpdateSqlBuilder) AddColumn(colArr []string) SqlBuilderI {
	if len(colArr) == 0 {
		return s
	}

	prepStatement := ""

	for i := 1; i <= len(colArr); i++ {
		prepStatement += " " + colArr[i-1] + " = $" + strconv.Itoa(i) + " ,"
	}

	prepStatementArr := strings.Split(prepStatement, " ")
	prepStatement = strings.Join(prepStatementArr[:len(prepStatementArr)-1], " ")

	s.col = prepStatement
	s.colArr = colArr

	return s
}

func (s *UpdateSqlBuilder) AddIdentifier(identifier string) SqlBuilderI {
	s.identifier = identifier
	return s
}

func (s *UpdateSqlBuilder) Build() string {
	prepStatement := "UPDATE "
	if s.tbName == "" || s.col == "" {
		return ""
	}

	prepStatement += s.tbName + " SET" + s.col

	if s.identifier != "" {
		prepStatement += " WHERE " + s.identifier + " = $" + strconv.Itoa(len(s.colArr)+1)
	}

	return prepStatement
}
