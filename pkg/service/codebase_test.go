package service

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/epam/edp-perf-operator/v2/pkg/util/common"
	"github.com/epam/edp-reconciler/v2/pkg/model/codebase"
	"github.com/epam/edp-reconciler/v2/pkg/repository"
	js "github.com/epam/edp-reconciler/v2/pkg/repository/jenkins-slave"
	jp "github.com/epam/edp-reconciler/v2/pkg/repository/job-provisioning"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestSetJobProvisioningId_JobProvisionerIsNilShouldBeExecutedSuccessfully(t *testing.T) {
	schema := "public"

	c := &codebase.Codebase{}

	err := setJobProvisioningId(nil, c, schema)
	assert.NoError(t, err)
	assert.Nil(t, c.JobProvisioningId)
}

func TestSetJobProvisioningId_ShouldBeExecutedSuccessfully(t *testing.T) {
	db, mock := newMock()

	schema := "public"
	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(jp.SelectJobProvisioningSql, schema)).
		ExpectQuery().
		WithArgs("default", "ci").
		WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow(1))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := &codebase.Codebase{
		JobProvisioning: common.GetStringP("default"),
	}

	err = setJobProvisioningId(tx, c, schema)
	assert.NoError(t, err)
	assert.Equal(t, 1, *c.JobProvisioningId)
}

func TestSetJobProvisioningId_ShouldReturnAnError(t *testing.T) {
	db, mock := newMock()

	schema := "public"
	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(jp.SelectJobProvisioningSql, schema)).
		ExpectQuery().
		WithArgs("default", "ci").
		WillReturnError(errors.New("error"))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := &codebase.Codebase{
		JobProvisioning: common.GetStringP("default"),
	}

	err = setJobProvisioningId(tx, c, schema)
	assert.Error(t, err)
}

func TestSetJenkinsSlaveId_JenkinsSlaveIsNilShouldBeExecutedSuccessfully(t *testing.T) {
	schema := "public"

	c := &codebase.Codebase{}

	err := setJenkinsSlaveId(nil, c, schema)
	assert.NoError(t, err)
	assert.Nil(t, c.JenkinsSlaveId)
}

func TestSetJenkinsSlaveId_ShouldBeExecutedSuccessfully(t *testing.T) {
	db, mock := newMock()

	schema := "public"
	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(js.SelectJenkinsSlaveSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow(1))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := &codebase.Codebase{
		JenkinsSlave: common.GetStringP("default"),
	}

	err = setJenkinsSlaveId(tx, c, schema)
	assert.NoError(t, err)
	assert.Equal(t, 1, *c.JenkinsSlaveId)
}

func TestSetJenkinsSlaveId_ShouldReturnAnError(t *testing.T) {
	db, mock := newMock()

	schema := "public"
	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(js.SelectJenkinsSlaveSql, schema)).
		ExpectQuery().
		WithArgs("default", "ci").
		WillReturnError(errors.New("error"))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := &codebase.Codebase{
		JenkinsSlave: common.GetStringP("default"),
	}

	err = setJenkinsSlaveId(tx, c, schema)
	assert.Error(t, err)
}

func TestUpdateCodebase_SelectJenkinsSlaveShouldReturnAnError(t *testing.T) {
	schema := "public"
	db, mock := newMock()

	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(js.SelectJenkinsSlaveSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := codebase.Codebase{
		JenkinsSlave: common.GetStringP("default"),
	}

	err = updateCodebase(tx, c, schema)
	assert.Error(t, err)
}

func TestUpdateCodebase_SelectJobProvisioningShouldReturnAnError(t *testing.T) {
	schema := "public"
	db, mock := newMock()

	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(js.SelectJenkinsSlaveSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(1))
	mock.ExpectPrepare(fmt.Sprintf(jp.SelectJobProvisioningSql, schema)).
		ExpectQuery().
		WithArgs("default", "ci").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := codebase.Codebase{
		JenkinsSlave:    common.GetStringP("default"),
		JobProvisioning: common.GetStringP("default"),
	}

	err = updateCodebase(tx, c, schema)
	assert.Error(t, err)
}

func TestCreateCodebase_SelectGitServerShouldReturnAnError(t *testing.T) {
	schema := "public"
	db, mock := newMock()

	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(repository.SelectGitServerSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := codebase.Codebase{
		GitServer:    "default",
		JenkinsSlave: common.GetStringP("default"),
	}

	_, err = CodebaseService{}.createBE(tx, c, schema)
	assert.Error(t, err)
}

func TestCreateCodebase_SelectJobProvisioningShouldReturnAnError1(t *testing.T) {
	schema := "public"
	db, mock := newMock()

	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(repository.SelectGitServerSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectPrepare(fmt.Sprintf(js.SelectJenkinsSlaveSql, schema)).
		ExpectQuery().
		WithArgs("default").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(1))
	mock.ExpectPrepare(fmt.Sprintf(jp.SelectJobProvisioningSql, schema)).
		ExpectQuery().
		WithArgs("default", "ci").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	c := codebase.Codebase{
		GitServer:    "default",
		JenkinsSlave: common.GetStringP("default"),
	}

	_, err = CodebaseService{}.createBE(tx, c, schema)
	assert.Error(t, err)
}
