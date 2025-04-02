package sqlbuilder

type DeleteSqlBuilder struct {
	tbName     string
	identifier string
}

func (s *DeleteSqlBuilder) AddColumn(colArr []string) SqlBuilderI {
	return s
}

func (s *DeleteSqlBuilder) AddValue(valArr []string) SqlBuilderI {
	return s
}

func (s *DeleteSqlBuilder) AddTable(tbName string) SqlBuilderI {
	s.tbName = tbName
	return s
}

func (s *DeleteSqlBuilder) AddIdentifier(identifier string) SqlBuilderI {
	s.identifier = identifier
	return s
}

func (s *DeleteSqlBuilder) Build() string {
	prepStatement := "DELETE FROM"
	if s.tbName == "" || s.identifier == "" {
		return ""
	}

	prepStatement += " " + s.tbName + " WHERE " + s.identifier + " = $1"

	return prepStatement
}
