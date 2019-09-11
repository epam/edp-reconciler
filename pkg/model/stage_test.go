package model

import (
	"fmt"
	edpv1alpha1 "github.com/epmd-edp/reconciler/v2/pkg/apis/edp/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

const (
	qualityGate     = "autotests"
	cdPipelineName  = "fake-name"
	jenkinsStepName = "fake-jenkins-step-name"
	fakeDecription  = "fake-description"
	triggerType     = "manual"
	stageAction     = "accept_cd_stage_registration"
)

func TestConvertMethodToCDStage(t *testing.T) {
	timeNow := time.Now()
	branchName := "fake-branch"
	autotestName := "fake-autotest"

	k8sObj := edpv1alpha1.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: edpv1alpha1.StageSpec{
			Name:        name,
			CdPipeline:  cdPipelineName,
			Description: fakeDecription,
			TriggerType: triggerType,
			Order:       1,
			QualityGates: []edpv1alpha1.QualityGate{
				{
					QualityGateType: qualityGate,
					BranchName:      &branchName,
					AutotestName:    &autotestName,
					StepName:        jenkinsStepName,
				},
			},
		},
		Status: edpv1alpha1.StageStatus{
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

	cdStage, err := ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	if cdStage.Name != name {
		t.Fatal(fmt.Sprintf("name is not %v", name))
	}

	if cdStage.CdPipelineName != cdPipelineName {
		t.Fatal(fmt.Sprintf("cdPipelineName is not %v", cdPipelineName))
	}

	if cdStage.Description != fakeDecription {
		t.Fatal(fmt.Sprintf("fakeDecription is not %v", fakeDecription))
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

	if cdStage.ActionLog.Event != formatStatus(event) {
		t.Fatal(fmt.Sprintf("event has incorrect status %v", event))
	}

	if cdStage.ActionLog.DetailedMessage != detailedMessage {
		t.Fatal(fmt.Sprintf("detailed message is incorrect %v", detailedMessage))
	}

	if cdStage.ActionLog.Username != username {
		t.Fatal(fmt.Sprintf("username is incorrect %v", username))
	}

	if !cdStage.ActionLog.UpdatedAt.Equal(timeNow) {
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
		acceptCdStageRegistration     = "accept_cd_stage_registration"
		fetchingUserSettingsConfigMap = "fetching_user_settings_config_map"
		openshiftProjectCreation      = "openshift_project_creation"
		jenkinsConfiguration          = "jenkins_configuration"
		setupDeploymentTemplates      = "setup_deployment_templates"
		nonExistedAction              = "fake-action"

		acceptCdStageRegistrationMsg     = "Accept CD Stage %v registration"
		fetchingUserSettingsConfigMapMsg = "Fetch User Settings from config map during CD Stage %v provision"
		openshiftProjectCreationMsg      = "Creation of Openshift Project %v"
		jenkinsConfigurationMsg          = "CI Jenkins pipelines codebase %v provisioning"
		setupDeploymentTemplatesMsg      = "Setup deployment templates for cd_stage %v"
		nonExistedActionMsg              = "fake message"

		branchName   = "fake-branch"
		autotestName = "fake-autotest"
	)

	k8sObj := edpv1alpha1.Stage{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: edpv1alpha1.StageSpec{
			Name:        name,
			CdPipeline:  cdPipelineName,
			Description: fakeDecription,
			TriggerType: triggerType,
			Order:       1,
			QualityGates: []edpv1alpha1.QualityGate{
				{
					QualityGateType: qualityGate,
					BranchName:      &branchName,
					AutotestName:    &autotestName,
					StepName:        jenkinsStepName,
				},
			},
		},
		Status: edpv1alpha1.StageStatus{
			Username:        username,
			DetailedMessage: detailedMessage,
			Value:           "active",
			Action:          acceptCdStageRegistration,
			Result:          result,
			Available:       true,
			LastTimeUpdated: time.Now(),
			Status:          event,
		},
	}

	cdStage, err := ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(acceptCdStageRegistrationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = fetchingUserSettingsConfigMap
	cdStage, err = ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(fetchingUserSettingsConfigMapMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = openshiftProjectCreation
	cdStage, err = ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(openshiftProjectCreationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = jenkinsConfiguration
	cdStage, err = ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(jenkinsConfigurationMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = setupDeploymentTemplates
	cdStage, err = ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(setupDeploymentTemplatesMsg, name), cdStage.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdStage.ActionLog.ActionMessage))

	k8sObj.Status.Action = nonExistedAction
	cdStage, err = ConvertToStage(k8sObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, nonExistedActionMsg, cdStage.ActionLog.ActionMessage)

}