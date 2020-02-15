APP_NAME=terraform-do
APP_PATH=cmd/${APP_NAME}

.PHONY: ${APP_PATH}/${APP_NAME}

all: build

build: ${APP_PATH}/${APP_NAME}

clean:
	[[ -e ${APP_PATH}/${APP_NAME} ]] && rm ${APP_PATH}/${APP_NAME} || true

${APP_PATH}/${APP_NAME}:
	go mod vendor
	cd ${APP_PATH} && go build -mod vendor
