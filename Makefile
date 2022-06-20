APP_NAME=monarch-be
APP_IMG=monarch-img
APP_BIN=server
APP_NETWORK=monarch-net
BUILD_DIR=bin

build:
	if mkdir $(BUILD_DIR); then \
		echo "bin directory created"; \
	else \
		echo "bin directory already exists";\
	fi
	
	if go build -o $(BUILD_DIR)/$(APP_BIN) main.go; then \
		echo "build successfully"; \
	else \
		echo "build failed"; \
	fi

clean:
	rm $(BUILD_DIR)/$(APP_BIN)

run:
	./$(BUILD_DIR)/$(APP_BIN)

build-docker:
	docker build -t $(APP_IMG) .

network-docker:
	if docker network create $(APP_NETWORK); then\
		echo "network created"; \
	else \
		echo "network already exists"; \
	fi

clean-docker:
	if docker stop $(APP_NAME); then \
		docker rm $(APP_NAME); \
	fi

run-docker: clean-docker
	docker run --name $(APP_NAME) -p 8000:8000 --restart unless-stopped -d $(APP_IMG)
	docker network connect monarch-net monarch-be