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
		return &SelectSqlbuilder{}
	case "insert":
		return &InsertSqlBuilder{}
	case "update":
		return &UpdateSqlBuilder{}
	case "delete":
		return &DeleteSqlBuilder{}
	default:
		return nil
	}
}
