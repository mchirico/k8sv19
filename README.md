
# k8sv19


## Steps for running

### Build KinD image for Kubernetes v1.19.0-rc.4
```
rm -rf ${GOPATH}/src/k8s.io/kubernetes
export GO111MODULE=off
go get k8s.io/kubernetes || true
export GO111MODULE=on

go get sigs.k8s.io/kind
cd ${GOPATH}/src/k8s.io/kubernetes && git checkout v1.19.0-rc.4
kind build node-image --image=v1.19.0-rc.4

```

### Create Cluster

```
kind delete cluster
kind create cluster --image=v1.19.0-rc.4 --config calico/kind-calico.yaml
kubectl apply -f calico/ingress-nginx.yaml
kubectl apply -f calico/tigera-operator.yaml
kubectl apply -f calico/calicoNetwork.yaml
kubectl apply -f calico/calicoctl.yaml
kubectl apply -f calico/cert-manager.yaml

```

### Build and Run Program
```
go run main.go
```

