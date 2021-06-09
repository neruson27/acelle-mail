BUILD_FILE=acelle-mail

RUN_FILE=./$(BUILD_FILE)

build:
	go build -ldflags "-s -w -extldflags '-static'" -o $(BUILD_FILE)

clean:
	rm -f $(BUILD_FILE)

all: clean build
	$(RUN_FILE) help

listener: clean build
	$(RUN_FILE) listener

server: clean build
	$(RUN_FILE) server

jobs: clean build
	$(RUN_FILE) jobs