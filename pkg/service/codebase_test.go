package service

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/epam/edp-perf-operator/v2/pkg/util/common"
	"github.com/epam/edp-reconciler/v2/pkg/model/codebase"
	jenkins_slave "github.com/epam/edp-reconciler/v2/pkg/repository/jenkins-slave"
	job_provisioning "github.com/epam/edp-reconciler/v2/pkg/repository/job-provisioning"
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
	mock.ExpectPrepare(fmt.Sprintf(job_provisioning.SelectJobProvisioningSql, schema)).
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
	mock.ExpectPrepare(fmt.Sprintf(job_provisioning.SelectJobProvisioningSql, schema)).
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
	println(c.JenkinsSlaveId)
	assert.Nil(t, c.JenkinsSlaveId)
}

func TestSetJenkinsSlaveId_ShouldBeExecutedSuccessfully(t *testing.T) {
	db, mock := newMock()

	schema := "public"
	mock.ExpectBegin()
	mock.ExpectPrepare(fmt.Sprintf(jenkins_slave.SelectJenkinsSlaveSql, schema)).
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
	mock.ExpectPrepare(fmt.Sprintf(jenkins_slave.SelectJenkinsSlaveSql, schema)).
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
