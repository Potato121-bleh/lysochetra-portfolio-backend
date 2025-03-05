package sqlbuilder

import "strings"

type selectSqlbuilder struct {
	// sqlStatement string
	col        string
	tbName     string
	identifier string
}

func (s *selectSqlbuilder) AddColumn(colArr []string) sqlBuilderI {
	if len(colArr) == 0 {
		return s
	}
	var newStatement string
	for i := 0; i < len(colArr); i++ {
		newStatement += " " + colArr[i] + " , "
	}
	tempStatement := strings.Split(newStatement, " ")
	var newStatementFormatted = strings.Join(tempStatement[:len(tempStatement)-2], " ")
	s.col = newStatementFormatted
	return s
}

func (s *selectSqlbuilder) AddTable(tbName string) sqlBuilderI {
	if tbName == "" {
		return nil
	}
	s.tbName = tbName
	return s
}

func (s *selectSqlbuilder) AddIdentifier(identifier string) sqlBuilderI {
	if identifier != "" {
		s.identifier = identifier
		return s
	}
	return nil
}

func (s *selectSqlbuilder) AddValue(valArr []string) sqlBuilderI {
	return s
}

func (s *selectSqlbuilder) Build() string {
	sqlStatement := "SELECT"
	if s.col == "" {
		sqlStatement += " * "
	} else {
		sqlStatement += " " + s.col + " "
	}
	if s.tbName == "" {
		return ""
	}
	sqlStatement += "FROM " + s.tbName + " "
	if s.identifier == "" {
		return sqlStatement
	} else {
		sqlStatement += "WHERE " + s.identifier + " = $1"
	}
	return sqlStatement
}
