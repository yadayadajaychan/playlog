.PHONY: all
all: frontend backend

.PHONY: frontend
frontend:
	npm run build

.PHONY: backend
backend:
	go build -ldflags "-X main.programVersion=$$(<VERSION)" -o playlog main.go

.PHONY: clean
clean:
	-rm playlog
	-rm -r build

.PHONY: test
test:
	go test -cover -count=1 ./...
