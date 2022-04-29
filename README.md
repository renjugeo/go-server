# go-server

## About
Golang based HTTP server which accepts GET requests with input parameter as "sortKey" and "limit". The server queries URLs in the config file, combines the results from all URLs, sorts them by the sortKey and returns the response

## How to build locally

```bash
# clone the repo
git clone git@github.com:renjugeo/go-server.git
cd go-server

# build
go build . -o ./go-server

# run 
./go-server -config-path ./config/config.yaml

# query endpoints
curl http://localhost:8080/api/v1/stats

curl 'http://localhost:8080/api/v1/stats?sortKey=views&limit=5'

```

## Steps to run in a local kubernetes cluster

### Install helm

Install helm by following the documetation https://helm.sh/docs/intro/install/#through-package-managers

### Install kind

Kind (Kubernetes IN Docker) can be used to create a kubernetes cluster using docker
Install kind following the documentation https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager

### Create a local kubernetes cluster

```bash
kind create cluster
```

### Build the docker image

From the go-server folder run the following command

```bash
docker build -t go-server:latest
```

### Load the docker image in to kind cluster

```bash
kind load docker-image go-server:latest
```

### Deploy the helm chart to kind cluster

From the go-server folder run the following commands

```bash
cd deploy/go-server
helm install app .

# verify the release
helm list

# verify the pod is up and running
kubectl get pod
```

### Access the api

```bash
# port forward the kubernets service to access the endpoint locally
kubectl port-forward svc/app-go-server 8080

# curl the api
curl http://localhost:8080/api/v1/stats
```