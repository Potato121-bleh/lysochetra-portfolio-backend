package sqlbuilder

type SqlBuilderI interface {
	AddColumn(colArr []string) SqlBuilderI
	AddTable(tbName string) SqlBuilderI
	AddIdentifier(identifier string) SqlBuilderI
	Build() string
}

func NewSqlBuilder(builderName string) SqlBuilderI {
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
