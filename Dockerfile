FROM golang:latest AS build


WORKDIR $GOPATH/src/github.com/mchirico/k8sv19
COPY . $GOPATH/src/github.com/mchirico/k8sv19

RUN cd examples/in-cluster-client-configuration && go build  -o /bin/project

