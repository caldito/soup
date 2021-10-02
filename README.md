# Soup

![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/caldito/soup.svg)](https://pkg.go.dev/github.com/caldito/soup)
[![Go Report Card](https://goreportcard.com/badge/github.com/caldito/soup)](https://goreportcard.com/report/github.com/caldito/soup)
[![release](https://img.shields.io/github/release/caldito/soup/all.svg)](https://github.com/caldito/soup/releases)

Soup is a GitOps operator for Kubernetes.

## Why Soup
* Focused on design and usage simplicity.
* Deployment to kubernetes performed with [Server-Side Apply](https://kubernetes.io/docs/reference/using-api/server-side-apply/), allows a simple yet fully capable deployment module for all type of resources.
* Match branch names with Regex.
* Option to create namespaces called the same way as the branch. Useful in combination with the regex branch selector.

Image [docs/images/overview.png](https://github.com/caldito/soup/blob/main/docs/images/overview.png) shows an overview on how the system works.

You can watch the following [demo](https://asciinema.org/a/439238) of Soup working.

## Getting started
### Prerequisites

`Kubernetes >= v1.20`

### Installation

Install in 3 easy steps:

1. First get the kubernetes kubernetes manifest for installing
```
curl -O https://raw.githubusercontent.com/caldito/soup/main/manifests/install.yml
```
2. Override the `repo` argument in the file you just downloaded
3. Apply to the cluster
```
kubectl apply -f install.yml
```

## Usage

### Deployment Configuration file
Soup reads a file named `.soup.yml` on the repository branches in order to know what to deploy and in which namespace. An example can be found in [docs/examples/.soup.yml](https://github.com/caldito/soup/blob/main/docs/examples/.soup.yml).

### Arguments
The command line arguments should be in this form:
```
-arg=value
```
Available arguments:
* repo: the url of the repo. This must be specified
* interval: the sync interval in seconds. By default is set to 120s.

## Contributing
Pull requests are welcomed and encouraged. For questions, feature requests and bug reports, please open an issue.

There is also a [TODO](https://github.com/caldito/soup/blob/main/TODO) file containing work planned to do.

An internal diagram which may be useful for contributors can be found in [docs/images/internals.png](https://github.com/caldito/soup/blob/main/docs/images/internals.png).
