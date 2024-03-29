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

package codebase

import (
	codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestConvert(t *testing.T) {
	frw := "spring-boot"
	k8sObject := codeBaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "fightclub",
			Name:      "fc-ui",
		},
		Spec: codeBaseApi.CodebaseSpec{
			Lang:      "java",
			Framework: &frw,
			BuildTool: "maven",
			Strategy:  codeBaseApi.Create,
		},
		Status: codeBaseApi.CodebaseStatus{
			Available:       true,
			LastTimeUpdated: metav1.Now(),
			Status:          "created",
		},
	}

	app, err := Convert(k8sObject, "foobar")
	if err != nil {
		t.Fatal(err)
	}

	if app.Name != "fc-ui" {
		t.Fatal("name is not fc-ui")
	}
}
