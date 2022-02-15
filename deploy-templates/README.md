# reconciler

![Version: 2.11.0-SNAPSHOT](https://img.shields.io/badge/Version-2.11.0--SNAPSHOT-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 2.11.0-SNAPSHOT](https://img.shields.io/badge/AppVersion-2.11.0--SNAPSHOT-informational?style=flat-square)

A Helm chart for EDP Reconciler

**Homepage:** <https://epam.github.io/edp-install/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| epmd-edp | SupportEPMD-EDP@epam.com | https://solutionshub.epam.com/solution/epam-delivery-platform |
| sergk |  | https://github.com/SergK |

## Source Code

* <https://github.com/epam/edp-reconciler>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| global.database.host | string | `"edp-db"` |  |
| global.database.name | string | `"edp-db"` |  |
| global.database.port | int | `5432` |  |
| global.edpName | string | `""` |  |
| global.platform | string | `"openshift"` |  |
| image.name | string | `"epamedp/reconciler"` |  |
| image.version | string | `nil` |  |
| imagePullPolicy | string | `"IfNotPresent"` |  |
| name | string | `"reconciler"` |  |
| nodeSelector | object | `{}` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"25m"` |  |
| resources.requests.memory | string | `"32Mi"` |  |
| tolerations | list | `[]` |  |

