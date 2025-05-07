.PHONY: all
all:
	go build -o playlog main.go

.PHONY: clean
clean:
	-rm playlog

.PHONY: test
test:
	go test -cover ./...
