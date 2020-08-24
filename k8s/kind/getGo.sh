#!/bin/bash
apt-get update
apt-get install -y wget git
mkdir -p /root/go/bin


wget https://golang.org/dl/go1.15.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz
rm -f go1.15.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go

git clone https://github.com/mchirico/k8sv19.git

cd k8sv19/
go install ./...
