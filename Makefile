
# SERVICE_NAME = api
#SERVICE_NAME = user
SERVICE_NAME = member
SERVER_PREFIX = src/server

run:
	cd $(SETVICE_PATH) && go run main.go

service:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto $(SERVICE_NAME)_model.proto $(SERVICE_NAME)_service.proto --svcout ../


api:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto api_model.proto api_service.proto --svcout ../

member:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto member_model.proto member_service.proto --svcout ../

all:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto member_model.proto member_service.proto --svcout ../
	cd $(SERVER_PREFIX)/pb && truss -v common.proto member_model.proto member_service.proto --svcout ../


model:
	cd $(SERVER_PREFIX)/pb && truss -v common.proto
	cd $(SERVER_PREFIX)/pb && truss -v *_model.proto

clean:
	rm -fr $(SERVER_PREFIX)/*-service
	cd $(SERVER_PREFIX)/pb && rm -fr *.go && rm -fr common.pb.go