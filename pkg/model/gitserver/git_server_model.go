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

package gitserver

import (
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/epam/edp-reconciler/v2/pkg/model"
)

var log = ctrl.Log.WithName("git-server-model")

type GitServer struct {
	GitHost                  string
	GitUser                  string
	HttpsPort                int32
	SshPort                  int32
	PrivateSshKey            string
	CreateCodeReviewPipeline bool
	ActionLog                model.ActionLog
	Tenant                   string
	Name                     string
}

func ConvertToGitServer(k8sObj codeBaseApi.GitServer, edpName string) (*GitServer, error) {
	log.Info("Start converting GitServer", "data", k8sObj.Name)

	spec := k8sObj.Spec

	actionLog := convertGitServerActionLog(k8sObj.Status)

	gitServer := GitServer{
		GitHost:                  spec.GitHost,
		GitUser:                  spec.GitUser,
		HttpsPort:                spec.HttpsPort,
		SshPort:                  spec.SshPort,
		PrivateSshKey:            spec.NameSshKeySecret,
		CreateCodeReviewPipeline: spec.CreateCodeReviewPipeline,
		ActionLog:                *actionLog,
		Tenant:                   edpName,
		Name:                     k8sObj.Name,
	}

	return &gitServer, nil
}

func convertGitServerActionLog(status codeBaseApi.GitServerStatus) *model.ActionLog {
	return &model.ActionLog{
		Event:           model.FormatStatus(status.Status),
		DetailedMessage: status.DetailedMessage,
		Username:        status.Username,
		UpdatedAt:       status.LastTimeUpdated.Time,
		Action:          status.Action,
		Result:          status.Result,
	}
}
