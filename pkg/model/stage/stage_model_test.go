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
	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/epam/edp-reconciler/v2/pkg/model"
)

const (
	name            = "fake-name"
	qualityGate     = "autotests"
	cdPipelineName  = "fake-name"
	jenkinsStepName = "fake-jenkins-step-name"
	fakeDescription = "fake-description"
	triggerType     = "manual"
	stageAction     = "accept_cd_stage_registration"
	edpN            = "foobar"
	username        = "fake-user"
	detailedMessage = "fake-detailed-message"
	result          = "success"
	event           = "created"
	jobProvisioning = "fake-job-provisioning"
)

func TestConvertMethodToCDStage(t *testing.T) {
	timeNow := metav1.Now()
	branchName := "fake-branch"
	autotestName := "fake-autotest"

	k8sObj := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: cdPipeApi.StageSpec{
			Name:        name,
			CdPipeline:  cdPipelineName,
			Description: fakeDescription,
			TriggerType: triggerType,
			Order:       1,
			QualityGates: []cdPipeApi.QualityGate{
				{
					QualityGateType: qualityGate,
					BranchName:      &branchName,
					AutotestName:    &autotestName,
					StepName:        jenkinsStepName,
				},
			},
			JobProvisioning: jobProvisioning,
		},
		Status: cdPipeApi.StageStatus{
			Username:        username,
			DetailedMessage: detailedMessage,
			Value:           "active",
			Action:          stageAction,
			Result:          result,
			Available:       true,
			LastTimeUpdated: timeNow,
			Status:          event,
		},
	}

	cdStage, err := ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	if cdStage.Name != name {
		t.Fatal(fmt.Sprintf("name is not %v", name))
	}

	if cdStage.CdPipelineName != cdPipelineName {
		t.Fatal(fmt.Sprintf("cdPipelineName is not %v", cdPipelineName))
	}

	if cdStage.Description != fakeDescription {
		t.Fatal(fmt.Sprintf("fakeDescription is not %v", fakeDescription))
	}

	if cdStage.TriggerType != triggerType {
		t.Fatal(fmt.Sprintf("triggerType is not %v", triggerType))
	}

	if cdStage.Order != 1 {
		t.Fatal(fmt.Sprintf("order is not %v", 1))
	}

	if len(cdStage.QualityGates) != 1 {
		t.Fatal(fmt.Sprintf("quality gates size is not %v", 1))
	}

	if cdStage.QualityGates[0].QualityGate != qualityGate {
		t.Fatal(fmt.Sprintf("quality gate is not %v", qualityGate))
	}

	if *cdStage.QualityGates[0].BranchName != branchName {
		t.Fatal(fmt.Sprintf("branch name is not %v", branchName))
	}

	if *cdStage.QualityGates[0].AutotestName != autotestName {
		t.Fatal(fmt.Sprintf("autotest name is not %v", autotestName))
	}

	if cdStage.QualityGates[0].JenkinsStepName != jenkinsStepName {
		t.Fatal(fmt.Sprintf("jenkinsStepName is not %v", jenkinsStepName))
	}
	if cdStage.Tenant != edpN {
		t.Errorf("ConvertToStage() expected - %v, actual - %v", edpN, cdStage.Tenant)
	}
	if cdStage.ActionLog.Event != model.FormatStatus(event) {
		t.Fatal(fmt.Sprintf("event has incorrect status %v", event))
	}

	if cdStage.ActionLog.DetailedMessage != detailedMessage {
		t.Fatal(fmt.Sprintf("detailed message is incorrect %v", detailedMessage))
	}

	if cdStage.ActionLog.Username != username {
		t.Fatal(fmt.Sprintf("username is incorrect %v", username))
	}

	if !cdStage.ActionLog.UpdatedAt.Equal(timeNow.Time) {
		t.Fatal(fmt.Sprintf("'updated at' is incorrect %v", username))
	}

	if cdStage.ActionLog.Action != stageAction {
		t.Fatal(fmt.Sprintf("action is incorrect %v", stageAction))
	}

	if cdStage.ActionLog.Result != result {
		t.Fatal(fmt.Sprintf("result is incorrect %v", result))
	}

	actionMessage := fmt.Sprintf(cdStageActionMessageMap[stageAction], name)
	if cdStage.ActionLog.ActionMessage != actionMessage {
		t.Fatal(fmt.Sprintf("action message is incorrect %v", actionMessage))
	}
}

func TestCDStageActionMessages(t *testing.T) {

	var (
		acceptCdStageRegistrationMsg     = "Accept CD Stage %v registration"
		fetchingUserSettingsConfigMapMsg = "Fetch User Settings from config map during CD Stage %v provision"
		openshiftProjectCreationMsg      = "Create Openshift Project for Stage %v"
		jenkinsConfigurationMsg          = "CI Jenkins pipelines %v provisioning"
		setupDeploymentTemplatesMsg      = "Setup deployment templates for cd_stage %v"
		nonExistedActionMsg              = "fake message"

		branchName   = "fake-branch"
		autotestName = "fake-autotest"
	)

	k8sObj := cdPipeApi.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: cdPipeApi.StageSpec{
			Name:        name,
			CdPipeline:  cdPipelineName,
			Description: fakeDescription,
			TriggerType: triggerType,
			Order:       1,
			QualityGates: []cdPipeApi.QualityGate{
				{
					QualityGateType: qualityGate,
					BranchName:      &branchName,
					AutotestName:    &autotestName,
					StepName:        jenkinsStepName,
				},
			},
		},
		Status: cdPipeApi.StageStatus{
			Username:        username,
			DetailedMessage: detailedMessage,
			Value:           "active",
			Action:          cdPipeApi.AcceptCDStageRegistration,
			Result:          result,
			Available:       true,
			LastTimeUpdated: metav1.Now(),
			Status:          event,
		},
	}

	cdStage, err := ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(acceptCdStageRegistrationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = "fetching_user_settings_config_map"
	cdStage, err = ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(fetchingUserSettingsConfigMapMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = "platform_project_creation"
	cdStage, err = ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(openshiftProjectCreationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = "jenkins_configuration"
	cdStage, err = ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(jenkinsConfigurationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = "setup_deployment_templates"
	cdStage, err = ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(setupDeploymentTemplatesMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status = cdPipeApi.StageStatus{}
	cdStage, err = ConvertToStage(k8sObj, edpN)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, nonExistedActionMsg, cdStage.ActionLog.ActionMessage)

}
