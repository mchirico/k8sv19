PROJECT = septapig
NAME = kindstarter
TAG = dev

export GOPATH=${HOME}/gopath
export GOBIN=${HOME}/gopath/bin


.PHONY: cert-manager-v1.19
cert-manager-v1.19:
	mkdir -p ${HOME}/gopath
	go get k8s.io/kubernetes || true
	cd ${HOME}/gopath/src/k8s.io/kubernetes && git checkout v1.19.0-rc.4
	go get sigs.k8s.io/kind
	export PATH=${HOME}/bin:${PATH}
#     Node image
	kind build node-image --image=v1.19.0-rc.4
