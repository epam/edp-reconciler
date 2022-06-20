<a name="unreleased"></a>
## [Unreleased]

### Routine

- Update current development version [EPMDEDP-8832](https://jiraeu.epam.com/browse/EPMDEDP-8832)
- Update chart annotation [EPMDEDP-9515](https://jiraeu.epam.com/browse/EPMDEDP-9515)


<a name="v2.11.0"></a>
## [v2.11.0] - 2022-05-25
### Features

- Update Makefile changelog target [EPMDEDP-8218](https://jiraeu.epam.com/browse/EPMDEDP-8218)
- Generate CRDs and helm docs automatically [EPMDEDP-8385](https://jiraeu.epam.com/browse/EPMDEDP-8385)

### Bug Fixes

- Fix changelog generation in GH Release Action [EPMDEDP-8468](https://jiraeu.epam.com/browse/EPMDEDP-8468)

### Routine

- Update release CI pipelines [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)
- Populate chart with Artifacthub annotations [EPMDEDP-8049](https://jiraeu.epam.com/browse/EPMDEDP-8049)
- Update base docker image to alpine 3.15.4 [EPMDEDP-8853](https://jiraeu.epam.com/browse/EPMDEDP-8853)
- Update "github.com/epam/edp-cd-pipeline-operator/v2" package [EPMDEDP-8929](https://jiraeu.epam.com/browse/EPMDEDP-8929)
- Update changelog [EPMDEDP-9185](https://jiraeu.epam.com/browse/EPMDEDP-9185)


<a name="v2.10.0"></a>
## [v2.10.0] - 2021-12-06
### Features

- Provide operator's build information [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)

### Bug Fixes

- Remove unknown field apiGroup in Openshift RB [EPMDEDP-7279](https://jiraeu.epam.com/browse/EPMDEDP-7279)
- Changelog links [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)

### Code Refactoring

- Expand reconciler role [EPMDEDP-7279](https://jiraeu.epam.com/browse/EPMDEDP-7279)
- Add namespace field in roleRef in RB, align RB name [EPMDEDP-7279](https://jiraeu.epam.com/browse/EPMDEDP-7279)
- Remove unnecessary namespace field in roleRef [EPMDEDP-7279](https://jiraeu.epam.com/browse/EPMDEDP-7279)
- Replace cluster-wide role/rolebinding to namespaced [EPMDEDP-7279](https://jiraeu.epam.com/browse/EPMDEDP-7279)
- Refactor pointers in sql transaction functions. [EPMDEDP-7943](https://jiraeu.epam.com/browse/EPMDEDP-7943)
- Address golangci-lint issues [EPMDEDP-7945](https://jiraeu.epam.com/browse/EPMDEDP-7945)

### Routine

- Add changelog generator [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)
- Add codecov report [EPMDEDP-7885](https://jiraeu.epam.com/browse/EPMDEDP-7885)
- Update docker image [EPMDEDP-7895](https://jiraeu.epam.com/browse/EPMDEDP-7895)
- Use custom go build step for operator [EPMDEDP-7932](https://jiraeu.epam.com/browse/EPMDEDP-7932)
- Update go to version 1.17 [EPMDEDP-7932](https://jiraeu.epam.com/browse/EPMDEDP-7932)

### Documentation

- Update the links on GitHub [EPMDEDP-7781](https://jiraeu.epam.com/browse/EPMDEDP-7781)


<a name="v2.9.0"></a>
## [v2.9.0] - 2021-12-03

<a name="v2.8.1"></a>
## [v2.8.1] - 2021-12-03

<a name="v2.8.0"></a>
## [v2.8.0] - 2021-12-03

<a name="v2.7.1"></a>
## [v2.7.1] - 2021-12-03

<a name="v2.7.0"></a>
## [v2.7.0] - 2021-12-03

[Unreleased]: https://github.com/epam/edp-reconciler/compare/v2.11.0...HEAD
[v2.11.0]: https://github.com/epam/edp-reconciler/compare/v2.10.0...v2.11.0
[v2.10.0]: https://github.com/epam/edp-reconciler/compare/v2.9.0...v2.10.0
[v2.9.0]: https://github.com/epam/edp-reconciler/compare/v2.8.1...v2.9.0
[v2.8.1]: https://github.com/epam/edp-reconciler/compare/v2.8.0...v2.8.1
[v2.8.0]: https://github.com/epam/edp-reconciler/compare/v2.7.1...v2.8.0
[v2.7.1]: https://github.com/epam/edp-reconciler/compare/v2.7.0...v2.7.1
[v2.7.0]: https://github.com/epam/edp-reconciler/compare/v2.3.0-112...v2.7.0
