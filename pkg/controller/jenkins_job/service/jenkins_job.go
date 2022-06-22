package service

import (
	"context"
	"database/sql"
	"fmt"

	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	jenkinsApi "github.com/epam/edp-jenkins-operator/v2/pkg/apis/v2/v1"
	"github.com/epam/edp-jenkins-operator/v2/pkg/util/consts"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/epam/edp-reconciler/v2/pkg/controller/helper"
	"github.com/epam/edp-reconciler/v2/pkg/db"
	"github.com/epam/edp-reconciler/v2/pkg/model"
	"github.com/epam/edp-reconciler/v2/pkg/repository"
	"github.com/epam/edp-reconciler/v2/pkg/util/cluster"
)

const ErrorStatus = "error"

var actions = map[string]string{
	"platform_project_creation": "Create Platform Project for Stage %v",
	"role_binding":              "Create Role Binding for project stage %v",
	"create_jenkins_pipeline":   "Create Jenkins pipeline for CD Stage %v",
}

type JenkinsJobService struct {
	DB     *sql.DB
	Client client.Client
}

var log = ctrl.Log.WithName("jenkins-job-service")

func (s JenkinsJobService) UpdateActionLog(jj *jenkinsApi.JenkinsJob) error {
	log.V(2).Info("start adding action log for jenkins job", "name", jj.Name)
	l, err := s.createActionLogModel(*jj)
	if err != nil {
		return err
	}

	edpN, err := helper.GetEDPName(s.Client, jj.Namespace)
	if err != nil {
		return errors.Wrap(err, "cannot get edp name")
	}

	stage, err := s.getStageInstanceOwner(*jj)
	if err != nil {
		return err
	}

	tx, err := db.Instance.Begin()
	if err != nil {
		return err
	}

	p, err := repository.GetCDPipeline(tx, stage.Spec.CdPipeline, *edpN)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrapf(err, "cannot get CD Pipeline %v", stage.Spec.CdPipeline)
	}

	if p == nil {
		_ = tx.Rollback()
		return fmt.Errorf("cd pipeline %v is not inserted into table yet", stage.Spec.CdPipeline)
	}

	alid, err := repository.CreateEventActionLog(tx, *l, *edpN)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = repository.CreateCDPipelineActionLog(tx, p.Id, *alid, *edpN); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.V(2).Info("action log record has been added", "name", jj.Name)
	return nil

}

func (s JenkinsJobService) createActionLogModel(jj jenkinsApi.JenkinsJob) (*model.ActionLog, error) {
	st := jj.Status
	l := &model.ActionLog{
		Username:        st.Username,
		UpdatedAt:       st.LastTimeUpdated.Time,
		Action:          fmt.Sprint(st.Action),
		Result:          fmt.Sprint(st.Result),
		DetailedMessage: st.DetailedMessage,
	}

	if st.Result == ErrorStatus {
		l.ActionMessage = st.DetailedMessage
		return l, nil
	}
	stage, err := s.getStageInstanceOwner(jj)
	if err != nil {
		return nil, err
	}
	l.ActionMessage = fmt.Sprintf(actions[string(st.Action)], stage.Name)
	return l, nil
}

func (s JenkinsJobService) getStageInstanceOwner(jj jenkinsApi.JenkinsJob) (*cdPipeApi.Stage, error) {
	log.V(2).Info("start getting stage owner cr", "stage", jj.Name)
	if ow := cluster.GetOwnerReference(consts.StageKind, jj.GetOwnerReferences()); ow != nil {
		log.V(2).Info("trying to fetch stage owner from reference", "stage", ow.Name)
		return s.getStageInstance(ow.Name, jj.Namespace)
	}
	if jj.Spec.StageName != nil {
		log.V(2).Info("trying to fetch stage owner from spec", "stage", jj.Spec.StageName)
		return s.getStageInstance(*jj.Spec.StageName, jj.Namespace)
	}
	return nil, fmt.Errorf("couldn't find stage owner for jenkins job %v", jj.Name)
}

func (s JenkinsJobService) getStageInstance(name, namespace string) (*cdPipeApi.Stage, error) {
	nsn := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	i := &cdPipeApi.Stage{}
	if err := s.Client.Get(context.TODO(), nsn, i); err != nil {
		return nil, errors.Wrapf(err, "failed to get instance by name %v", name)
	}
	return i, nil
}
