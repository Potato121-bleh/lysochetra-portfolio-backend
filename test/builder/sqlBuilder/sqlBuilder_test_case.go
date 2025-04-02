package builderTest

import (
	sqlbuilder "profile-portfolio/internal/builder/sqlBuilder"
	"strings"
)

type SqlTestCase struct {
	name       string
	setup      func(builder sqlbuilder.SqlBuilderI) string
	expectedRs string
}

var SelectBuilderTestCases = []SqlTestCase{
	{
		name: "Classic SELECT query, no extra WHERE Clause",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").Build())
		},
		expectedRs: "SELECT * FROM mytb",
	}, {
		name: "Classic SELECT query, With extra WHERE Clause",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddColumn([]string{"userid"}).AddTable("mytb").AddIdentifier("username").Build())
		},
		expectedRs: "SELECT userid FROM mytb WHERE username = $1",
	},
	{
		name: "Doesn't include table name, Expecting Empty string",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddColumn([]string{"userid"}).AddIdentifier("username").Build())
		},
		expectedRs: "",
	},
}

var InsertBuilderTestCases = []SqlTestCase{
	{
		name: "Basic INSERT transaction with some column",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").AddColumn([]string{"userid", "myname"}).Build())
		},
		expectedRs: "INSERT INTO mytb ( userid , myname ) VALUES ( $1 , $2 )",
	},
	{
		name: "Basic INSERT transaction with no table, expected empty string retured",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddColumn([]string{"userid"}).Build())
		},
		expectedRs: "",
	},
	{
		name: "Basic INSERT transaction with no column, expected empty string retured",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").Build())
		},
		expectedRs: "",
	},
}

var UpdateBuilderTestCases = []SqlTestCase{
	{
		name: "Basic UPDATE transaction with some column",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").AddColumn([]string{"user_id", "user_name", "price"}).Build())
		},
		expectedRs: "UPDATE mytb SET user_id = $1 , user_name = $2 , price = $3",
	},
	{
		name: "Basic UPDATE transaction with WHERE Clause",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").AddColumn([]string{"user_name"}).AddIdentifier("user_id").Build())
		},
		expectedRs: "UPDATE mytb SET user_name = $1 WHERE user_id = $2",
	},
	{
		name: "Basic INSERT transaction with no table, expected empty string retured",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddColumn([]string{"user_name"}).AddIdentifier("user_id").Build())
		},
		expectedRs: "",
	},
	{
		name: "Basic INSERT transaction with no column, expected empty string retured",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").AddIdentifier("user_id").Build())
		},
		expectedRs: "",
	},
}

var DeleteBuilderTestCases = []SqlTestCase{
	{
		name: "Basic DELETE transaction with WHERE Clause",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").AddIdentifier("user_id").Build())
		},
		expectedRs: "DELETE FROM mytb WHERE user_id = $1",
	},
	{
		name: "Basic DELETE transaction without WHERE Clause, expect empty string",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddTable("mytb").Build())
		},
		expectedRs: "",
	},
	{
		name: "Basic DELETE transaction with no table, expect empty string",
		setup: func(builder sqlbuilder.SqlBuilderI) string {
			return strings.TrimSpace(builder.AddIdentifier("user_id").Build())
		},
		expectedRs: "",
	},
}
