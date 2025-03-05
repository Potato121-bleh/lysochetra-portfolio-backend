package sqlbuilder

type deleteSqlBuilder struct {
	tbName     string
	identifier string
}

func (s *deleteSqlBuilder) AddColumn(colArr []string) sqlBuilderI {
	return s
}

func (s *deleteSqlBuilder) AddValue(valArr []string) sqlBuilderI {
	return s
}

func (s *deleteSqlBuilder) AddTable(tbName string) sqlBuilderI {
	s.tbName = tbName
	return s
}

func (s *deleteSqlBuilder) AddIdentifier(identifier string) sqlBuilderI {
	s.identifier = identifier
	return s
}

func (s *deleteSqlBuilder) Build() string {
	prepStatement := "DELETE FROM"
	if s.tbName == "" || s.identifier == "" {
		return ""
	}

	prepStatement += " " + s.tbName + " WHERE " + s.identifier + " = $1"

	return prepStatement
}
