PROJECT = septapig
NAME = k8sv19
TAG = dev


.PHONY: docker-build
docker-build:
        docker build --no-cache -t us.gcr.io/$(PROJECT)/$(NAME):$(TAG) -f Dockerfile .

.PHONY: kind
kind:
	kind load docker-image gcr.io/$(PROJECT)/$(NAME):$(TAG)


.PHONY: calico
calico:
	kind delete cluster
	kind create cluster --config calico/kind-calico.yaml
	kubectl apply -f calico/ingress-nginx.yaml
	kubectl apply -f calico/tigera-operator.yaml
	kubectl apply -f calico/calicoNetwork.yaml
	kubectl apply -f calico/calicoctl.yaml


.PHONY: jaeger
jaeger:
	kubectl create namespace observability 
	kubectl create -f jaeger/jaegertracing.io_jaegers_crd.yaml # <2>
	kubectl create -n observability -f jaeger/service_account.yaml
	kubectl create -n observability -f jaeger/role.yaml
	kubectl create -n observability -f jaeger/role_binding.yaml
	kubectl create -n observability -f jaeger/operator.yaml
	kubectl create -f jaeger/cluster_role.yaml
	kubectl create -f jaeger/cluster_role_binding.yaml


.PHONY: calicoctl
calicoctl:
	 kubectl exec -ti -n kube-system calicoctl -- /calicoctl get profiles -o wide


.PHONY: cert-manager
cert-manager:
	kind delete cluster
	kind create cluster --config calico/kind-calico.yaml
	kubectl apply -f calico/ingress-nginx.yaml
	kubectl apply -f calico/tigera-operator.yaml
	kubectl apply -f calico/calicoNetwork.yaml
	kubectl apply -f calico/calicoctl.yaml
	kubectl apply -f calico/cert-manager.yaml
#	kubectl kudo init



.PHONY: buildImage
buildImage:
	rm -rf ${GOPATH}/src/k8s.io/kubernetes
	export GO111MODULE=off
	go get k8s.io/kubernetes || true
	export GO111MODULE=on
	go get sigs.k8s.io/kind
	cd ${GOPATH}/src/k8s.io/kubernetes && git checkout v1.19.0-rc.4
	kind build node-image --image=v1.19.0-rc.4



.PHONY: cert-manager-v1.19
cert-manager-v1.19:
	kind delete cluster
	kind create cluster --image=v1.19.0-rc.4 --config calico/kind-calico.yaml
	kubectl apply -f calico/ingress-nginx.yaml
	kubectl apply -f calico/tigera-operator.yaml
	kubectl apply -f calico/calicoNetwork.yaml
	kubectl apply -f calico/calicoctl.yaml
	kubectl apply -f calico/cert-manager.yaml
#	kubectl kudo init


.PHONY: cert-manager
cert-manager-mce:
	kind delete cluster
	kind create cluster --config calico-mce/kind-calico.yaml
	kubectl apply -f calico-mce/ingress-nginx.yaml
	kubectl apply -f calico/tigera-operator.yaml
	kubectl apply -f calico/calicoNetwork.yaml
	kubectl apply -f calico/calicoctl.yaml
	kubectl apply -f calico/cert-manager.yaml
#	kubectl kudo init



.PHONY: packages
packages:
	kubectl kudo init
	kubectl kudo install zookeeper
	kubectl kudo install kafka
	kubectl kudo install redis
	kubectl kudo install mysql
	kubectl kudo install rabbitmq


.PHONY: sample
sample:
	rm -rf express.zip
	rm -rf node-starter-express
	curl -LO https://github.com/mchirico/node-starter/archive/express.zip
	unzip express.zip
	cd node-starter-express
	npm install


# Danger This Cleans Up Everything!
.PHONY: cleanup
cleanup:
	docker stop $(docker ps -aq)
	docker rm $(docker ps -aq)
	docker rmi $(docker images -q)


.PHONY: docker-build
docker-build:
	docker build --no-cache -t gcr.io/$(PROJECT)/$(NAME):$(TAG) -f Dockerfile .

.PHONY: kind
kind:
	kind load docker-image gcr.io/$(PROJECT)/$(NAME):$(TAG)

.PHONY: push
push:
	docker push gcr.io/$(PROJECT)/$(NAME):$(TAG) 

.PHONY: pull
pull:
	docker pull gcr.io/$(PROJECT)/$(NAME):$(TAG) 



.PHONY: run
run:
	docker run -p 3000:3000 --rm -it -d --name $(NAME) gcr.io/$(PROJECT)/$(NAME):$(TAG) 


.PHONY: runnod
runnod:
	docker run -p 3000:3000 --rm -it --name $(NAME) gcr.io/$(PROJECT)/$(NAME):$(TAG) 

.PHONY: stop
stop:
	docker stop $(NAME)


.PHONY: logs
logs:
	docker logs $(NAME)



