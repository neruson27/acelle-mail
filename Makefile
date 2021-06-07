all: clean build
	./acelle-mail help

build:
	go build -ldflags "-s -w -extldflags '-static'" -o acelle-mail

clean:
	rm -f acelle-mail

listener: clean build
	./acelle-mail listener

server: clean build
	./acelle-mail server