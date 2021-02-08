SERVICE_NAME = echo
SERVER_PREFIX = src/server

run:
	cd $(SETVICE_PATH) && go run main.go

service:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto $(SERVICE_NAME)_model.proto $(SERVICE_NAME)_service.proto --svcout ../

model:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto
	cd $(SERVER_PREFIX)/pb && truss -v *_model.proto

clean:
	rm -fr $(SERVER_PREFIX)/$(SERVICE_NAME)-service
	cd $(SERVER_PREFIX)/pb && rm -fr $(SERVICE_NAME)*.go && rm -fr common.pb.go
	rm -fr $(SERVER_PREFIX)/$(SERVICE_NAME)-service