[![codecov](https://codecov.io/gh/epam/edp-reconciler/branch/master/graph/badge.svg?token=NISU3J07BI)](https://codecov.io/gh/epam/edp-reconciler)

# Reconciler Operator

| :heavy_exclamation_mark: Please refer to [EDP documentation](https://epam.github.io/edp-install/) to get the notion of the main concepts and guidelines. |
| --- |

Get acquainted with the Reconciler Operator and the installation process as well as the local development.

## Overview

Reconciler Operator is an EDP operator that is responsible for saving state of CR's in EDP database. Operator installation can be applied on two container orchestration platforms: OpenShift and Kubernetes.

_**NOTE:** Operator is platform-independent, that is why there is a unified instruction for deploying._

## Prerequisites

* Linux machine or Windows Subsystem for Linux instance with [Helm 3](https://helm.sh/docs/intro/install/) installed;
* Cluster admin access to the cluster;
* EDP project/namespace is deployed by following the [Install EDP](https://epam.github.io/edp-install/operator-guide/install-edp/) instruction.

## Installation

In order to install the EDP Reconciler Operator, follow the steps below:

1. To add the Helm EPAMEDP Charts for local client, run "helm repo add":
     ```bash
     helm repo add epamedp https://epam.github.io/edp-helm-charts/stable
     ```
2. Choose available Helm chart version:
     ```bash
     helm search repo epamedp/reconciler -l
     ```
   Example response:
     ```bash
     NAME              	CHART VERSION	APP VERSION	DESCRIPTION
     epamedp/reconciler	2.11.0       	2.11.0     	A Helm chart for EDP Reconciler
     epamedp/reconciler	2.10.0       	2.10.0     	A Helm chart for EDP Reconciler
     ```

    _**NOTE:** It is highly recommended to use the latest released version._

3. Full chart parameters available in [deploy-templates/README.md](deploy-templates/README.md).

4. Install operator in the <edp-project> namespace with the helm command; find below the installation command example:
    ```bash
    helm install reconciler epamedp/reconciler --namespace <edp-project> --version <chart_version> --set name=reconciler --set global.edpName=<edp-project> --set global.platform=<platform_type> --set global.database.name=<db-name> --set global.database.host=<db-name>.<namespace_name> --set global.database.port=<port>
    ```
5. Check the <edp-project> namespace that should contain operator deployment with your operator in a running status.

## Local Development

In order to develop the operator, first set up a local environment. For details, please refer to the [Local Development](https://epam.github.io/edp-install/developer-guide/local-development/) page.

For development process, are available snapshot versions of component. For details, please refer to the [snapshot helm chart repository](https://epam.github.io/edp-helm-charts/snapshot/) page.

### Related Articles

* [Install EDP](https://epam.github.io/edp-install/operator-guide/install-edp/)
