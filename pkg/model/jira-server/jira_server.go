package jira_server

import codeBaseApi "github.com/epam/edp-codebase-operator/v2/pkg/apis/edp/v1"

type JiraServer struct {
	Name      string
	Available bool
	Tenant    string
}

func ConvertSpecToJira(jira codeBaseApi.JiraServer, tenant string) JiraServer {
	return JiraServer{
		Name:      jira.Name,
		Available: jira.Status.Available,
		Tenant:    tenant,
	}
}
