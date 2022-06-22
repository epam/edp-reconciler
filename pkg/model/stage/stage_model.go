/*
 * Copyright 2019 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stage

import (
	"fmt"
	"strings"

	cpPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"

	"github.com/epam/edp-reconciler/v2/pkg/model"
)

type Stage struct {
	Id              int
	Name            string
	Tenant          string
	Namespace       string
	CdPipelineName  string
	Description     string
	TriggerType     string
	Order           int
	ActionLog       model.ActionLog
	Status          string
	QualityGates    []QualityGate
	Source          Source
	JobProvisioning string
}

type Source struct {
	Type    string
	Library Library
}

type Library struct {
	Id       *int
	BranchId *int
	Name     string
	Branch   string
}

type QualityGate struct {
	QualityGate     string
	JenkinsStepName string
	AutotestName    *string
	BranchName      *string
}

var cdStageActionMessageMap = map[string]string{
	"accept_cd_stage_registration":      "Accept CD Stage %v registration",
	"fetching_user_settings_config_map": "Fetch User Settings from config map during CD Stage %v provision",
	"platform_project_creation":         "Create Openshift Project for Stage %v",
	"jenkins_configuration":             "CI Jenkins pipelines %v provisioning",
	"setup_deployment_templates":        "Setup deployment templates for cd_stage %v",
	"create_jenkins_pipeline":           "Create Jenkins pipeline for CD Stage %v",
}

// ConvertToStage returns converted to DTO Stage object from K8S and provided edp name
// An error occurs if method received nil instead of k8s object
func ConvertToStage(k8sObject cpPipeApi.Stage, edpName string) (*Stage, error) {
	spec := k8sObject.Spec
	actionLog := convertStageActionLog(k8sObject.Name, k8sObject.Status)
	stage := Stage{
		Name:           spec.Name,
		Tenant:         edpName,
		Namespace:      k8sObject.Namespace,
		CdPipelineName: spec.CdPipeline,
		Description:    spec.Description,
		TriggerType:    strings.ToLower(spec.TriggerType),
		Order:          spec.Order,
		ActionLog:      *actionLog,
		Status:         k8sObject.Status.Value,
		QualityGates:   convertQualityGatesFromRequest(spec.QualityGates),
		Source: Source{
			Type: spec.Source.Type,
			Library: Library{
				Name:   spec.Source.Library.Name,
				Branch: spec.Source.Library.Branch,
			},
		},
		JobProvisioning: spec.JobProvisioning,
	}
	return &stage, nil
}

func convertQualityGatesFromRequest(gates []cpPipeApi.QualityGate) []QualityGate {
	var result []QualityGate

	for _, val := range gates {
		gate := QualityGate{
			QualityGate:     strings.ToLower(val.QualityGateType),
			JenkinsStepName: strings.ToLower(val.StepName),
		}

		if gate.QualityGate == "autotests" {
			gate.AutotestName = val.AutotestName
			gate.BranchName = val.BranchName
		}

		result = append(result, gate)
	}

	return result
}

func convertStageActionLog(cdStageName string, status cpPipeApi.StageStatus) *model.ActionLog {

	al := &model.ActionLog{
		Event:           model.FormatStatus(status.Status),
		DetailedMessage: status.DetailedMessage,
		Username:        status.Username,
		UpdatedAt:       status.LastTimeUpdated.Time,
		Action:          fmt.Sprint(status.Action),
		Result:          fmt.Sprint(status.Result),
	}

	if status.Result == "error" {
		al.ActionMessage = status.DetailedMessage
		return al
	}

	al.ActionMessage = fmt.Sprintf(cdStageActionMessageMap[fmt.Sprint(status.Action)], cdStageName)
	return al
}
