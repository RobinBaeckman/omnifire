## Local Setup
```sh
$ cd deploy/infra/environment/dev
$ terraform init
$ terraform apply
```

### Kafka
after all instances are running:
1. port-forward mykafka-cp-kafka-connect-* pod on port 8083
2. install postman/Omnifire.postman_collection.json and postman/Kafka.postman_environment.json
3. run "Add postgres connector" from kafka connect postman directory
4. run "Get postgres connector status" to make sure the connector is running and working
5. create an insert record in postgres student table
6. port-forward mykafka-cp-control-center-* pod and visit http://localhost:9021, choose cluster, topic, message to see the message that was created from the insert

####notes:
* checkout deploy/infra/charts/cp-helm-charts/templates/NOTES.txt how to smoke test kafka instances.
* to be able to add connectors to kafka connect a new image is created from the original and connectors are copied to the new image, use `make build-kafka` to build and push the image

## Architecture
tracing:
hopbox-1 -|
	        -----> open-telementry -> tempo -> grafana
hopbox-n -|

logging:
hopbox-1 -|
	        -----> promtail -> loki -> grafana
hopbox-n -|


### TODO/Thoughts (create as issues and/or add to project manager tool)
* do we need telegraf?
* Do we need influxdb, what for?
* * configure influxdb default values
* * add it as a datasource in grafana helm above
* add tracing open telemetry
* add data analytics tools
* install grafana tempo and loki
* make all apps accessible through the browser with hostname
* traefik prometheus metrics?
* add access points for influxdb and prometheus
* make postgres data source available from beginning in grafana
* make ports accessible from beginning on remote k8 host for important services like postgres etc
* fix so traefiks external service address is set automaticly
* TODO fix local init: $ make local-init
