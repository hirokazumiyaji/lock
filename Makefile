VERSION=1.0.0

build: cmd/lock/main.go lock/*.go
	go build -ldflags "-X main.version=${VERSION}" -o bin/lock cmd/lock/main.go

clean:
	rm -rf bin/*
