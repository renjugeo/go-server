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

```