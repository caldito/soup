# Soup

![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/caldito/soup.svg)](https://pkg.go.dev/github.com/caldito/soup)
[![Go Report Card](https://goreportcard.com/badge/github.com/caldito/soup)](https://goreportcard.com/report/github.com/caldito/soup)
[![release](https://img.shields.io/github/release/caldito/soup/all.svg)](https://github.com/caldito/soup/releases)

Soup is a GitOps operator for Kubernetes.

## Requirements

`Kubernetes >= v1.20`

## Installation

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


The command line arguments should be in this form:
```
-arg=value
```

### Available arguments
* repo: the url of the repo. This must be specified
* interval: the sync interval in seconds. By default is set to 120s.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
