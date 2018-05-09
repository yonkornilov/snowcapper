.PHONY : binary vagrant clean

default: binary vagrant

get:
	go get ./...

test:
	go test -v ./...

binary:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

vagrant:
	vagrant up --provision

clean:
	vagrant destroy --force && rm -f snowcapper
