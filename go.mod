module github.com/epam/edp-reconciler/v2

go 1.14

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-sigs/reference-docs v0.0.0-20170929004150-fcf65347b256
	github.com/markbates/inflect => github.com/markbates/inflect v1.0.4
	github.com/openshift/api => github.com/openshift/api v0.0.0-20210416130433-86964261530c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	k8s.io/api => k8s.io/api v0.20.7-rc.0
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/epam/edp-cd-pipeline-operator/v2 v2.3.0-58.0.20210730114857-7c52e393d03d
	github.com/epam/edp-codebase-operator/v2 v2.3.0-95.0.20210615090821-27e73211c5f4
	github.com/epam/edp-component-operator v0.1.1-0.20210427065236-c7dce7f4ea2b
	github.com/epam/edp-gerrit-operator/v2 v2.3.0-73.0.20210427073621-8b07da9960b2
	github.com/epam/edp-jenkins-operator/v2 v2.3.0-130.0.20210611080440-fd54e8b0eda3
	github.com/epam/edp-perf-operator/v2 v2.0.0-20210615084859-4e4202d10e93
	github.com/go-logr/logr v0.4.0
	github.com/lib/pq v1.8.0
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.21.0-rc.0
	k8s.io/apimachinery v0.21.0-rc.0
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)
