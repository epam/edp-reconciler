package jenkins_slave

import (
	"database/sql"
	jenkinsV2Api "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	"github.com/epam/edp-reconciler/v2/pkg/repository/jenkins-slave"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("jenkins-slave-service")

type JenkinsSlaveService struct {
	DB *sql.DB
}

func (s JenkinsSlaveService) CreateSlavesOrDoNothing(slaves []jenkinsV2Api.Slave, schemaName string) error {
	log.Info("Start executing CreateSlavesOrDoNothing method... ")

	txn, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	for _, s := range slaves {
		if len(s.Name) == 0 {
			continue
		}

		id, err := jenkins_slave.SelectJenkinsSlave(txn, s.Name, schemaName)
		if err != nil {
			return err
		}

		if id != nil {
			log.Info("Jenkins Slave already exists. Skip adding into db", "name", s)
			continue
		}

		if err := jenkins_slave.CreateJenkinsSlave(txn, s.Name, schemaName); err != nil {
			return err
		}
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	log.Info("End executing CreateSlavesOrDoNothing method... ")

	return err
}
