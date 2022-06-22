package model

import (
	edpCompApi "github.com/epam/edp-component-operator/pkg/apis/v1/v1"
)

type EDPComponent struct {
	Type    string
	Url     string
	Icon    string
	Visible bool
}

func ConvertToEDPComponent(k8sObj edpCompApi.EDPComponent) (*EDPComponent, error) {
	s := k8sObj.Spec

	return &EDPComponent{
		Type:    s.Type,
		Url:     s.Url,
		Icon:    s.Icon,
		Visible: s.Visible,
	}, nil
}
