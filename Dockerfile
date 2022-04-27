FROM golang:1.17.9-buster as build
RUN apt update -y && apt upgrade -y && update-ca-certificates
WORKDIR /go/src/github.com/renjugeo/go-server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go-server

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go-server /
ENTRYPOINT [ "/go-server" ]