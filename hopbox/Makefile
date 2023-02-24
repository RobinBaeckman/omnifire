SHELL=/bin/bash -o pipefail
help:
	# make local-deploy: use the local env overlays and deploy k8 manifests
	# make migrateup   : add migration database schemas
	# make migratedown : remove migration database schemas

.PHONY: deploy
deploy:
	v=$$(echo $$RANDOM | md5sum | head -c 20) && \
	echo "version: " $$v && \
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build && \
	docker buildx build --platform linux/amd64 -t robinbaeckman/hopbox:$$v . && \
	docker push robinbaeckman/hopbox:$$v && \
	rm hopbox && \
	cd deploy/env/local/box-1 && \
	kustomize edit set image robinbaeckman/hopbox=robinbaeckman/hopbox:$$v && kustomize build --load_restrictor=none | kubectl apply -f - && \
	cd ../../../.. && \
	cd deploy/env/local/box-2 && \
	kustomize edit set image robinbaeckman/hopbox=robinbaeckman/hopbox:$$v && kustomize build --load_restrictor=none | kubectl apply -f - && \
	cd ../../../..

.PHONY: migrateup
migrateup:
	migrate -path storage/migration -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" -verbose up

.PHONY: migratedown
migratedown:
	migrate -path storage/migration -database "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable" -verbose down
