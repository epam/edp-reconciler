package repository

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/epam/edp-reconciler/v2/pkg/model/codebase"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestUpdate(t *testing.T) {
	db, mock := NewMock()
	c := codebase.Codebase{
		Name:        "test",
		Description: "test",
	}

	schema := "public"

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(fmt.Sprintf(updateCodebase, schema)))
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(updateCodebase, schema))).WithArgs(c.Type, strings.ToLower(c.Language), c.Framework,
		strings.ToLower(c.BuildTool), strings.ToLower(c.Strategy), c.RepositoryUrl, c.RouteSite, c.RoutePath,
		c.Status, c.TestReportFramework, c.Description,
		getIntOrNil(c.GitServerId), getStringOrNil(c.GitUrlPath), getIntOrNil(c.JenkinsSlaveId),
		getIntOrNil(c.JobProvisioningId), c.DeploymentScript, getStatus(c.Strategy), c.VersioningType,
		c.StartVersioningFrom, getIntOrNil(c.JiraServerId), getStringOrNil(c.CommitMessagePattern),
		getStringOrNil(c.TicketNamePattern), c.CiTool, getPerfIdOrNil(c.Perf), c.DefaultBranch,
		getStringOrNil(c.JiraIssueMetadataPayload), c.EmptyProject, c.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	if err := Update(*tx, c, schema); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
