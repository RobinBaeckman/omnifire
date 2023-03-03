SHELL=/bin/bash -o pipefail
help:
	# make build-kafka: builds the kafka image to be able to have connectors on startup

.PHONY: build-kafka
build-kafka:
	v=$$(echo $$RANDOM | md5sum | head -c 20) && \
	echo "version: " $$v && \
	docker buildx build --platform linux/amd64 -t robinbaeckman/kafka-connect:$$v . && \
	docker push robinbaeckman/kafka-connect:$$v
