all: install

clean:
	go clean ./...

doc:
	godoc -http=:6060

install:
	go get github.com/labstack/echo
	go get gopkg.in/mgo.v2
	go get gopkg.in/mgo.v2/bson
	go get gopkg.in/yaml.v2
	go get github.com/spf13/viper
	go get github.com/jordan-wright/email

test-install: install
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

build-install: install test-install
	go get github.com/mitchellh/gox

test:
	go test -cover ./...

test-travis:
	./test-cover.sh

bench:
	go test -run=none -bench=. -benchmem ./...

build:
	go build -o $(GOPATH)/bin/scds .

# Build and tag binaries for each OS and architecture.
build-all: build
	mkdir -p bin

	gox -output="bin/scds-{{.OS}}.{{.Arch}}" \
		-os="linux windows darwin" \
		-arch="amd64" \
		./cmd/origins > /dev/null

fmt:
	go vet ./...
	go fmt ./...

lint:
	golint ./...


.PHONY: test
