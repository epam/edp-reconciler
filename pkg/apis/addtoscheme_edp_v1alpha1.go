package apis

import (
	edpv1alpha1Codebase "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1alpha1"
	pipeApi "github.com/epmd-edp/cd-pipeline-operator/v2/pkg/apis/edp/v1alpha1"
	edpComponentV1Api "github.com/epmd-edp/edp-component-operator/pkg/apis/v1/v1alpha1"
	jenkinsV2Api "github.com/epmd-edp/jenkins-operator/v2/pkg/apis/v2/v1alpha1"
	v1 "github.com/openshift/api/template/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, pipeApi.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, v1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, jenkinsV2Api.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, edpv1alpha1Codebase.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, edpComponentV1Api.SchemeBuilder.AddToScheme)
}
