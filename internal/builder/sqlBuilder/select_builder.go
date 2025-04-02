package sqlbuilder

import "strings"

type SelectSqlbuilder struct {
	col        string
	tbName     string
	identifier string
}

func (s *SelectSqlbuilder) AddColumn(colArr []string) SqlBuilderI {
	if len(colArr) == 0 {
		return s
	}
	var newStatement string
	for i := 0; i < len(colArr); i++ {
		newStatement += "" + colArr[i] + " , "
	}
	tempStatement := strings.Split(newStatement, " ")
	var newStatementFormatted = strings.Join(tempStatement[:len(tempStatement)-2], " ")
	s.col = newStatementFormatted
	return s
}

func (s *SelectSqlbuilder) AddTable(tbName string) SqlBuilderI {
	if tbName == "" {
		return nil
	}
	s.tbName = tbName
	return s
}

func (s *SelectSqlbuilder) AddIdentifier(identifier string) SqlBuilderI {
	if identifier != "" {
		s.identifier = identifier
		return s
	}
	return nil
}

func (s *SelectSqlbuilder) AddValue(valArr []string) SqlBuilderI {
	return s
}

func (s *SelectSqlbuilder) Build() string {
	sqlStatement := "SELECT"
	if s.col == "" {
		sqlStatement += " * "
	} else {
		sqlStatement += " " + s.col + " "
	}
	if s.tbName == "" {
		return ""
	}
	sqlStatement += "FROM " + s.tbName
	if s.identifier == "" {
		return sqlStatement
	} else {
		sqlStatement += " WHERE " + s.identifier + " = $1"
	}
	return sqlStatement
}
