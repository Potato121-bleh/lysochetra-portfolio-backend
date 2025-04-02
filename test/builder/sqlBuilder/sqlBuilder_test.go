package builderTest

import (
	sqlbuilder "profile-portfolio/internal/builder/sqlBuilder"
	"testing"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlBuilderConstructor(t *testing.T) {
	builderNameArr := []string{"select", "insert", "update", "delete", ""}
	for i := range builderNameArr {
		builderName := builderNameArr[i]
		expectedBuilder := sqlbuilder.NewSqlBuilder(builderName)

		switch builderName {
		case "select":
			require.IsType(t, &sqlbuilder.SelectSqlbuilder{}, expectedBuilder)
		case "insert":
			require.IsType(t, &sqlbuilder.InsertSqlBuilder{}, expectedBuilder)
		case "update":
			require.IsType(t, &sqlbuilder.UpdateSqlBuilder{}, expectedBuilder)
		case "delete":
			require.IsType(t, &sqlbuilder.DeleteSqlBuilder{}, expectedBuilder)
		case "":
			require.Nil(t, expectedBuilder)
		}
	}
}

// -------------------------------------

func TestSelectBuilder(t *testing.T) {

	for _, tc := range SelectBuilderTestCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				selectBuilder := sqlbuilder.NewSqlBuilder("select")
				resultStmt := tc.setup(selectBuilder)
				assert.Equal(t, tc.expectedRs, resultStmt)
			},
		)
	}

}

func TestInsertBuilder(t *testing.T) {

	for _, tc := range InsertBuilderTestCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				insertBuilder := sqlbuilder.NewSqlBuilder("insert")
				testRs := tc.setup(insertBuilder)
				assert.Equal(t, tc.expectedRs, testRs)
			},
		)
	}

}

func TestUpdateBuilder(t *testing.T) {

	for _, tc := range UpdateBuilderTestCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				updateBuilder := sqlbuilder.NewSqlBuilder("update")
				testRs := tc.setup(updateBuilder)
				assert.Equal(t, tc.expectedRs, testRs)
			},
		)
	}

}

func TestDeleteBuilder(t *testing.T) {

	for _, tc := range DeleteBuilderTestCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				deleteBuilder := sqlbuilder.NewSqlBuilder("delete")
				testRs := tc.setup(deleteBuilder)
				assert.Equal(t, tc.expectedRs, testRs)
			},
		)
	}

}
