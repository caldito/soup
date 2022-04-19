# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/caldito/soup/tree/develop)

## [v0.4.0](https://github.com/caldito/soup/releases/tag/v0.4.0) - 2022-04-19
### Added
- Add ARM64 support
- Use patterns for getting k8s manifests
- Unit tests

## [v0.3.1](https://github.com/caldito/soup/releases/tag/v0.3.1) - 2021-10-02
### Added
- Refactor application to separate the modules clearly
- Improve the Readme

## [v0.3.0](https://github.com/caldito/soup/releases/tag/v0.3.0) - 2021-07-24
### Added
- Improve error handling when the .soup.yml file does not exist and can not be parsed
- Extracted DoSSA and DeclareNamespaceSSA functions to `/pkg/k8s` directory to be imported easily by other projects
- Refactoring
### Removed
- Unused Makefile targets

## [0.2.0](https://github.com/caldito/soup/releases/tag/0.2.0) - 2021-06-30
### Added
- Deploy to cluster (main functionality working)
- Cluster role and cluster role binding for soup pod
- Make target for managing dependencies
### Removed
- Unnecesary make installation in pipeline
### Fixed
- Installation in cluster

## [0.1.0](https://github.com/caldito/soup/releases/tag/0.1.0) - 2021-04-02
### Added
- Program that clones repo and processes the build configuration in each branch
- Makefile
- Container support
- Installation in clusters
- Readme, license and changelog
- Github workflow to push container images
