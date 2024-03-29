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

package cdpipeline

import (
	"fmt"
	cdPipeApi "github.com/epam/edp-cd-pipeline-operator/v2/pkg/apis/edp/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/epam/edp-reconciler/v2/pkg/model"
)

const (
	name                  = "fake-name"
	username              = "fake-user"
	detailedMessage       = "fake-detailed-message"
	inputDockerStream     = "fake-docker-stream-verified"
	applicationsToPromote = "fake-application"
	result                = "success"
	cdPipelineAction      = "setup_initial_structure"
	event                 = "created"
	edpName               = "foobar"
)

func TestConvertMethodToCDPipeline(t *testing.T) {
	timeNow := metav1.Now()

	k8sObj := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:                  name,
			InputDockerStreams:    []string{inputDockerStream},
			ApplicationsToPromote: []string{applicationsToPromote},
		},
		Status: cdPipeApi.CDPipelineStatus{
			Username:        username,
			DetailedMessage: detailedMessage,
			Value:           "active",
			Action:          cdPipelineAction,
			Result:          result,
			Available:       true,
			LastTimeUpdated: timeNow,
			Status:          event,
		},
	}

	cdPipeline, err := ConvertToCDPipeline(k8sObj, edpName)
	if err != nil {
		t.Fatal(err)
	}

	if cdPipeline.Name != name {
		t.Fatal(fmt.Sprintf("name is not %v", name))
	}

	checkSpecField(t, cdPipeline.InputDockerStreams, inputDockerStream, "input docker stream")

	checkSpecField(t, cdPipeline.ApplicationsToPromote, applicationsToPromote, "applications to promote")

	if cdPipeline.ActionLog.Event != model.FormatStatus(event) {
		t.Fatal(fmt.Sprintf("event has incorrect status %v", event))
	}

	if cdPipeline.ActionLog.DetailedMessage != detailedMessage {
		t.Fatal(fmt.Sprintf("detailed message is incorrect %v", detailedMessage))
	}

	if cdPipeline.ActionLog.Username != username {
		t.Fatal(fmt.Sprintf("username is incorrect %v", username))
	}

	if !cdPipeline.ActionLog.UpdatedAt.Equal(timeNow.Time) {
		t.Fatal(fmt.Sprintf("'updated at' is incorrect %v", username))
	}

	if cdPipeline.ActionLog.Action != cdPipelineAction {
		t.Fatal(fmt.Sprintf("action is incorrect %v", cdPipelineAction))
	}

	if cdPipeline.ActionLog.Result != result {
		t.Fatal(fmt.Sprintf("result is incorrect %v", result))
	}

	actionMessage := fmt.Sprintf(cdPipelineActionMessageMap[cdPipelineAction], name)
	if cdPipeline.ActionLog.ActionMessage != actionMessage {
		t.Fatal(fmt.Sprintf("action message is incorrect %v", actionMessage))
	}

}

func checkSpecField(t *testing.T, src []string, toCheck string, entityName string) {
	if len(src) != 1 {
		t.Fatal(fmt.Sprintf("%v has incorrect size", entityName))
	}

	if src[0] != toCheck {
		t.Fatal(fmt.Sprintf("%v name is not %v", entityName, toCheck))
	}
}

func TestCDPipelineActionMessages(t *testing.T) {

	var (
		acceptCdPipelineRegistrationMsg = "Accept CD Pipeline %v registration"
		jenkinsConfigurationMsg         = "CI Jenkins pipelines %v provisioning"
		setupInitialStructureMsg        = "Initial structure for CD Pipeline %v is created"
	)

	k8sObj := cdPipeApi.CDPipeline{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fake-namespace",
			Name:      name,
		},
		Spec: cdPipeApi.CDPipelineSpec{
			Name:                  name,
			InputDockerStreams:    []string{inputDockerStream},
			ApplicationsToPromote: []string{applicationsToPromote},
		},
		Status: cdPipeApi.CDPipelineStatus{
			Username:        username,
			DetailedMessage: detailedMessage,
			Value:           "active",
			Result:          result,
			Available:       true,
			LastTimeUpdated: metav1.Now(),
			Status:          event,
			Action:          "accept_cd_pipeline_registration",
		},
	}

	cdPipeline, err := ConvertToCDPipeline(k8sObj, edpName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(acceptCdPipelineRegistrationMsg, name), cdPipeline.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdPipeline.ActionLog.ActionMessage))

	k8sObj.Status.Action = "jenkins_configuration"
	cdPipeline, err = ConvertToCDPipeline(k8sObj, edpName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(jenkinsConfigurationMsg, name), cdPipeline.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdPipeline.ActionLog.ActionMessage))

	k8sObj.Status.Action = cdPipeApi.SetupInitialStructureForCDPipeline
	cdPipeline, err = ConvertToCDPipeline(k8sObj, edpName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(setupInitialStructureMsg, name), cdPipeline.ActionLog.ActionMessage,
		fmt.Sprintf("converted action is incorrect %v", cdPipeline.ActionLog.ActionMessage))

	k8sObj.Status = cdPipeApi.CDPipelineStatus{}
	_, err = ConvertToCDPipeline(k8sObj, edpName)
	if err != nil {
		t.Fatal(err)
	}
}
