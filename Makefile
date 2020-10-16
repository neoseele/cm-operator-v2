IMAGE_NAME ?= cm-operator
IMAGE_TAG ?= v0.0.1
REPO_NAME ?= neoseele

IMAGE := $(shell docker image inspect $(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG) &>/dev/null || echo missing)

.PHONY: build
build:

ifeq ($(IMAGE),missing)
	@echo "building image [$(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)] ..."
	@operator-sdk build $(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)
else
	@echo "image [$(IMAGE_NAME):$(IMAGE_TAG)] already exists."
endif

.PHONY: build-dockerhub
build-dockerhub: build # depend on build
	@echo "pushing image [$(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)] ..."
	@docker push $(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: clean
clean:
	@echo "removing images ..."
	-@docker rmi $(REPO_NAME)/$(IMAGE_NAME):$(IMAGE_TAG)

deploy-operator:
	@kubectl apply -f deploy/service_account.yaml
	@kubectl apply -f deploy/role.yaml
	@kubectl apply -f deploy/role_binding.yaml
	@kubectl apply -f deploy/operator.yaml

teardown-operator:
	-@kubectl apdeleteply -f deploy/operator.yaml
	-@kubectl delete -f deploy/role_binding.yaml
	-@kubectl delete -f deploy/role.yaml
	-@kubectl delete -f deploy/service_account.yaml

deploy-cr:
	@kubectl apply -f deploy/crds/cm.example.com_custommetrics_crd.yaml
	@kubectl apply -f deploy/crds/cm.example.com_v1alpha1_custommetric_cr.yaml

teardown-cr:
	-@kubectl delete -f deploy/crds/cm.example.com_v1alpha1_custommetric_cr.yaml
	-@kubectl delete -f deploy/crds/cm.example.com_custommetrics_crd.yaml
