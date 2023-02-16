doc:
	# make local-init: initialize infrastructure for kubernetes
	# make proto     : generate proto files

.PHONY: local-init
local-init:
	cd deploy/infra; \
	terraform init; \
	terraform apply -auto-approve; \
	cd ../../; \
	cd api && $(MAKE) local-deploy; \
	cd ..; \
	cd box && $(MAKE) local-deploy;

.PHONY: proto
proto:
	protoc -I ./proto --go_out ./proto --go_opt paths=source_relative --go-grpc_out ./proto --go-grpc_opt paths=source_relative --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative ./proto/box/box.proto
	protoc -I ./proto --go_out ./proto --go_opt paths=source_relative --go-grpc_out ./proto --go-grpc_opt paths=source_relative --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative ./proto/api/api.proto
