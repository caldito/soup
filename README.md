# Soup

Soup is a GitOps operator for Kubernetes.

## Installation

Install in 3 easy steps:

1. First get the kubernetes kubernetes manifest for installing
```
curl -O https://raw.githubusercontent.com/caldito/soup/main/install.yml
```
2. Override the `repo` argument in the file you just downloaded
3. Apply to the cluster
```
kubectl apply -f install.yml
```

## Usage


The command line arguments can be in this forms:
```
-arg value
```
```
-arg=value
```

### Available arguments
* repo: the url of the repo. This must be specified
* interval: the sync interval in seconds. By default is set to 120s.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
![Apache 2.0 License](https://img.shields.io/hexpm/l/plug.svg)

This project is licensed under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license - see the [LICENSE](LICENSE) file for details.