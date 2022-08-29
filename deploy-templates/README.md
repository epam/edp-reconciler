# reconciler

![Version: 2.13.0-SNAPSHOT](https://img.shields.io/badge/Version-2.13.0--SNAPSHOT-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 2.13.0-SNAPSHOT](https://img.shields.io/badge/AppVersion-2.13.0--SNAPSHOT-informational?style=flat-square)

A Helm chart for EDP Reconciler

**Homepage:** <https://epam.github.io/edp-install/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| epmd-edp | <SupportEPMD-EDP@epam.com> | <https://solutionshub.epam.com/solution/epam-delivery-platform> |
| sergk |  | <https://github.com/SergK> |

## Source Code

* <https://github.com/epam/edp-reconciler>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| global.database.host | string | `"edp-db"` | database host |
| global.database.name | string | `"edp-db"` | database name |
| global.database.port | int | `5432` | database port |
| global.edpName | string | `""` | namespace or a project name (in case of OpenShift) |
| global.platform | string | `"openshift"` | platform type that can be "kubernetes" or "openshift" |
| image.repository | string | `"epamedp/reconciler"` | EDP reconciler Docker image name. The released image can be found on [Dockerhub](https://hub.docker.com/r/epamedp/reconciler) |
| image.tag | string | `nil` | EDP reconciler Docker image tag. The released image can be found on [Dockerhub](https://hub.docker.com/r/epamedp/reconciler/tags) |
| imagePullPolicy | string | `"IfNotPresent"` |  |
| name | string | `"reconciler"` | component name |
| nodeSelector | object | `{}` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"25m"` |  |
| resources.requests.memory | string | `"32Mi"` |  |
| tolerations | list | `[]` |  |

