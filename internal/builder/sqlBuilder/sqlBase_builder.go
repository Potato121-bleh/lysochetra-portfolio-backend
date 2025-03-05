package sqlbuilder

type sqlBuilderI interface {
	AddColumn(colArr []string) sqlBuilderI
	AddTable(tbName string) sqlBuilderI
	AddIdentifier(identifier string) sqlBuilderI
	Build() string
}

func NewSqlBuilder(builderName string) sqlBuilderI {
	switch builderName {
	case "select":
		return &selectSqlbuilder{}
	case "insert":
		return &insertSqlBuilder{}
	case "update":
		return &updateSqlBuilder{}
	case "delete":
		return &deleteSqlBuilder{}
	default:
		return nil
	}
}
